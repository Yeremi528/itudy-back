package oauth

import (
	"context"

	"github.com/Yeremi528/itudy-back/user"
)

type Service interface {
	GoogleLogin(ctx context.Context, token string) (user.User, error)
}
