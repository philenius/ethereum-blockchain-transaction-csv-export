package main

import (
	"bufio"
	"fmt"
	"os"
	"philenius/ethereum-transaction-export/work"
	"time"

	"github.com/mgutz/logxi/v1"
)

func exportAsCSV(jobs chan *work.Job) {

	now := time.Now().Format("2006-01-02-15-04-05")
	f, err := os.Create(fmt.Sprintf("geth_tx_export_%s.csv", now))
	if err != nil {
		log.Fatal("failed to open output file", "err", err.Error())
	}

	w := bufio.NewWriter(f)

	w.WriteString("tx_hash,tx_nonce,tx_block_hash,tx_block_number,tx_index,tx_from,tx_to,tx_value,tx_gas,tx_gas_price,tx_input,tx_timestamp\n")
	w.Flush()

	lineBuf := 0
	for job := range jobs {
		tx := job.Tx

		lineBuf++
		line := fmt.Sprintf(
			"%s,%d,%s,%d,%d,%s,%s,%d,%d,%d,%s,%d\n",
			tx.Hash, tx.Nonce, tx.BlockHash, *tx.BlockNumber, *tx.TransactionIndex, tx.From, tx.To,
			tx.Value.Int64(), tx.Gas, tx.GasPrice.Int64(), tx.Input, job.Timestamp,
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

func exportFailedBlockJobs(jobs chan *work.Job) {

	f, err := os.Create("failedBlocks.txt")
	if err != nil {
		log.Fatal("failed to create block error file", "err", err.Error())
	}

	w := bufio.NewWriter(f)

	for job := range jobs {
		w.WriteString(fmt.Sprintf("%d\n", job.BlockHeight))
		w.Flush()
	}

	log.Info("finished writing to block error file")
	w.Flush()
}

func exportFailedTxJobs(jobs chan *work.Job) {

	f, err := os.Create("failedTransactions.txt")
	if err != nil {
		log.Fatal("failed to create transaction error file", "err", err.Error())
	}

	w := bufio.NewWriter(f)

	for job := range jobs {
		w.WriteString(fmt.Sprintf("%s\n", job.TxHash))
		w.Flush()
	}

	log.Info("finished writing to transaction error file")
	w.Flush()
}
