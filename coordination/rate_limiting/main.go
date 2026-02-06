// rate_limiting.go
//
// This example demonstrates a simple token-bucket rate limiter using time.Ticker.
//
// Key ideas illustrated:
//
//   - A ticker refills tokens at a fixed rate
//   - Callers must acquire a token before proceeding
//
package main

import (
	"context"
	"fmt"
	"time"
)

type RateLimiter struct {
	tokens chan struct{}
	stop   chan struct{}
}

func NewRateLimiter(rps int, burst int) *RateLimiter {
	rl := &RateLimiter{
		tokens: make(chan struct{}, burst),
		stop:   make(chan struct{}),
	}

	// fill burst initially
	for i := 0; i < burst; i++ {
		rl.tokens <- struct{}{}
	}

	interval := time.Second / time.Duration(rps)
	t := time.NewTicker(interval)
	go func() {
		defer t.Stop()
		for {
			select {
			case <-rl.stop:
				return
			case <-t.C:
				select {
				case rl.tokens <- struct{}{}:
				default:
				}
			}
		}
	}()
	return rl
}

func (rl *RateLimiter) Stop() { close(rl.stop) }

func (rl *RateLimiter) Acquire(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-rl.tokens:
		return nil
	}
}

func main() {
	rl := NewRateLimiter(5, 2) // 5 req/s, burst 2
	defer rl.Stop()

	ctx := context.Background()
	start := time.Now()
	for i := 1; i <= 10; i++ {
		_ = rl.Acquire(ctx)
		fmt.Println("request", i, "at", time.Since(start).Truncate(10*time.Millisecond))
	}
}
