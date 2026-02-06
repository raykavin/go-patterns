// retry_with_backoff.go
//
// This example demonstrates retry with exponential backoff and jitter.
//
// Key ideas illustrated:
//
//   - Retry loop with max attempts
//   - Exponential backoff with random jitter
//   - Respect context cancellation
//
package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

func sometimesFails() error {
	if rand.Intn(4) != 0 { // succeed ~25%
		return errors.New("transient error")
	}
	return nil
}

func retry(ctx context.Context, max int, base time.Duration, fn func() error) error {
	backoff := base
	for attempt := 1; attempt <= max; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}

		if attempt == max {
			return err
		}

		jitter := time.Duration(rand.Intn(100)) * time.Millisecond
		sleep := backoff + jitter

		select {
		case <-time.After(sleep):
		case <-ctx.Done():
			return ctx.Err()
		}

		if backoff < 800*time.Millisecond {
			backoff *= 2
		}
	}
	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := retry(ctx, 6, 100*time.Millisecond, sometimesFails)
	fmt.Println("final err:", err)
}
