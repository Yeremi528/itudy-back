package assignments

import (
	"context"
	"fmt"
)

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{repository: r}
}

func (s *service) AssignmentsAvailables(ctx context.Context, tech, level string) ([]AssignmentTest, error) {
	workers, err := s.repository.QueryWorkersByTechAndLevel(ctx, tech, level)
	if err != nil {
		return nil, fmt.Errorf("assignments.AssignmentsAvailables: %w", err)
	}

	for _, worker := range workers {
		assignments, err := s.repository.Assignments(ctx, tech, level, worker.DayAvailable)
		if err != nil {
			return nil, fmt.Errorf("assignments.AssignmentsAvailables: %w", err)
		}

		fmt.Printf("Worker: %+v\nAssignments: %+v\n", worker, assignments)
	}
	return nil, nil
}

func (s *service) AssignmentTestByUserID(ctx context.Context, userID string) ([]AssignmentTest, error) {
	return nil, nil
}
