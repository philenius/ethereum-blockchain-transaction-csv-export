package main

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/mgutz/logxi/v1"
	"github.com/onrik/ethrpc"

	"github.com/philenius/ethereum-blockchain-transaction-csv-export/work"
)

func main() {

	start := time.Now()

	blockConcurr := flag.Int("blockConcurr", 10, "The count of concurrent workers for fetching blocks.")
	blockCount := flag.Int("count", 1000, "The total amount of blocks to fetch.")
	hostname := flag.String("host", "127.0.0.1", "The hostname / IP address of the Ethereum node.")
	port := flag.Int("port", 8545, "The port number of the Ethereum node.")
	startHeight := flag.Int("start", 0, "The height / number of the first block to fetch.")
	statsInterval := flag.Int("statsIntervalSec", 5, "The invertal in seconds to display stats.")
	txConcurr := flag.Int("txConcurr", 20, "The count of concurrent workers for fetching transactions.")

	flag.Parse()

	clientAddr := fmt.Sprintf("http://%s:%d", *hostname, *port)
	client := ethrpc.NewEthRPC(clientAddr)
	version, err := client.Web3ClientVersion()
	if err != nil {
		log.Fatal("failed to connect to Ethereum node", "err", err.Error())
	}
	log.Info("successfully connected to Ethereum node", "host", *hostname, "port", *port, "version", version)

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

		for range time.NewTicker(time.Second * time.Duration(*statsInterval)).C {
			blockDiff := fetchedBlockCount - lastFetchedBlockCount
			blockRate := float32(blockDiff) / float32(*statsInterval)

			txDiff := fetchedTxCount - lastTransactionCount
			txRate := float32(txDiff) / float32(*statsInterval)
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
		endHeight := *startHeight + *blockCount
		for i := *startHeight; i < endHeight; i++ {
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
			*blockConcurr, clientAddr, blockHeightChan, txHashChan, failedBlockChan, statsChan,
		)
		p.Run()
		wt.Done()
	}()
	wt.Add(1)

	// fetch all transactions
	go func() {
		p := work.NewTxWorkerPool(*txConcurr, clientAddr, txHashChan, txChan, failedTxChan, statsChan)
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
		"startHeight", *startHeight,
		"endHeight", *startHeight+*blockCount,
		"blockCount", *blockCount,
	)
	close(statsChan)

}
