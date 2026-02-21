package movements

import (
	"context"
	"fmt"
)

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{repository: r}
}

func (s *service) Query(ctx context.Context, rut string) ([]Movement, error) {
	movements, err := s.repository.Query(ctx, rut)
	if err != nil {
		return nil, fmt.Errorf("movements.Query: repository.query: %w", err)
	}

	return movements, nil
}

func (s *service) Insert(ctx context.Context, movement Movement) error {

	if err := s.repository.Insert(ctx, movement); err != nil {
		return fmt.Errorf("movements.insert: repository.insert: %w", err)
	}

	return nil
}
