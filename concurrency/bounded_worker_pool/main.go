// bounded_worker_pool.go
//
// This example demonstrates concurrency limiting by bounding in-flight work.
// We accept arbitrary submissions, but only allow N tasks to run concurrently.
//
// Key ideas illustrated:
//
//   - A buffered "tokens" channel acts as a semaphore
//   - Each task acquires a token before running and releases it after
//
package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func runBounded(ctx context.Context, limit int, tasks []func(context.Context) error) error {
	tokens := make(chan struct{}, limit)
	var wg sync.WaitGroup

	errCh := make(chan error, 1)

	for _, t := range tasks {
		wg.Add(1)
		go func(task func(context.Context) error) {
			defer wg.Done()

			select {
			case tokens <- struct{}{}: // acquire
			case <-ctx.Done():
				return
			}
			defer func() { <-tokens }() // release

			if err := task(ctx); err != nil {
				select {
				case errCh <- err:
				default:
				}
			}
		}(t)
	}

	wg.Wait()
	close(errCh)

	return <-errCh
}

func main() {
	ctx := context.Background()

	tasks := make([]func(context.Context) error, 0, 8)
	for i := 1; i <= 8; i++ {
		i := i
		tasks = append(tasks, func(ctx context.Context) error {
			fmt.Println("start", i)
			time.Sleep(200 * time.Millisecond)
			fmt.Println("done ", i)
			return nil
		})
	}

	if err := runBounded(ctx, 3, tasks); err != nil {
		fmt.Println("error:", err)
	}
}
