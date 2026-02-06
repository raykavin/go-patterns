// circuit_breaker.go
//
// This example demonstrates a small circuit breaker to protect a failing dependency.
//
// Key ideas illustrated:
//
//   - CLOSED -> OPEN after failure threshold
//   - OPEN -> HALF-OPEN after cool-down
//   - HALF-OPEN -> CLOSED on success (or back to OPEN on failure)
//
package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type State int

const (
	Closed State = iota
	Open
	HalfOpen
)

type CircuitBreaker struct {
	mu sync.Mutex

	state          State
	failures       int
	failureThresh  int
	openUntil      time.Time
	coolDown       time.Duration
	halfOpenTrials int
}

func NewCircuitBreaker(thresh int, coolDown time.Duration, halfOpenTrials int) *CircuitBreaker {
	return &CircuitBreaker{
		state:          Closed,
		failureThresh:  thresh,
		coolDown:       coolDown,
		halfOpenTrials: halfOpenTrials,
	}
}

var ErrOpen = errors.New("circuit breaker is open")

func (cb *CircuitBreaker) Allow() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := time.Now()
	switch cb.state {
	case Open:
		if now.After(cb.openUntil) {
			cb.state = HalfOpen
			cb.failures = 0
			return nil
		}
		return ErrOpen
	default:
		return nil
	}
}

func (cb *CircuitBreaker) OnSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures = 0
	cb.state = Closed
}

func (cb *CircuitBreaker) OnFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures++
	if cb.state == HalfOpen && cb.failures >= 1 {
		cb.trip()
		return
	}
	if cb.state == Closed && cb.failures >= cb.failureThresh {
		cb.trip()
	}
}

func (cb *CircuitBreaker) trip() {
	cb.state = Open
	cb.openUntil = time.Now().Add(cb.coolDown)
}

func (cb *CircuitBreaker) State() State {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.state
}

func main() {
	cb := NewCircuitBreaker(3, 400*time.Millisecond, 1)

	// dependency: fail first 5 calls, then succeed
	failUntil := 5
	call := 0
	dep := func() error {
		call++
		if call <= failUntil {
			return errors.New("dependency down")
		}
		return nil
	}

	for i := 1; i <= 12; i++ {
		if err := cb.Allow(); err != nil {
			fmt.Println(i, "blocked:", err, "state=", cb.State())
			time.Sleep(120 * time.Millisecond)
			continue
		}

		err := dep()
		if err != nil {
			fmt.Println(i, "call failed:", err, "state=", cb.State())
			cb.OnFailure()
		} else {
			fmt.Println(i, "call ok", "state=", cb.State())
			cb.OnSuccess()
		}
		time.Sleep(120 * time.Millisecond)
	}
}
