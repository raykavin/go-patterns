// event_dispatcher.go
//
// This example demonstrates an event dispatcher (in-process event bus).
//
// Key ideas illustrated:
//
//   - Register handlers for event types
//   - Dispatch events to all handlers
//
package main

import (
	"fmt"
	"sync"
)

type Event interface {
	Name() string
}

type UserCreated struct{ ID string }

func (UserCreated) Name() string { return "UserCreated" }

type Handler func(Event)

type Dispatcher struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{handlers: map[string][]Handler{}}
}

func (d *Dispatcher) On(eventName string, h Handler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers[eventName] = append(d.handlers[eventName], h)
}

func (d *Dispatcher) Dispatch(e Event) {
	d.mu.RLock()
	hs := append([]Handler(nil), d.handlers[e.Name()]...)
	d.mu.RUnlock()

	for _, h := range hs {
		h(e)
	}
}

func main() {
	d := NewDispatcher()

	d.On("UserCreated", func(e Event) {
		uc := e.(UserCreated)
		fmt.Println("send email to", uc.ID)
	})
	d.On("UserCreated", func(e Event) {
		uc := e.(UserCreated)
		fmt.Println("audit log:", uc.ID)
	})

	d.Dispatch(UserCreated{ID: "u7"})
}
