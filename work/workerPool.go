package work

import (
	"sync"

	log "github.com/mgutz/logxi/v1"
	"github.com/onrik/ethrpc"
)

type Pool struct {
	wt         *sync.WaitGroup
	workers    []Worker
	results    chan *Job
	failedJobs chan *Job
}

// NewTxWorkerPool instantiates a pool of workers with given concurrency for fetching transactions.
func NewTxWorkerPool(concurrency int, clientAddr string, jobs chan *Job, results chan *Job,
	failedJobs chan *Job, stats chan *Stat) *Pool {

	wt := &sync.WaitGroup{}

	workerArr := make([]Worker, concurrency)
	for i := 0; i < concurrency; i++ {
		workerArr[i] = &TxWorker{ethrpc.NewEthRPC(clientAddr), jobs, results, failedJobs, stats, wt}
	}

	return &Pool{
		wt:         wt,
		workers:    workerArr,
		results:    results,
		failedJobs: failedJobs,
	}

}

// NewBlockWorkerPool instantiates a pool of workers with given concurrency for fetching blocks.
func NewBlockWorkerPool(concurrency int, clientAddr string, jobs chan *Job, results chan *Job,
	failedJobs chan *Job, stats chan *Stat) *Pool {

	wt := &sync.WaitGroup{}

	workerArr := make([]Worker, concurrency)
	for i := 0; i < concurrency; i++ {
		workerArr[i] = &BlockWorker{ethrpc.NewEthRPC(clientAddr), jobs, results, failedJobs, stats, wt}
	}

	return &Pool{
		wt:         wt,
		workers:    workerArr,
		results:    results,
		failedJobs: failedJobs,
	}

}

// Run starts the processing of incoming jobs
func (p *Pool) Run() {
	for _, w := range p.workers {
		go w.doWork()
		p.wt.Add(1)
	}
	p.wt.Wait()

	close(p.results)
	close(p.failedJobs)

	log.Info("all worker completed")
}
