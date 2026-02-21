package assignments

import "context"

type Service interface {
	AssignmentsAvailables(ctx context.Context, tech, level string) ([]AssignmentTest, error)
	AssignmentTestByUserID(ctx context.Context, userID string) ([]AssignmentTest, error)
	CreateAssignmentsByUserID(ctx context.Context, assignment AssignmentTest) (string, error)
}

type Repository interface {
	QueryWorkersByTechAndLevel(ctx context.Context, tech, level string) ([]Worker, error)
	AssignmentsByWorker(ctx context.Context, workerID string) ([]AssignmentTest, error)
	AssignmentTestByUserID(ctx context.Context, userID string) ([]AssignmentTest, error)
	CreateAssignment(ctx context.Context, assignment AssignmentTest) (string, error)
}
