// backpressure_handling.go
//
// This example demonstrates backpressure handling using bounded channels.
// When the downstream is slow, upstream blocks naturally, limiting memory growth.
//
// Key ideas illustrated:
//
//   - Buffered channel acts as a queue with a hard cap
//   - Producer will block when queue is full
//
package main

import (
	"fmt"
	"time"
)

func main() {
	queue := make(chan int, 3)

	// slow consumer
	go func() {
		for v := range queue {
			fmt.Println("consume:", v)
			time.Sleep(200 * time.Millisecond)
		}
	}()

	start := time.Now()
	for i := 1; i <= 8; i++ {
		queue <- i // blocks when queue full -> backpressure
		fmt.Println("produce:", i, "t=", time.Since(start).Truncate(10*time.Millisecond))
	}
	close(queue)
	time.Sleep(500 * time.Millisecond)
}
