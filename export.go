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

	for transaction := range txChan {
		tx := transaction.tx
		w.WriteString(
			fmt.Sprintf(
				"%s,%d,%s,%d,%d,%s,%s,%d,%d,%d,%s,%d\n",
				tx.Hash, tx.Nonce, tx.BlockHash, tx.BlockNumber, tx.TransactionIndex, tx.From, tx.To, tx.Value.Int64(), tx.Gas, tx.GasPrice.Int64(), tx.Input, transaction.timestamp,
			),
		)
		w.Flush()
	}

	log.Info("finished writing to file", "file", f.Name())
	w.Flush()
}
