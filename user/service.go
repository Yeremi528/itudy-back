package user

import (
	"context"
	"fmt"
	"io"
	"time"
)

// type service is the implementation of Service interface containing all the business logic
// and dependencies required to complete the given tasks without exposing the implementation.
type service struct {
	repository Repository
	storage    Storage
}

func NewService(r Repository, s Storage) Service {
	return &service{repository: r, storage: s}
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

func (s *service) UploadProfileImage(ctx context.Context, userID, filename string, data io.Reader, contentType string) (string, error) {
	if s.storage == nil {
		return "", fmt.Errorf("user.UploadProfileImage: almacenamiento no configurado")
	}

	objectName := fmt.Sprintf("profiles/%s", userID)
	url, err := s.storage.UploadPublic(ctx, objectName, data, contentType)
	if err != nil {
		return "", fmt.Errorf("user.UploadProfileImage: %w", err)
	}

	if err := s.repository.UpdateImageURL(ctx, userID, url); err != nil {
		return "", fmt.Errorf("user.UploadProfileImage: %w", err)
	}

	return url, nil
}

func (s *service) UpdateStreak(ctx context.Context, userID string) error {
	u, err := s.repository.GetUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("user.UpdateStreak: %w", err)
	}

	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	lastLogin := u.Stats.LastLoginAt
	lastLoginDay := time.Date(lastLogin.Year(), lastLogin.Month(), lastLogin.Day(), 0, 0, 0, 0, time.UTC)

	var newStreak int
	switch {
	case lastLogin.IsZero() || lastLoginDay.Before(today.AddDate(0, 0, -1)):
		// Nunca inició sesión o faltó más de un día → reiniciar racha
		newStreak = 1
	case lastLoginDay.Equal(today):
		// Ya inició sesión hoy → sin cambios
		return nil
	default:
		// Inició sesión ayer → incrementar racha
		newStreak = u.Stats.StreakDays + 1
	}

	if err := s.repository.UpdateStreak(ctx, userID, newStreak, now); err != nil {
		return fmt.Errorf("user.UpdateStreak: %w", err)
	}

	return nil
}

func (s *service) AddAchievement(ctx context.Context, userID string, achievement Achievement) error {
	if err := s.repository.AddAchievement(ctx, userID, achievement); err != nil {
		return fmt.Errorf("user.AddAchievement: %w", err)
	}

	return nil
}
