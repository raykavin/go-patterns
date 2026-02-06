// cqrs_lite.go
//
// This example demonstrates CQRS-lite: separate read and write models.
// It's "lite" because both models are in-process and small.
//
// Key ideas illustrated:
//
//   - Commands mutate state (write side)
//   - Queries read from a read model (read side)
//
package main

import (
	"context"
	"fmt"
	"sync"
)

type Store struct {
	mu    sync.RWMutex
	users map[string]string // id -> name
}

func NewStore() *Store { return &Store{users: map[string]string{}} }

// Command side
func (s *Store) CreateUser(ctx context.Context, id, name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[id] = name
}
func (s *Store) RenameUser(ctx context.Context, id, name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[id] = name
}

// Query side
func (s *Store) GetUserName(ctx context.Context, id string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n, ok := s.users[id]
	return n, ok
}

func main() {
	ctx := context.Background()
	s := NewStore()

	s.CreateUser(ctx, "u1", "Linus")
	s.RenameUser(ctx, "u1", "Linus T.")

	name, ok := s.GetUserName(ctx, "u1")
	fmt.Println("ok?", ok, "name:", name)
}
