package exam

import (
	"context"
	"fmt"

	"github.com/Yeremi528/itudy-back/user"
)

type service struct {
	repository     Repository
	userRepository user.Repository
	exam           map[string]Exam
	examByID       map[string]Exam
}

func NewService(r Repository, userRepository user.Repository) Service {

	exam, err := r.GetAllExam(context.Background())
	if err != nil {
		panic(fmt.Sprintf("exam.NewService: %v", err))
	}

	exams := make(map[string]Exam)
	for _, e := range exam {
		exams[fmt.Sprintf("%s_%s", e.Language, e.DifficultyLevel)] = e
	}

	examsByID := make(map[string]Exam)
	for _, e := range exam {
		examsByID[e.ID] = e
	}

	return &service{userRepository: userRepository, exam: exams, examByID: examsByID}
}

func (s *service) ExambyID(ctx context.Context, ID string) Exam {
	return s.examByID[ID]
}

// Tengo que probar
func (s *service) ExamByUserID(ctx context.Context, email string) ([]Exam, error) {
	userInfo, err := s.userRepository.GetUser(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("exam.ExamByUserID: %w", err)
	}

	var exams []Exam
	for _, list := range userInfo.CoursesInfo.List {
		if list.Progress >= 60 {
			exams = append(exams, s.exam[list.ID])
		}
	}

	return exams, nil
}
