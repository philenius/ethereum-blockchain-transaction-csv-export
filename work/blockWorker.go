package work

import (
	"sync"

	log "github.com/mgutz/logxi/v1"
	"github.com/onrik/ethrpc"
)

type BlockWorker struct {
	client     *ethrpc.EthRPC
	jobs       chan *Job
	result     chan *Job
	failedJobs chan *Job
	wt         *sync.WaitGroup
}

func (w *BlockWorker) doWork() {
	for job := range w.jobs {
		//latestBlock = blockHeight
		block, err := w.client.EthGetBlockByNumber(job.BlockHeight, true)
		if err != nil {
			w.failedJobs <- job
			log.Error("failed to get block", "blockNumber", job.BlockHeight, "err", err.Error())
			continue
		}
		if log.IsDebug() {
			log.Debug("successfully got block", "blockNumber", block.Number)
		}
		for _, tx := range block.Transactions {
			w.result <- &Job{TxHash: tx.Hash, Timestamp: block.Timestamp}
		}
	}
	log.Info("block worker finished")
	w.wt.Done()
}
