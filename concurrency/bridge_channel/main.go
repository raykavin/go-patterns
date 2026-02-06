// bridge_channel_dynamic_channel_of_channels.go
//
// This example demonstrates the "bridge" pattern: flattening a channel-of-channels
// into a single output stream. The producer can dynamically create sub-streams.
//
// Key ideas illustrated:
//
//   - out <-chan <-chan T as dynamic streams
//   - Bridge consumes each inner channel sequentially
//   - Uses a done channel to support cancellation
//
package main

import (
	"fmt"
	"time"
)

func bridge[T any](done <-chan struct{}, chans <-chan (<-chan T)) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for {
			var ch <-chan T
			select {
			case <-done:
				return
			case c, ok := <-chans:
				if !ok {
					return
				}
				ch = c
			}

			for v := range ch {
				select {
				case <-done:
					return
				case out <- v:
				}
			}
		}
	}()
	return out
}

func main() {
	done := make(chan struct{})
	defer close(done)

	streams := make(chan (<-chan int))

	go func() {
		defer close(streams)
		for s := 1; s <= 3; s++ {
			s := s
			ch := make(chan int)
			go func() {
				defer close(ch)
				for i := 1; i <= 3; i++ {
					ch <- s*10 + i
					time.Sleep(40 * time.Millisecond)
				}
			}()
			streams <- ch
		}
	}()

	for v := range bridge(done, streams) {
		fmt.Println(v)
	}
}
