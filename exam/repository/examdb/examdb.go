package examdb

import (
	"context"

	"github.com/Yeremi528/itudy-back/exam"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	examCollection *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	examCollection := db.Collection("tests")
	return &Repository{
		examCollection: examCollection,
	}
}

// GetAllQuizzes recupera todos los exámenes de la colección
func (r *Repository) GetAllExam(ctx context.Context) ([]exam.Exam, error) {

	// 2. Slice donde guardaremos los resultados
	// Se inicializa vacío para devolver [] en lugar de nil si no hay datos
	quizzes := make([]exam.Exam, 0)

	// 3. Ejecutar el Find con un filtro vacío (bson.M{}) para traer todo
	cursor, err := r.examCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	// Aseguramos cerrar el cursor al terminar
	defer cursor.Close(ctx)

	// 4. Iterar y decodificar todos los documentos directamente en el slice
	if err := cursor.All(ctx, &quizzes); err != nil {
		return nil, err
	}

	return quizzes, nil
}
