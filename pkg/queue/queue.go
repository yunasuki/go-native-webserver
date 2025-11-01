package queue

import (
	"context"
	"sync"
)

// Job represents a unit of work to be processed by the worker pool.
type Job interface {
	Process(ctx context.Context)
}

// Queue manages job distribution to a pool of workers.
type Queue struct {
	jobs   chan Job
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

// NewQueue creates a new Queue with the given number of workers and buffer size.
func NewQueue(workerCount, bufferSize int) *Queue {
	ctx, cancel := context.WithCancel(context.Background())
	q := &Queue{
		jobs:   make(chan Job, bufferSize),
		ctx:    ctx,
		cancel: cancel,
	}
	for i := 0; i < workerCount; i++ {
		q.wg.Add(1)
		go q.worker()
	}
	return q
}

// Enqueue adds a job to the queue.
func (q *Queue) Enqueue(job Job) {
	select {
	case q.jobs <- job:
	case <-q.ctx.Done():
	}
}

// worker processes jobs from the queue.
func (q *Queue) worker() {
	defer q.wg.Done()
	for {
		select {
		case job := <-q.jobs:
			if job != nil {
				job.Process(q.ctx)
			}
		case <-q.ctx.Done():
			return
		}
	}
}

// Shutdown gracefully stops the queue and waits for all workers to finish.
func (q *Queue) Shutdown() {
	q.cancel()
	q.wg.Wait()
	close(q.jobs)
}
