// dependency_injection_manual.go
//
// This example demonstrates manual dependency injection in Go (no frameworks).
//
// Key ideas illustrated:
//
//   - Constructors wire dependencies explicitly
//   - Interfaces are used for testability
//
package main

import (
	"context"
	"fmt"
)

type Clock interface {
	NowUnix() int64
}

type RealClock struct{}

func (RealClock) NowUnix() int64 { return 1700000000 } // fixed for demo

type Greeter struct {
	clock Clock
}

func NewGreeter(clock Clock) *Greeter { return &Greeter{clock: clock} }

func (g *Greeter) Greet(ctx context.Context, name string) string {
	return fmt.Sprintf("hello %s (t=%d)", name, g.clock.NowUnix())
}

func main() {
	greeter := NewGreeter(RealClock{})
	fmt.Println(greeter.Greet(context.Background(), "world"))
}
