// graceful_shutdown.go
//
// This example demonstrates graceful shutdown for a long-running server-like loop.
//
// Key ideas illustrated:
//
//   - os.Signal notification
//   - context cancellation to stop background goroutines
//   - WaitGroup to wait for in-flight work before exiting
//
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				fmt.Println("worker: stop")
				return
			case <-ticker.C:
				fmt.Println("worker: tick")
			}
		}
	}()

	fmt.Println("Press Ctrl+C to shutdown...")
	<-sigCh
	fmt.Println("signal received, shutting down...")
	cancel()
	wg.Wait()
	fmt.Println("clean exit")
}
