package payin

import (
	"context"

	"github.com/Yeremi528/itudy-back/movements"
	"github.com/Yeremi528/itudy-back/user"
)

type Service interface {
	WebHook(ctx context.Context, ID, topic string) error
	RechargeLink(ctx context.Context, rut, IDAssignment, examName, fechaAsignacion string, amount int) (string, error)
	Payin(ctx context.Context, movement movements.Movement) error
}

type Repository interface {
	Update(ctx context.Context, pay Pay) error
	QueryByEmail(ctx context.Context, email string) (Pay, error)
	Audit(ctx context.Context, mercadoPago MercadoPago) error
}

type repositoryAssignments interface {
	UpdateAssignment(ctx context.Context, ID string) error
	QueryAssignmentTestByID(ctx context.Context, ID string) (string, error)
}

type userService interface {
	GetUser(ctx context.Context, idOrEmail string) (user.User, error)
	AddAchievement(ctx context.Context, userID string, achievement user.Achievement) error
}
