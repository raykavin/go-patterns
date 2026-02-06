// batching_coalescing.go
//
// This example demonstrates batching/coalescing events to amortize work.
// We flush when:
//   - batch reaches max size, OR
//   - a flush interval elapses
//
// Common uses:
//   - Bulk inserts
//   - Coalescing repeated updates
//   - Reducing per-item overhead
//
package main

import (
	"fmt"
	"time"
)

func batcher(in <-chan int, max int, every time.Duration) <-chan []int {
	out := make(chan []int)
	go func() {
		defer close(out)

		ticker := time.NewTicker(every)
		defer ticker.Stop()

		var buf []int
		flush := func() {
			if len(buf) == 0 {
				return
			}
			b := make([]int, len(buf))
			copy(b, buf)
			out <- b
			buf = buf[:0]
		}

		for {
			select {
			case v, ok := <-in:
				if !ok {
					flush()
					return
				}
				buf = append(buf, v)
				if len(buf) >= max {
					flush()
				}
			case <-ticker.C:
				flush()
			}
		}
	}()
	return out
}

func main() {
	in := make(chan int)
	go func() {
		defer close(in)
		for i := 1; i <= 12; i++ {
			in <- i
			time.Sleep(60 * time.Millisecond)
		}
	}()

	for b := range batcher(in, 4, 200*time.Millisecond) {
		fmt.Println("batch:", b)
	}
}
