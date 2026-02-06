// or_channel_combine_cancellation.go
//
// This example demonstrates an "or-channel" that closes when any of the input
// done channels closes (fan-in of cancellation signals).
//
// Key ideas illustrated:
//
//   - Combine multiple cancellation sources
//   - Useful for waiting on "any" completion condition
//
package main

import (
	"fmt"
	"time"
)

func or(dones ...<-chan struct{}) <-chan struct{} {
	switch len(dones) {
	case 0:
		ch := make(chan struct{})
		close(ch)
		return ch
	case 1:
		return dones[0]
	}

	out := make(chan struct{})
	go func() {
		defer close(out)
		select {
		case <-dones[0]:
		case <-dones[1]:
		case <-or(dones[2:]...):
		}
	}()
	return out
}

func main() {
	a := make(chan struct{})
	b := make(chan struct{})
	c := make(chan struct{})

	go func() { time.Sleep(150 * time.Millisecond); close(b) }()
	go func() { time.Sleep(300 * time.Millisecond); close(a) }()
	go func() { time.Sleep(450 * time.Millisecond); close(c) }()

	start := time.Now()
	<-or(a, b, c)
	fmt.Println("first done after:", time.Since(start))
}
