// drop_pattern_load_shedding.go
//
// This example demonstrates a drop/load-shedding pattern: when downstream is slow,
// we drop new items instead of blocking producers.
//
// Key ideas illustrated:
//
//   - Non-blocking send using select { case ch <- v: default: }
//   - Keeping the system responsive under load
//
package main

import (
	"fmt"
	"time"
)

func main() {
	events := make(chan int, 3) // small buffer simulates limited capacity

	// slow consumer
	go func() {
		for v := range events {
			fmt.Println("consume:", v)
			time.Sleep(200 * time.Millisecond)
		}
	}()

	dropped := 0
	for i := 1; i <= 15; i++ {
		select {
		case events <- i:
		default:
			dropped++
		}
		time.Sleep(40 * time.Millisecond)
	}

	close(events)
	fmt.Println("dropped:", dropped)
}
