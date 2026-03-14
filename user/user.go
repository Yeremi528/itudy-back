package user

import (
	"context"
	"io"
	"time"
)

// Storage define el contrato para subir archivos al almacenamiento en la nube.
type Storage interface {
	UploadPublic(ctx context.Context, objectName string, data io.Reader, contentType string) (string, error)
}

// Service interface defines the set of functions that will containt the business logic, allowing
// CRUD operations and more.
type Service interface {
	CreateUser(ctx context.Context, user User) error
	GetUser(ctx context.Context, idOrEmail string) (User, error)
	UpdateUser(ctx context.Context, user User) error
	UploadProfileImage(ctx context.Context, userID string, filename string, data io.Reader, contentType string) (string, error)
	UpdateStreak(ctx context.Context, userID string) error
	AddAchievement(ctx context.Context, userID string, achievement Achievement) error
}

type Repository interface {
	CreateUser(ctx context.Context, user User) error
	GetUser(ctx context.Context, idOrEmail string) (User, error)
	UpdateUser(ctx context.Context, user User) error
	UpdateImageURL(ctx context.Context, userID, url string) error
	UpdateStreak(ctx context.Context, userID string, streakDays int, lastLoginAt time.Time) error
	AddAchievement(ctx context.Context, userID string, achievement Achievement) error
}
