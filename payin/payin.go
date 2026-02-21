package payin

import (
	"context"

	"github.com/Yeremi528/itudy-back/movements"
)

type Service interface {
	WebHook(ctx context.Context, ID, topic string) error
	RechargeLink(ctx context.Context, rut, IDAssignment, examName string, amount int) (string, error)
	Payin(ctx context.Context, movement movements.Movement) error
}

type Repository interface {
	Update(ctx context.Context, pay Pay) error
	QueryByEmail(ctx context.Context, email string) (Pay, error)
	Audit(ctx context.Context, mercadoPago MercadoPago) error
}

type repositoryAssignments interface {
	UpdateAssignment(ctx context.Context, ID string) error
}
