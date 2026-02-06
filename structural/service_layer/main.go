// service_layer.go
//
// This example demonstrates a Service Layer that coordinates use cases over repositories.
//
// Key ideas illustrated:
//
//   - Service contains business operations
//   - Repository interfaces are injected into service
//
package main

import (
	"context"
	"errors"
	"fmt"
)

type User struct {
	ID   string
	Name string
}

type UserRepository interface {
	Get(ctx context.Context, id string) (User, bool, error)
	Save(ctx context.Context, u User) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Rename(ctx context.Context, id, newName string) error {
	u, ok, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("user not found")
	}
	u.Name = newName
	return s.repo.Save(ctx, u)
}

// in-memory repo (for demo)
type memRepo struct{ store map[string]User }

func (m memRepo) Get(ctx context.Context, id string) (User, bool, error) {
	u, ok := m.store[id]
	return u, ok, nil
}
func (m memRepo) Save(ctx context.Context, u User) error {
	m.store[u.ID] = u
	return nil
}

func main() {
	ctx := context.Background()
	repo := memRepo{store: map[string]User{"u1": {ID: "u1", Name: "Grace"}}}

	svc := NewUserService(repo)
	_ = svc.Rename(ctx, "u1", "Grace Hopper")

	u, _, _ := repo.Get(ctx, "u1")
	fmt.Println(u)
}
