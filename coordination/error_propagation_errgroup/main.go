// error_propagation_errgroup.go
//
// This example demonstrates error propagation and cancellation using an errgroup-like
// helper (implemented locally to keep the example dependency-free).
//
// Key ideas illustrated:
//
//   - First error cancels sibling goroutines (via derived context)
//   - Wait returns that first error
//
package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type Group struct {
	wg     sync.WaitGroup
	cancel context.CancelFunc

	errOnce sync.Once
	err     error
}

func WithContext(parent context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(parent)
	return &Group{cancel: cancel}, ctx
}

func (g *Group) Go(fn func() error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		if err := fn(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				g.cancel()
			})
		}
	}()
}

func (g *Group) Wait() error {
	g.wg.Wait()
	g.cancel()
	return g.err
}

func main() {
	ctx := context.Background()
	g, ctx := WithContext(ctx)

	g.Go(func() error {
		select {
		case <-time.After(150 * time.Millisecond):
			fmt.Println("task1 ok")
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	g.Go(func() error {
		select {
		case <-time.After(250 * time.Millisecond):
			return errors.New("task2 failed")
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("task3 canceled:", ctx.Err())
				return ctx.Err()
			case <-time.After(80 * time.Millisecond):
				fmt.Println("task3 heartbeat")
			}
		}
	})

	if err := g.Wait(); err != nil {
		fmt.Println("group error:", err)
	}
}
