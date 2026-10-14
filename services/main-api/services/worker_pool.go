package services

import (
	"context"
	"log"
	"sync"
)

// PricingJob represents a single job for pricing calculation
type PricingJob struct {
	ID      int
	Request PricingRequest
	Result  chan PricingResponse
	Err     chan error
}

// PricingWorkerPool manages a pool of goroutines to process pricing calculations concurrently
type PricingWorkerPool struct {
	jobs    chan PricingJob
	wg      sync.WaitGroup
	workers int
}

// Global instance of the worker pool
var GlobalPricingPool *PricingWorkerPool

// InitPricingWorkerPool initializes the pool with a fixed number of workers
func InitPricingWorkerPool(numWorkers int) {
	GlobalPricingPool = &PricingWorkerPool{
		jobs:    make(chan PricingJob, 1000), // Buffer size 1000
		workers: numWorkers,
	}

	for i := 0; i < numWorkers; i++ {
		GlobalPricingPool.wg.Add(1)
		go GlobalPricingPool.worker(i)
	}
}

// worker processes incoming jobs
func (p *PricingWorkerPool) worker(id int) {
	defer p.wg.Done()
	for job := range p.jobs {
		// log.Printf("Worker %d starting job %d", id, job.ID)
		
		resp, err := CalculateQuote(job.Request)
		
		if err != nil {
			job.Err <- err
		} else {
			job.Result <- resp
		}
	}
}

// SubmitJob sends a job to the worker pool and returns channels for the result/error
func (p *PricingWorkerPool) SubmitJob(ctx context.Context, req PricingRequest) (<-chan PricingResponse, <-chan error) {
	resChan := make(chan PricingResponse, 1)
	errChan := make(chan error, 1)

	job := PricingJob{
		ID:      1, // Could use atomic counter or UUID
		Request: req,
		Result:  resChan,
		Err:     errChan,
	}

	select {
	case p.jobs <- job:
		// Successfully queued
	case <-ctx.Done():
		errChan <- ctx.Err()
	}

	return resChan, errChan
}

// Shutdown gracefully shuts down the worker pool
func (p *PricingWorkerPool) Shutdown() {
	close(p.jobs)
	p.wg.Wait()
	log.Println("Pricing Worker Pool shut down gracefully")
}
