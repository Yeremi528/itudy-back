package movements

import "context"

type Service interface {
	Query(ctx context.Context, rut string) ([]Movement, error)
}

type Repository interface {
	Query(ctx context.Context, rut string) ([]Movement, error)
	Insert(ctx context.Context, movement Movement) error
	Update(ctx context.Context, movement Movement) error
}
