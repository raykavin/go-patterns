// tee_channel.go
//
// This example demonstrates a tee channel pattern that duplicates a stream into two outputs.
//
// Key ideas illustrated:
//
//   - A single input is broadcast to two outputs
//   - Both outputs receive all items in order
//   - Stops cleanly when input is closed
//
package main

import "fmt"

func tee[T any](in <-chan T) (<-chan T, <-chan T) {
	out1 := make(chan T)
	out2 := make(chan T)

	go func() {
		defer close(out1)
		defer close(out2)
		for v := range in {
			// ensure both get v
			out1 <- v
			out2 <- v
		}
	}()

	return out1, out2
}

func main() {
	in := make(chan int)
	go func() {
		defer close(in)
		for i := 1; i <= 5; i++ {
			in <- i
		}
	}()

	a, b := tee(in)
	for i := 0; i < 5; i++ {
		fmt.Println(<-a, <-b)
	}
}
