package work

import (
	"philenius/ethereum-transaction-export/models"
	"sync"

	log "github.com/mgutz/logxi/v1"
	"github.com/onrik/ethrpc"
)

type Worker struct {
	client     *ethrpc.EthRPC
	jobs       chan *models.TxHash
	result     chan *models.Tx
	failedJobs chan *models.TxHash
	wt         *sync.WaitGroup
}

func (w *Worker) doWork() {
	for txHash := range w.jobs {
		//latestTransactionCount++
		transaction, err := w.client.EthGetTransactionByHash(txHash.Hash)
		if err != nil || transaction.BlockNumber == nil || transaction.TransactionIndex == nil {
			w.failedJobs <- txHash
			log.Error("failed to get transaction", "txHash", txHash.Hash)
			continue
		}

		w.result <- &models.Tx{
			Timestamp: txHash.Timestamp,
			Tx:        transaction,
		}
		if log.IsDebug() {
			log.Debug("successfully got transaction", "txHash", transaction.Hash)
		}
	}
	log.Info("finished fetching all transactions")
	w.wt.Done()
}
