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
	// Pasamos "" como audience para que la librería valide firma, expiración e issuer,
	// y nosotros validamos el Audience manualmente contra nuestros 3 Client IDs.
	payload, err := idtoken.Validate(ctx, token, "")
	if err != nil {
		return user.User{}, fmt.Errorf("oauth.GoogleLogin: token inválido o expirado: %w", err)
	}

	// Validamos que el Audience del token corresponda a Web, Android o iOS
	validAudience := false
	for _, clientID := range s.Cfg.GoogleClientIDs {
		// Verificamos si la Audiencia coincide. 
		// En algunos casos de iOS/Android, el Audience es el WebClientID, pero el 'azp' (Authorized Party)
		// es el ClientID específico de iOS o Android. Validamos ambos.
		if payload.Audience == clientID {
			validAudience = true
			break
		}
		
		if azp, ok := payload.Claims["azp"].(string); ok && azp == clientID {
			validAudience = true
			break
		}
	}

	if !validAudience {
		// Log para debuguear fácilmente qué Audience llegó
		fmt.Printf("Token Audience rechazado: %s. azp: %v\n", payload.Audience, payload.Claims["azp"])
		return user.User{}, fmt.Errorf("oauth.GoogleLogin: Client ID no autorizado (Audience: %s)", payload.Audience)
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

	// Actualizar racha de días consecutivos
	if err := s.userService.UpdateStreak(ctx, u.ID); err != nil {
		// No bloqueamos el login por un error en la racha
		fmt.Printf("oauth.GoogleLogin: UpdateStreak: %v\n", err)
	}

	return u, nil

}
