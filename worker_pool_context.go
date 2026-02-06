// worker_pool_context.go
//
// This example demonstrates an idiomatic worker pool pattern in Go using
// context-based cancellation and cooperative shutdown.
//
// Key concepts illustrated:
//
//   - Fixed number of persistent workers consuming jobs
//   - Backpressure via job channels
//   - Context propagation for cancellation
//   - Preventing blocked sends on shutdown
//   - Coordinated worker lifecycle using sync.WaitGroup
//   - Graceful channel closing
//
// This pattern is commonly used for:
//   - Background processing pipelines
//   - Job queues
//   - Batch workloads
//   - Limiting concurrency under load
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID int
}

type Result struct {
	JobID int
	Value int
}

// Worker consumes jobs until the context is cancelled or the job channel closes
func worker(
	ctx context.Context,
	id int,
	jobs <-chan Job,
	results chan<- Result,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				fmt.Println("worker", id, "shutting down (jobs closed)")
				return
			}

			// Simulate work
			time.Sleep(500 * time.Millisecond)

			res := Result{
				JobID: job.ID,
				Value: job.ID * 2,
			}

			// Send result cooperatively
			select {
			case results <- res:
			case <-ctx.Done():
				fmt.Println("worker", id, "cancelled")
				return
			}

		case <-ctx.Done():
			fmt.Println("worker", id, "cancelled")
			return
		}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	const workerCount = 4
	const jobCount = 20

	jobs := make(chan Job)
	results := make(chan Result)

	var wg sync.WaitGroup

	// Start worker pool
	wg.Add(workerCount)
	for w := 0; w < workerCount; w++ {
		go worker(ctx, w, jobs, results, &wg)
	}

	// Produce jobs
	go func() {
		defer close(jobs)
		for i := 0; i < jobCount; i++ {
			select {
			case jobs <- Job{ID: i}:
			case <-ctx.Done():
				return
			}
		}
	}()

	// Close results after workers exit
	go func() {
		wg.Wait()
		close(results)
	}()

	// Consume results
	for {
		select {
		case r, ok := <-results:
			if !ok {
				fmt.Println("all workers finished")
				return
			}
			fmt.Printf("result: job=%d value=%d\n", r.JobID, r.Value)

		case <-ctx.Done():
			fmt.Println("main cancelled:", ctx.Err())
			return
		}
	}
}
