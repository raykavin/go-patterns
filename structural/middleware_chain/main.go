// middleware_chain.go
//
// This example demonstrates a middleware chain (common in HTTP, RPC, queues).
//
// Key ideas illustrated:
//
//   - Middleware wraps a handler
//   - Compose middleware in order
//
package main

import (
	"context"
	"fmt"
	"time"
)

type Handler func(ctx context.Context, req string) (string, error)
type Middleware func(Handler) Handler

func Chain(h Handler, mws ...Middleware) Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

func logging() Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, req string) (string, error) {
			start := time.Now()
			resp, err := next(ctx, req)
			fmt.Println("log:", req, "took", time.Since(start).Truncate(time.Millisecond), "err=", err)
			return resp, err
		}
	}
}

func auth() Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, req string) (string, error) {
			if req == "admin" {
				return "", fmt.Errorf("unauthorized")
			}
			return next(ctx, req)
		}
	}
}

func main() {
	base := func(ctx context.Context, req string) (string, error) {
		return "ok:" + req, nil
	}

	h := Chain(base, logging(), auth())
	fmt.Println(h(context.Background(), "user"))
	fmt.Println(h(context.Background(), "admin"))
}
