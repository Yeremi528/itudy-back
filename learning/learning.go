package learning

import "context"

type Service interface {
	CreateLesson(ctx context.Context, enrollment Lesson) error
	UpdateLesson(ctx context.Context, userID, courseID, lessonID, nextID string) error
	GetLessonByID(ctx context.Context, userID, courseID string) (Lesson, error)
}

type Repository interface {
	CreateLesson(ctx context.Context, lesson Lesson) error
	CompleteLesson(ctx context.Context, userID, courseID, lessonID, nextID string) error
	GetLessonByID(ctx context.Context, userID, courseID string) (Lesson, error)
}
