package learning

import (
	"context"
	"fmt"
)

type service struct {
	reportory   Repository
	userUpdater userProgressUpdater
}

func NewService(r Repository, u userProgressUpdater) Service {
	return &service{reportory: r, userUpdater: u}
}

func (s *service) CreateLesson(ctx context.Context, lesson Lesson) error {
	if err := s.reportory.CreateLesson(ctx, lesson); err != nil {
		return fmt.Errorf("learning.CreateLesson: %w", err)
	}
	return nil
}

func (s *service) UpdateLesson(ctx context.Context, userID, courseID, lessonID, nextID string, totalLessons int) (LessonUpdateResult, error) {
	lesson, err := s.reportory.GetLessonByID(ctx, userID, courseID)
	if err != nil {
		return LessonUpdateResult{}, fmt.Errorf("learning.UpdateLesson: %w", err)
	}

	// Si la lección ya estaba completada, devolvemos el estado actual sin sumar XP
	if validateLastLesson(lesson, lessonID) {
		return buildResult(lesson, 0), nil
	}

	// Obtener racha del usuario para el multiplicador de XP
	u, err := s.userUpdater.GetUser(ctx, userID)
	if err != nil {
		return LessonUpdateResult{}, fmt.Errorf("learning.UpdateLesson: getUser: %w", err)
	}
	xpEarned := calculateXP(u.Stats.StreakDays)

	// Completar lección en DB
	if err := s.reportory.CompleteLesson(ctx, userID, courseID, lessonID, nextID, totalLessons); err != nil {
		return LessonUpdateResult{}, fmt.Errorf("learning.UpdateLesson: completeLesson: %w", err)
	}

	// Acumular XP en la lección
	if err := s.reportory.IncrementLessonXP(ctx, userID, courseID, xpEarned); err != nil {
		return LessonUpdateResult{}, fmt.Errorf("learning.UpdateLesson: incrementLessonXP: %w", err)
	}

	// Actualizar XP total del usuario
	if err := s.userUpdater.IncrementXP(ctx, userID, xpEarned); err != nil {
		return LessonUpdateResult{}, fmt.Errorf("learning.UpdateLesson: incrementXP: %w", err)
	}

	// Calcular progreso del curso
	updatedLesson, err := s.reportory.GetLessonByID(ctx, userID, courseID)
	if err != nil {
		return LessonUpdateResult{}, fmt.Errorf("learning.UpdateLesson: reloadLesson: %w", err)
	}

	completed := countCompleted(updatedLesson.LessonsProgress)
	total := updatedLesson.TotalLessons
	var progress float64
	if total > 0 {
		progress = float64(completed) / float64(total) * 100
	}
	isCompleted := total > 0 && completed >= total

	if err := s.userUpdater.UpdateCourseProgress(ctx, userID, courseID, progress, isCompleted); err != nil {
		// No bloqueamos el flujo por un error de progreso
		fmt.Printf("learning.UpdateLesson: UpdateCourseProgress: %v\n", err)
	}

	return LessonUpdateResult{
		XPEarned:    xpEarned,
		TotalXP:     u.Stats.TotalXP + xpEarned,
		Progress:    progress,
		IsCompleted: isCompleted,
	}, nil
}

// calculateXP devuelve los XP a otorgar según los días de racha del usuario.
//
//	 Racha    Multiplicador   XP
//	 1–2 días     ×1.00       20
//	 3–6 días     ×1.25       25
//	 7–13 días    ×1.50       30
//	 14+ días     ×2.00       40
func calculateXP(streakDays int) int {
	const base = 20
	switch {
	case streakDays >= 14:
		return base * 2
	case streakDays >= 7:
		return int(float64(base) * 1.5)
	case streakDays >= 3:
		return int(float64(base) * 1.25)
	default:
		return base
	}
}

func countCompleted(progress map[string]LessonProgress) int {
	n := 0
	for _, p := range progress {
		if p.Status == StatusCompleted {
			n++
		}
	}
	return n
}

func buildResult(lesson Lesson, xpEarned int) LessonUpdateResult {
	completed := countCompleted(lesson.LessonsProgress)
	var progress float64
	if lesson.TotalLessons > 0 {
		progress = float64(completed) / float64(lesson.TotalLessons) * 100
	}
	return LessonUpdateResult{
		XPEarned:    xpEarned,
		TotalXP:     lesson.CurrentXP,
		Progress:    progress,
		IsCompleted: lesson.TotalLessons > 0 && completed >= lesson.TotalLessons,
	}
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
