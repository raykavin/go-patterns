// timeout_pattern.go
//
// This example demonstrates a timeout pattern using context.WithTimeout to bound
// the duration of an operation.
//
// Key ideas illustrated:
//
//   - Deriving a timeout context
//   - Select on ctx.Done() to stop waiting
//
package main

import (
	"context"
	"fmt"
	"time"
)

func slowOperation(ctx context.Context) (string, error) {
	select {
	case <-time.After(800 * time.Millisecond):
		return "completed", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	res, err := slowOperation(ctx)
	fmt.Println("res:", res, "err:", err)
}
