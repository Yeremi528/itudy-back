package user

import "context"

// Service interface defines the set of functions that will containt the business logic, allowing
// CRUD operations and more.
type Service interface {
	CreateUser(ctx context.Context, user User) error
	GetUser(ctx context.Context, idOrEmail string) (User, error)
	UpdateUser(ctx context.Context, user User) error
}

type Repository interface {
	CreateUser(ctx context.Context, user User) error
	GetUser(ctx context.Context, idOrEmail string) (User, error)
	UpdateUser(ctx context.Context, user User) error
}
