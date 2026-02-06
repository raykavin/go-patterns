// repository_pattern.go
//
// This example demonstrates the Repository pattern: an interface that abstracts
// data persistence from business logic.
//
// Key ideas illustrated:
//
//   - Domain model separated from storage concerns
//   - Repository interface with a concrete in-memory implementation
//
package main

import (
	"context"
	"fmt"
	"sync"
)

type User struct {
	ID   string
	Name string
}

type UserRepository interface {
	Get(ctx context.Context, id string) (User, bool, error)
	Save(ctx context.Context, u User) error
}

type InMemoryUserRepo struct {
	mu sync.RWMutex
	m  map[string]User
}

func NewInMemoryUserRepo() *InMemoryUserRepo {
	return &InMemoryUserRepo{m: map[string]User{}}
}

func (r *InMemoryUserRepo) Get(ctx context.Context, id string) (User, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.m[id]
	return u, ok, nil
}

func (r *InMemoryUserRepo) Save(ctx context.Context, u User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.m[u.ID] = u
	return nil
}

func main() {
	ctx := context.Background()
	repo := NewInMemoryUserRepo()

	_ = repo.Save(ctx, User{ID: "u1", Name: "Ada"})
	u, ok, _ := repo.Get(ctx, "u1")

	fmt.Println("found?", ok, "user=", u)
}
