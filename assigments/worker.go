package assignments

import "context"

type Service interface {
	AssignmentsAvailables(ctx context.Context, xtech, level string) ([]Assignment, error)
	TaskByUserID(ctx context.Context, userID string) (Task, error)
}

type Repository interface {
	QueryWorkersByTechAndLevel(ctx context.Context, tech, level string) ([]Worker, error)
	Assignments(ctx context.Context, tech, level, day string) ([]Assignment, error)
	TaskByUserID(ctx context.Context, userID string) (Task, error)
}
