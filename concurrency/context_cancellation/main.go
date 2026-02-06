// context_cancellation_propagation.go
//
// This example demonstrates cooperative cancellation propagation via context.Context.
//
// Key ideas illustrated:
//
//   - A parent context canceled on first error
//   - Workers observe ctx.Done() to avoid goroutine leaks
//   - Safe fan-in collection that respects cancellation
//
package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type item struct{ id int }

func worker(ctx context.Context, in <-chan item, out chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case it, ok := <-in:
			if !ok {
				return
			}
			// simulate work
			time.Sleep(80 * time.Millisecond)
			select {
			case out <- it.id * 10:
			case <-ctx.Done():
				return
			}
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	in := make(chan item)
	out := make(chan int)

	const workers = 3
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go worker(ctx, in, out, &wg)
	}

	go func() {
		defer close(in)
		for i := 1; i <= 10; i++ {
			in <- item{id: i}
		}
	}()

	// close out when workers are done
	go func() {
		wg.Wait()
		close(out)
	}()

	// cancel on a condition (simulate an error)
	errCh := make(chan error, 1)
	go func() {
		time.Sleep(350 * time.Millisecond)
		errCh <- errors.New("simulated failure: cancel everything")
	}()

	for {
		select {
		case v, ok := <-out:
			if !ok {
				fmt.Println("all done")
				return
			}
			fmt.Println("result:", v)
		case err := <-errCh:
			fmt.Println("error:", err)
			cancel()
		}
	}
}
