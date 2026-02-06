// supervisor_restart_strategy.go
//
// This example demonstrates a simple supervisor that restarts a worker when it fails.
//
// Key ideas illustrated:
//
//   - A supervisor loop runs the child worker
//   - On error, it waits (backoff) and restarts
//   - Stops when parent context is canceled
//
package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

func flakyWorker(ctx context.Context) error {
	// do some "work"
	select {
	case <-time.After(120 * time.Millisecond):
	case <-ctx.Done():
		return ctx.Err()
	}

	// randomly fail
	if rand.Intn(3) == 0 {
		return errors.New("worker crashed")
	}
	fmt.Println("worker: completed one cycle")
	return nil
}

func supervise(ctx context.Context) {
	backoff := 100 * time.Millisecond
	for {
		select {
		case <-ctx.Done():
			fmt.Println("supervisor: stop")
			return
		default:
		}

		err := flakyWorker(ctx)
		if err == nil {
			backoff = 100 * time.Millisecond
			continue
		}
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			fmt.Println("supervisor: worker ended:", err)
			return
		}

		fmt.Println("supervisor: restart after error:", err)
		time.Sleep(backoff)
		if backoff < 800*time.Millisecond {
			backoff *= 2
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	ctx, cancel := context.WithTimeout(context.Background(), 1200*time.Millisecond)
	defer cancel()

	supervise(ctx)
}
