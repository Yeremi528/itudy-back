package assigmentsdb

import (
	"context"
	"strings"
	"time"

	assignments "github.com/Yeremi528/itudy-back/assigments"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	assignmentsCollection *mongo.Collection
	workersCollection     *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	assignmentsCollection := db.Collection("assignments")
	workersCollection := db.Collection("workers")
	return &Repository{
		assignmentsCollection: assignmentsCollection,
		workersCollection:     workersCollection,
	}
}

func (r *Repository) QueryWorkersByTechAndLevel(ctx context.Context, tech, level string) ([]assignments.Worker, error) {
	var allowedLevels []string

	switch strings.ToLower(level) {
	case "junior":
		allowedLevels = []string{"Junior", "SemiSenior", "Senior"}
	case "semisenior", "semi-senior":
		allowedLevels = []string{"SemiSenior", "Senior"}
	case "senior":
		allowedLevels = []string{"Senior"}
	default:
		allowedLevels = []string{level}
	}

	filter := bson.M{
		"estado":              "activo",
		"perfil_tecnico.tech": tech,
		"perfil_tecnico.level": bson.M{
			"$in": allowedLevels,
		},
	}

	cursor, err := r.workersCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	results := make([]assignments.Worker, 0)
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *Repository) AssignmentsByWorker(ctx context.Context, workerID string) ([]assignments.AssignmentTest, error) {
	startDate := time.Now().Format("2006-01-02T15:04:05Z")
	endDate := time.Now().AddDate(0, 0, 8).Format("2006-01-02T15:04:05Z")

	filter := bson.M{
		"worker_id": workerID,
		"fecha_asignacion": bson.M{
			"$gte": startDate,
			"$lt":  endDate,
		},
	}

	cursor, err := r.assignmentsCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	results := make([]assignments.AssignmentTest, 0)
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *Repository) AssignmentTestByUserID(ctx context.Context, userID string) ([]assignments.AssignmentTest, error) {
	// 1. Obtener hora actual en UTC y formatear a STRING
	// Es CRÍTICO que este formato coincida con el que tienes en Mongo ("2026-02-07T12:00:00Z")
	nowStr := time.Now().UTC().Format("2006-01-02T15:04:05Z")

	filter := bson.M{
		"user_id": userID,
		"fecha_asignacion": bson.M{
			"$gte": nowStr, // "Dame fechas que sean MAYORES o IGUALES a ahora"
		},
	}

	// Opcional pero recomendado: Ordenar por fecha ascendente (la más próxima primero)
	opts := options.Find().SetSort(bson.D{{Key: "fecha_asignacion", Value: 1}})

	// Pasamos 'opts' al Find
	cursor, err := r.assignmentsCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	results := make([]assignments.AssignmentTest, 0)
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *Repository) CreateAssignment(ctx context.Context, assignment assignments.AssignmentTest) error {

	// 1. Convertimos la fecha time.Time a STRING ISO-8601
	// Esto asegura que se guarde igual que los otros registros: "2026-02-11T16:30:00Z"
	dateStr := assignment.FechaAsignacion.UTC().Format("2006-01-02T15:04:05Z")

	// 2. Creamos un mapa (bson.M) para tener control total de qué enviamos
	// En lugar de enviar el struct directo, enviamos este mapa.
	newDocument := bson.M{
		"worker_id":        assignment.WorkerID,
		"user_id":          assignment.UserID, // Agregado según tu última foto
		"test_id":          assignment.TestID,
		"fecha_asignacion": dateStr, // <--- AQUÍ ESTÁ EL CAMBIO (String)
		"estado":           assignment.Estado,
	}

	// 3. Manejo opcional del ID
	// Si tu struct ya trae un ID generado, lo usas. Si no, Mongo lo crea solo.
	if assignment.ID != "" {
		newDocument["_id"] = assignment.ID
	}

	// Insertamos el mapa, no el struct
	_, err := r.assignmentsCollection.InsertOne(ctx, newDocument)
	return err
}
