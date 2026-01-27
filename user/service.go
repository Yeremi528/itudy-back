package user

import (
	"context"
	"fmt"
)

// type service is the implementation of Service interface containing all the business logic
// and dependencies required to complete the given tasks without exposing the implementation.
type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{repository: r}
}

func (s *service) CreateUser(ctx context.Context, user User) error {
	if err := s.repository.CreateUser(ctx, user); err != nil {
		return fmt.Errorf("user.CreateUser: %w", err)
	}

	return nil
}

func (s *service) GetUser(ctx context.Context, idOrEmail string) (User, error) {
	u, err := s.repository.GetUser(ctx, idOrEmail)
	if err != nil {
		return User{}, fmt.Errorf("user.GetUser: %w", err)
	}

	return u, nil
}

func (s *service) UpdateUser(ctx context.Context, user User) error {
	fmt.Printf("%+v", user)
	if err := s.repository.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("user.UpdateUser: %w", err)
	}

	return nil
}
