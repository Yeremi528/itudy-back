package oauth

import (
	"context"
	"fmt"

	"github.com/Yeremi528/itudy-back/user"
	"github.com/google/uuid"

	"google.golang.org/api/idtoken"
)

type service struct {
	Cfg         Config
	userService user.Service
}

func NewService(cfg Config, userService user.Service) Service {
	return &service{Cfg: cfg, userService: userService}
}

func (s *service) GoogleLogin(ctx context.Context, token string) (user.User, error) {
	payload, err := idtoken.Validate(ctx, token, s.Cfg.GoogleClientID)
	if err != nil {
		return user.User{}, fmt.Errorf("oauth.GoogleLogin: %w", err)
	}

	email := payload.Claims["email"].(string)
	u, err := s.userService.GetUser(ctx, email)
	if err != nil {
		return user.User{}, fmt.Errorf("oauth.GoogleLogin: %w", err)
	}

	if u.Email == "" {
		u = user.User{
			ID:    uuid.New().String(),
			Email: email,
			Name:  payload.Claims["name"].(string),
		}
		if err := s.userService.CreateUser(ctx, u); err != nil {
			return user.User{}, fmt.Errorf("oauth.GoogleLogin: %w", err)
		}

		return u, nil
	}

	return u, nil

}
