package work

import (
	"philenius/ethereum-transaction-export/models"
	"sync"

	log "github.com/mgutz/logxi/v1"
	"github.com/onrik/ethrpc"
)

type Pool struct {
	wt         *sync.WaitGroup
	worker     []*Worker
	jobs       chan *models.TxHash
	results    chan *models.Tx
	failedJobs chan *models.TxHash
}

// NewPool instantiates a pool of workers with given concurrency
func NewPool(concurrency int, clientAddr string, jobs chan *models.TxHash, results chan *models.Tx, failedJobs chan *models.TxHash) *Pool {

	wt := &sync.WaitGroup{}

	workerArr := make([]*Worker, concurrency)
	for i := 0; i < concurrency; i++ {
		workerArr[i] = &Worker{ethrpc.NewEthRPC(clientAddr), jobs, results, failedJobs, wt}
	}

	return &Pool{
		wt:         wt,
		worker:     workerArr,
		jobs:       jobs,
		results:    results,
		failedJobs: failedJobs,
	}

}

// Run starts the processing of incoming jobs
func (p *Pool) Run() {
	for _, w := range p.worker {
		go w.doWork()
		p.wt.Add(1)
	}
	p.wt.Wait()

	close(p.results)
	close(p.failedJobs)

	log.Info("all worker completed")
}
