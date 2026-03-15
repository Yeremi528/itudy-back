package learningdb

import (
	"context"
	"fmt"
	"time"

	"github.com/Yeremi528/itudy-back/learning" // Ajusta este import a tu ruta real
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		db: db,
	}
}

// CreateLesson guarda una nueva lección (contenido) en la base de datos
func (r *Repository) CreateLesson(ctx context.Context, lesson learning.Lesson) error {
	collection := r.db.Collection("lessons")

	// Opcional: Si quieres manejar timeouts específicos para la DB
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, lesson)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) CompleteLesson(ctx context.Context, userID, courseID, completedLessonID, nextLessonID string, totalLessons int) error {
	collection := r.db.Collection("lessons")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID, "course_id": courseID}

	fieldsToUpdate := bson.M{
		"last_lesson_id": nextLessonID,
		"last_practiced": time.Now(),
		fmt.Sprintf("lessons_progress.%s.status", completedLessonID): "completed",
	}

	if completedLessonID != nextLessonID {
		fieldsToUpdate[fmt.Sprintf("lessons_progress.%s.status", nextLessonID)] = "active"
	}

	// Guardar total de lecciones si se conoce y aún no está guardado
	if totalLessons > 0 {
		fieldsToUpdate["total_lessons"] = totalLessons
	}

	result, err := collection.UpdateOne(ctx, filter, bson.M{"$set": fieldsToUpdate}, options.Update().SetUpsert(false))
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("no progress found for user %s in course %s", userID, courseID)
	}

	return nil
}

func (r *Repository) IncrementLessonXP(ctx context.Context, userID, courseID string, xp int) error {
	collection := r.db.Collection("lessons")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID, "course_id": courseID}
	update := bson.M{"$inc": bson.M{"current_xp": xp}}

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

// GetLessonByID busca una lección por su ID (ej: "go_001")
func (r *Repository) GetLessonByID(ctx context.Context, userID, courseID string) (learning.Lesson, error) {
	collection := r.db.Collection("lessons")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var lesson learning.Lesson
	filter := bson.M{"user_id": userID, "course_id": courseID}

	err := collection.FindOne(ctx, filter).Decode(&lesson)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Puedes retornar un error personalizado de dominio aquí si prefieres
			// ej: return learning.Lesson{}, learning.ErrLessonNotFound
			return learning.Lesson{}, nil
		}
		return learning.Lesson{}, err
	}

	return lesson, nil
}
