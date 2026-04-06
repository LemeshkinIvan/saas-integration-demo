package worker

import (
	custom_logger "daos_core/internal/utils/logger"
	"fmt"
)

const jobsCount = 10

var Pool *WorkerPool

func InitGlobalPool() {
	Pool = NewWorkerPool(5, 1000)
	Pool.Start()
}

type Job interface {
	Do() error
}

type WorkerPool struct {
	jobs    chan Job
	workers int
}

func NewWorkerPool(workers, queueSize int) *WorkerPool {
	return &WorkerPool{
		jobs:    make(chan Job, queueSize),
		workers: workers,
	}
}

// запуск пула
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workers; i++ {
		go func() {
			for job := range wp.jobs {
				if err := job.Do(); err != nil {
					custom_logger.AsyncLog(1, fmt.Sprintf("job error: %v", err))
				}
			}
		}()
	}
}

// добавление задачи в очередь
func (wp *WorkerPool) Enqueue(job Job) {
	wp.jobs <- job
}

// остановка пула
func (wp *WorkerPool) Stop() {
	close(wp.jobs)
}
