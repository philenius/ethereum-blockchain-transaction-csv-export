package work

import (
	"sync"

	log "github.com/mgutz/logxi/v1"
	"github.com/onrik/ethrpc"
)

type TxWorker struct {
	client     *ethrpc.EthRPC
	jobs       chan *Job
	result     chan *Job
	failedJobs chan *Job
	stats      chan *Stat
	wt         *sync.WaitGroup
}

func (w *TxWorker) doWork() {
	for job := range w.jobs {
		w.stats <- &Stat{FetchedTransaction: true}
		transaction, err := w.client.EthGetTransactionByHash(job.TxHash)
		if err != nil || transaction.BlockNumber == nil || transaction.TransactionIndex == nil {
			w.failedJobs <- job
			log.Error("failed to get transaction", "txHash", job.TxHash)
			continue
		}

		w.result <- &Job{
			Timestamp: job.Timestamp,
			Tx:        transaction,
		}
		if log.IsDebug() {
			log.Debug("successfully got transaction", "txHash", transaction.Hash)
		}
	}
	log.Info("transaction worker finished")
	w.wt.Done()
}
