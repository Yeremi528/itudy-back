package learning

import (
	"context"

	"github.com/Yeremi528/itudy-back/user"
)

type Service interface {
	CreateLesson(ctx context.Context, enrollment Lesson) error
	// totalLessons: total de lecciones del curso (enviado por la app).
	UpdateLesson(ctx context.Context, userID, courseID, lessonID, nextID string, totalLessons int) (LessonUpdateResult, error)
	GetLessonByID(ctx context.Context, userID, courseID string) (Lesson, error)
}

type Repository interface {
	CreateLesson(ctx context.Context, lesson Lesson) error
	CompleteLesson(ctx context.Context, userID, courseID, lessonID, nextID string, totalLessons int) error
	IncrementLessonXP(ctx context.Context, userID, courseID string, xp int) error
	GetLessonByID(ctx context.Context, userID, courseID string) (Lesson, error)
}

// userProgressUpdater es el subconjunto de user.Repository que necesita el módulo learning.
type userProgressUpdater interface {
	GetUser(ctx context.Context, idOrEmail string) (user.User, error)
	IncrementXP(ctx context.Context, userID string, xp int) error
	UpdateCourseProgress(ctx context.Context, userID, courseID string, progress float64, isCompleted bool) error
}
