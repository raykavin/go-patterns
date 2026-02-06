// semaphore_channel_throttling.go
//
// This example shows a classic channel-based semaphore for throttling concurrent work.
//
// Key ideas illustrated:
//
//   - Buffered channel as token bucket
//   - Defer release to avoid leaks
//
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	const maxConcurrent = 2
	sem := make(chan struct{}, maxConcurrent)

	var wg sync.WaitGroup
	for i := 1; i <= 6; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			sem <- struct{}{}        // acquire
			defer func() { <-sem }() // release

			fmt.Println("working", i)
			time.Sleep(150 * time.Millisecond)
			fmt.Println("done   ", i)
		}()
	}

	wg.Wait()
}
