package email

import (
	"context"

	"github.com/Yeremi528/itudy-back/exam"
)

type Service interface {
	SendEmail(ctx context.Context, examInfo exam.Exam, date, email string) error
}
