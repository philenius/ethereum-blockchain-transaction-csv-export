package main

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"github.com/mgutz/logxi/v1"
	"github.com/onrik/ethrpc"

	"philenius/ethereum-transaction-export/work"
)

func main() {

	start := time.Now()

	hostname := flag.String("host", "127.0.0.1", "The hostname / IP address of the Ethereum node.")
	port := flag.Int("port", 8545, "The port number of the Ethereum node.")
	startHeight := flag.Int("start", 0, "The height / number of the block where to start.")
	blockCount := flag.Int("count", 1000, "The total amount of blocks to fetch.")
	blockConcurr := flag.Int("blockConcurr", 10, "The count of concurrent workers for fetching block.")
	txConcurr := flag.Int("txConcurr", 10, "The count of concurrent workers for fetching transactions.")
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
	latestBlock := 0
	latestTransactionCount := int64(0)
	wt := sync.WaitGroup{}

	go func() {

		latestLatestBlock := 0
		latestLatestTransactionCount := int64(0)

		for range time.NewTicker(time.Second * 5).C {
			blockDiff := latestBlock - latestLatestBlock
			blockRate := float32(blockDiff) / 5
			txDiff := latestTransactionCount - latestLatestTransactionCount
			txRate := float32(txDiff) / 5
			log.Warn("stats",
				"lastestFetchedBlock", latestBlock,
				"blockRatePerSec", blockRate,
				"fetchedTransactions", latestTransactionCount,
				"txRatePerSec", txRate,
				"txHashChanLen", len(txHashChan),
			)
			latestLatestBlock = latestBlock
			latestLatestTransactionCount = latestTransactionCount
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
		p := work.NewBlockWorkerPool(*blockConcurr, clientAddr, blockHeightChan, txHashChan, failedBlockChan)
		p.Run()
		wt.Done()
	}()
	wt.Add(1)

	// fetch all transactions
	go func() {
		p := work.NewTxWorkerPool(*txConcurr, clientAddr, txHashChan, txChan, failedTxChan)
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
	log.Warn("application successfully completed",
		"durationInMinutes", int(elapsed.Minutes()),
		"startHeight", *startHeight,
		"endHeight", *startHeight+*blockCount,
		"blockCount", *blockCount,
	)

}
