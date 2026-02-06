// select_based_multiplexing.go
//
// This example demonstrates select-based multiplexing: consuming from multiple
// input channels and producing a single output stream.
//
// Key ideas illustrated:
//
//   - select to receive from multiple sources
//   - graceful handling when inputs close
//
package main

import (
	"fmt"
	"time"
)

func source(prefix string, n int, every time.Duration) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for i := 1; i <= n; i++ {
			out <- fmt.Sprintf("%s-%d", prefix, i)
			time.Sleep(every)
		}
	}()
	return out
}

func multiplex(a, b <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for a != nil || b != nil {
			select {
			case v, ok := <-a:
				if !ok {
					a = nil
					continue
				}
				out <- v
			case v, ok := <-b:
				if !ok {
					b = nil
					continue
				}
				out <- v
			}
		}
	}()
	return out
}

func main() {
	a := source("A", 4, 70*time.Millisecond)
	b := source("B", 4, 110*time.Millisecond)

	for v := range multiplex(a, b) {
		fmt.Println(v)
	}
}
