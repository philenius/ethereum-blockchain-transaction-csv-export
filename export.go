package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/mgutz/logxi/v1"
)

func exportAsCSV(txChan chan *Tx) {

	now := time.Now().Format("2006-01-02-15-04-05")
	f, err := os.Create(fmt.Sprintf("geth_tx_export_%s.csv", now))
	if err != nil {
		log.Fatal("failed to open output file", "err", err.Error())
	}

	w := bufio.NewWriter(f)

	w.WriteString("tx_hash,tx_nonce,tx_block_hash,tx_block_number,tx_index,tx_from,tx_to,tx_value,tx_gas,tx_gas_price,tx_input,tx_timestamp\n")
	w.Flush()

	lineBuf := 0
	for transaction := range txChan {
		tx := transaction.tx

		lineBuf++
		line := fmt.Sprintf(
			"%s,%d,%s,%d,%d,%s,%s,%d,%d,%d,%s,%d\n",
			tx.Hash, tx.Nonce, tx.BlockHash, *tx.BlockNumber, *tx.TransactionIndex, tx.From, tx.To, tx.Value.Int64(), tx.Gas, tx.GasPrice.Int64(), tx.Input, transaction.timestamp,
		)
		w.WriteString(line)

		if lineBuf > 50 {
			lineBuf = 0
			w.Flush()
		}
	}
	log.Info("finished writing to file", "file", f.Name())
	w.Flush()
}

func exportFailedBlocks(blockChan chan int) {

	f, err := os.Create("failedBlocks.txt")
	if err != nil {
		log.Fatal("failed to create block error file", "err", err.Error())
	}

	w := bufio.NewWriter(f)

	for block := range blockChan {
		w.WriteString(fmt.Sprintf("%d\n", block))
		w.Flush()
	}

	log.Info("finished writing to block error file")
	w.Flush()
}

func exportFailedTx(txChan chan string) {

	f, err := os.Create("failedTransactions.txt")
	if err != nil {
		log.Fatal("failed to create transaction error file", "err", err.Error())
	}

	w := bufio.NewWriter(f)

	for tx := range txChan {
		fmt.Println(tx)
		w.WriteString(fmt.Sprintf("%s\n", tx))
		w.Flush()
	}

	log.Info("finished writing to transaction error file")
	w.Flush()
}
