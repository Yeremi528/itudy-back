package email

import (
	"context"
)

type Service interface {
	SendEmail(ctx context.Context, date, email, nameTest string) error
}
