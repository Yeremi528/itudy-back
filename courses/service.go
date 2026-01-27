package courses

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

func (s *service) GetAllCourses(ctx context.Context, lang string) ([]TechAvailability, error) {

	courses, err := s.repository.GetAllAvailableTechsByLang(context.Background(), lang)
	if err != nil {
		return nil, fmt.Errorf("courses.GetAllCourses: %w", err)
	}
	return courses, nil

}

func (s *service) GetCourseByID(ctx context.Context, id string, lv string, lang string) (CourseByID, error) {
	courses, err := s.repository.GetCoursePath(ctx, id, lv, lang)
	if err != nil {
		return CourseByID{}, err
	}
	return courses, nil
}

func (s *service) GetCourseContent(courseID string) (Content, error) {
	// Implementa la lógica para obtener el contenido del curso
	coursesContent, err := s.repository.GetCourseContent(context.Background(), courseID, nil)
	if err != nil {
		return Content{}, err
	}
	return coursesContent, nil

}

func (s *service) GetLanguages(ctx context.Context) ([]Language, error) {
	languages, err := s.repository.GetLanguages(ctx)
	if err != nil {
		return nil, err
	}
	return languages, nil
}
