package courses

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Service interface {
	GetAllCourses(ctx context.Context, lang string) ([]TechAvailability, error)
	GetCourseByID(ctx context.Context, id string, lv string, lang string) (CourseByID, error)
	GetCourseContent(courseID string) (Content, error)
	GetLanguages(ctx context.Context) ([]Language, error)
}

type Repository interface {
	GetAllAvailableTechsByLang(ctx context.Context, langCode string) ([]TechAvailability, error)
	GetLanguages(ctx context.Context) ([]Language, error)
	GetCoursePath(ctx context.Context, tech string, level string, lang string) (CourseByID, error)
	GetCourseContent(ctx context.Context, courseID string, contentCollection *mongo.Collection) (Content, error)
}
