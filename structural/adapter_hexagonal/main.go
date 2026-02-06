// adapter_ports_and_adapters_hexagonal.go
//
// This example demonstrates a simple Ports & Adapters (Hexagonal) structure.
//
// Key ideas illustrated:
//
//   - "Port" is an interface representing a dependency
//   - "Adapter" implements the port (e.g., HTTP, DB, external API)
//   - Core use case depends only on the port
//
package main

import (
	"context"
	"fmt"
)

type NotifierPort interface {
	Notify(ctx context.Context, userID, message string) error
}

// core use case
type WelcomeUseCase struct {
	notifier NotifierPort
}

func NewWelcomeUseCase(n NotifierPort) *WelcomeUseCase { return &WelcomeUseCase{notifier: n} }

func (uc *WelcomeUseCase) Run(ctx context.Context, userID string) error {
	return uc.notifier.Notify(ctx, userID, "welcome!")
}

// adapter
type ConsoleNotifier struct{}

func (ConsoleNotifier) Notify(ctx context.Context, userID, message string) error {
	fmt.Println("notify:", userID, message)
	return nil
}

func main() {
	uc := NewWelcomeUseCase(ConsoleNotifier{})
	_ = uc.Run(context.Background(), "u123")
}
