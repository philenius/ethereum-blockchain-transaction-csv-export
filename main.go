package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/mgutz/logxi/v1"
	"github.com/onrik/ethrpc"
)

func main() {

	hostname := flag.String("host", "127.0.0.1", "The hostname / IP address of the Ethereum node.")
	port := flag.Int("port", 8545, "The port number of the Ethereum node.")
	startHeight := flag.Int("start", 0, "The height / number of the block where to start.")
	blockCount := flag.Int("count", 1000, "The total amount of blocks to fetch.")
	flag.Parse()

	client := ethrpc.New(fmt.Sprintf("http://%s:%d", *hostname, *port))
	version, err := client.Web3ClientVersion()
	if err != nil {
		log.Fatal("failed to connect to Ethereum node", "err", err.Error())
	}
	log.Info("connected to Ethereum node", "version", version)

	blockHeightChan := make(chan int, 10000)
	txHashChan := make(chan *TxHash, 10000)
	txChan := make(chan *Tx, 10000)
	wt := sync.WaitGroup{}

	// list all blocks to fetch
	go func() {
		endHeight := *startHeight + *blockCount
		for i := *startHeight; i < endHeight; i++ {
			blockHeightChan <- i
		}
		log.Info("finished listing block numbers")
		close(blockHeightChan)
		wt.Done()
	}()
	wt.Add(1)

	// fetch all blocks
	go func() {
		for blockHeight := range blockHeightChan {
			block, err := client.EthGetBlockByNumber(blockHeight, true)
			if err != nil {
				log.Error("failed to get block", "blockNumber", *blockCount, "err", err.Error())
				continue
			}
			log.Debug("successfully got block", "blockNumber", block.Number)
			for _, tx := range block.Transactions {
				txHashChan <- &TxHash{block.Timestamp, tx.Hash}
			}
		}
		log.Info("finished fetching all blocks")
		close(txHashChan)
		wt.Done()
	}()
	wt.Add(1)

	// fetch all transactions
	go func() {
		for txHash := range txHashChan {
			transaction, err := client.EthGetTransactionByHash(txHash.hash)
			if err != nil {
				log.Error("failed to get transaction", "txiHash", txHash)
				continue
			}
			txChan <- &Tx{txHash.timestamp, transaction}
			log.Debug("successfully got transaction", "txHash", transaction.Hash)
		}
		close(txChan)
		log.Info("finished fetching all transactions")
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
	log.Info("application completed")

}
