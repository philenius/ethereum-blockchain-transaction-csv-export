package main

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/mgutz/logxi/v1"
	"github.com/onrik/ethrpc"

	"ethereum-blockchain-transaction-csv-export/work"
)

func main() {

	start := time.Now()

	var c config
	err := envconfig.Process("", &c)
	if err != nil {
		log.Error("configuration error", "err", err.Error())
		buf := new(bytes.Buffer)
		envconfig.Usagef("", &c, buf, envconfig.DefaultListFormat)
		log.Info("usage", buf.String())
		os.Exit(1)
	}

	clientAddr := fmt.Sprintf("http://%s:%d", c.Hostname, c.Port)
	client := ethrpc.NewEthRPC(clientAddr)
	version, err := client.Web3ClientVersion()
	if err != nil {
		log.Error("failed to connect to Ethereum node", "err", err.Error())
		os.Exit(2)
	}
	log.Info("successfully connected to Ethereum node", "host", c.Hostname, "port", c.Port, "version", version)

	blockHeightChan := make(chan *work.Job, 10000)
	txHashChan := make(chan *work.Job, 10000)
	txChan := make(chan *work.Job, 10000)
	failedBlockChan := make(chan *work.Job, 10000)
	failedTxChan := make(chan *work.Job, 10000)
	statsChan := make(chan *work.Stat, 10000)
	wt := sync.WaitGroup{}

	latestBlockHeight := 0
	fetchedBlockCount := 0
	fetchedTxCount := int64(0)

	// receive all stats and store result in vars
	go func() {
		for stat := range statsChan {
			if stat.FetchedBlock {
				fetchedBlockCount++
				latestBlockHeight = stat.BlockHeight
			}
			if stat.FetchedTransaction {
				fetchedTxCount++
			}
		}
	}()

	// display stats periodically
	go func() {
		lastFetchedBlockCount := 0
		lastTransactionCount := int64(0)

		for range time.NewTicker(time.Second * time.Duration(c.StatsIntervalInSeconds)).C {
			blockDiff := fetchedBlockCount - lastFetchedBlockCount
			blockRate := float32(blockDiff) / float32(c.StatsIntervalInSeconds)

			txDiff := fetchedTxCount - lastTransactionCount
			txRate := float32(txDiff) / float32(c.StatsIntervalInSeconds)
			log.Warn("stats",
				"totalFetchedBlocks", fetchedBlockCount,
				"lastestFetchedBlock", latestBlockHeight,
				"blockFetchRatePerSec", blockRate,
				"totalFetchedTransactions", fetchedTxCount,
				"txFetchRatePerSec", txRate,
				"chanLengthRemainingBlockJobs", len(blockHeightChan),
				"chanLengthRemainingTransactionJobs", len(txHashChan),
				"chanLengthRemainingTransactionFileWrites", len(txChan),
			)
			lastFetchedBlockCount = fetchedBlockCount
			lastTransactionCount = fetchedTxCount
		}
	}()

	// list all blocks to fetch
	go func() {
		endHeight := c.StartBlockHeight + c.BlockCount
		for i := c.StartBlockHeight; i < endHeight; i++ {
			blockHeightChan <- &work.Job{BlockHeight: i}
		}
		log.Info("finished listing block numbers")
		close(blockHeightChan)
		wt.Done()
	}()
	wt.Add(1)

	// fetch all blocks
	go func() {
		p := work.NewBlockWorkerPool(
			c.WorkerCountForBlocks, clientAddr, blockHeightChan, txHashChan, failedBlockChan, statsChan,
		)
		p.Run()
		wt.Done()
	}()
	wt.Add(1)

	// fetch all transactions
	go func() {
		p := work.NewTxWorkerPool(c.WorkerCountForTransactions, clientAddr, txHashChan, txChan, failedTxChan, statsChan)
		p.Run()
		wt.Done()
	}()
	wt.Add(1)

	go func() {
		exportFailedBlockJobs(failedBlockChan)
		wt.Done()
	}()
	wt.Add(1)

	go func() {
		exportFailedTxJobs(failedTxChan)
		wt.Done()
	}()
	wt.Add(1)

	go func() {
		exportAsCSV(txChan)
		wt.Done()
	}()
	wt.Add(1)

	log.Info("waiting for application to fetch all data from Ethereum node")
	wt.Wait()

	elapsed := time.Since(start).Round(time.Minute)
	log.Warn("application completed successfully",
		"durationInMinutes", int(elapsed.Minutes()),
		"startHeight", c.StartBlockHeight,
		"endHeight", c.StartBlockHeight+c.BlockCount,
		"blockCount", c.BlockCount,
	)
	close(statsChan)

}
