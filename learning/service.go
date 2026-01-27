package learning

import (
	"context"
	"fmt"
)

type service struct {
	reportory Repository
}

func NewService(r Repository) Service {
	return &service{reportory: r}
}

func (s *service) CreateLesson(ctx context.Context, lesson Lesson) error {
	if err := s.reportory.CreateLesson(ctx, lesson); err != nil {
		return fmt.Errorf("learning.CreateLesson: %w", err)
	}

	return nil
}

func (s *service) UpdateLesson(ctx context.Context, userID, courseID, lessonID, nextID string) error {
	lesson, err := s.reportory.GetLessonByID(ctx, userID, courseID)
	if err != nil {
		return fmt.Errorf("learning.UpdateLesson: %w", err)
	}

	if validateLastLesson(lesson, nextID) {
		return nil
	}

	if err := s.reportory.CompleteLesson(ctx, userID, courseID, lessonID, nextID); err != nil {
		return fmt.Errorf("learning.UpdateLesson: %w", err)
	}

	return nil
}

func (s *service) GetLessonByID(ctx context.Context, userID, courseID string) (Lesson, error) {
	lesson, err := s.reportory.GetLessonByID(ctx, userID, courseID)
	if err != nil {
		return Lesson{}, fmt.Errorf("learning.GetLessonByID: %w", err)
	}

	if lesson.ID == "" {
		s.CreateLesson(ctx, Lesson{
			UserID:          userID,
			CourseID:        courseID,
			CurrentXP:       0,
			LastLessonID:    "",
			LessonsProgress: make(map[string]LessonProgress),
		})
		lesson, err = s.reportory.GetLessonByID(ctx, userID, courseID)
		if err != nil {
			return Lesson{}, fmt.Errorf("learning.GetLessonByID: %w", err)
		}
	}

	return lesson, nil
}

func validateLastLesson(lesson Lesson, lessonID string) bool {
	for k, v := range lesson.LessonsProgress {
		if k == lessonID && v.Status == "completed" {
			return true
		}
	}
	return false
}
