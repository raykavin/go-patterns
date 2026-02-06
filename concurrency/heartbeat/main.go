// heartbeat_liveness_signaling.go
//
// This example demonstrates a heartbeat pattern where a worker periodically emits
// liveness signals while doing work.
//
// Key ideas illustrated:
//
//   - A heartbeat channel indicates "still alive"
//   - A results channel carries actual outputs
//   - The caller can detect stalled workers via timeouts
//
package main

import (
	"context"
	"fmt"
	"time"
)

func doWork(ctx context.Context) (<-chan struct{}, <-chan int) {
	heartbeat := make(chan struct{}, 1)
	results := make(chan int)

	go func() {
		defer close(heartbeat)
		defer close(results)

		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()

		for i := 1; i <= 5; i++ {
			// simulate work chunk
			select {
			case <-time.After(180 * time.Millisecond):
				results <- i
			case <-ctx.Done():
				return
			}

			// best-effort heartbeat
			select {
			case <-ticker.C:
				select {
				case heartbeat <- struct{}{}:
				default:
				}
			default:
			}
		}
	}()

	return heartbeat, results
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	heartbeat, results := doWork(ctx)

	for {
		select {
		case _, ok := <-heartbeat:
			if ok {
				fmt.Println("heartbeat")
			}
		case v, ok := <-results:
			if !ok {
				fmt.Println("done")
				return
			}
			fmt.Println("result:", v)
		case <-time.After(250 * time.Millisecond):
			fmt.Println("no heartbeat recently -> cancel")
			cancel()
		}
	}
}
