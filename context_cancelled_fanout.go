// fanout_context_cancellation.go
//
// This example demonstrates an idiomatic Go concurrency pattern where a
// parent operation spawns multiple goroutines (fan-out) to perform work
// in parallel while supporting cooperative cancellation using context.Context.
//
// Key ideas illustrated:
//
//   - Propagating cancellation signals through context
//   - Preventing goroutine leaks by observing ctx.Done()
//   - Coordinating completion using sync.WaitGroup
//   - Safely closing channels after all workers finish
//   - Avoiding blocked sends when cancellation occurs
//   - Collecting results while respecting timeouts
//
// This pattern is commonly used in real-world systems for:
//   - Parallel HTTP/API requests
//   - Database shard queries
//   - Scatter/gather workloads
//   - Aggregation services
//
// NOTE:
// This is NOT a worker pool implementation. Workers are short-lived and
// scoped to a single operation. For persistent job processing, a worker
// pool pattern should be used instead.
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res, err := largeOp(ctx)
	fmt.Println("result:", res, "error:", err)
}

func largeOp(ctx context.Context) ([]string, error) {
	const workers = 10

	results := make(chan string)
	var wg sync.WaitGroup

	wg.Add(workers)

	// Simulate parallels externals calls
	for i := 0; i < workers; i++ {
		go func(id int) {
			defer wg.Done()

			// Simulate a large op
			select {
			case <-time.After(time.Duration(1+id) * time.Second):

				// Before send the result, checks the cancellation
				select {
				case results <- fmt.Sprintf("worker %d finished", id):
				case <-ctx.Done():
					return
				}

			case <-ctx.Done():
				// Cancelled during the work
				return
			}
		}(i)
	}

	// Close the channel after all workers completed
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect the results
	var collected []string

	for {
		select {
		case r, ok := <-results:
			if !ok {
				return collected, nil
			}
			collected = append(collected, r)

		case <-ctx.Done():
			// External cancellation
			return collected, ctx.Err()
		}
	}
}
