package exam

import "context"

type Service interface {
	ExamByUserID(ctx context.Context, userID string) ([]Exam, error)
	ExambyID(ctx context.Context, ID string) Exam
}

type Repository interface {
	GetAllExam(ctx context.Context) ([]Exam, error)
}
