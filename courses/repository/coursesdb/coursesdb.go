package coursesdb

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Yeremi528/itudy-back/courses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		db: db,
	}
}

// GetAllAvailableTechsByLang filtra por idioma (ej: "es") y devuelve
// las tecnologías con sus niveles disponibles para ese idioma.
func (r *Repository) GetAllAvailableTechsByLang(ctx context.Context, langCode string) ([]courses.TechAvailability, error) {
	collection := r.db.Collection("courses")

	pipeline := mongo.Pipeline{
		// 1. FILTRAR ($match): Solo documentos que coincidan con el idioma solicitado
		{{Key: "$match", Value: bson.D{
			{Key: "lang", Value: langCode},
		}}},

		// 2. AGRUPAR ($group): Por tecnología, acumulando los niveles únicos
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$tech"}, // Agrupar por 'tech' (ej: Golang)
			{Key: "levels", Value: bson.D{
				{Key: "$addToSet", Value: "$level"}, // Crear array de niveles sin duplicados
			}},
		}}},

		// 3. PROYECTAR ($project): Dar formato final al JSON
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},       // Quitar _id interno
			{Key: "tech", Value: "$_id"}, // Poner el nombre de la tech
			{Key: "levels", Value: 1},    // Mantener el array de niveles
		}}},

		// 4. ORDENAR ($sort): Opcional, para que salga ordenado alfabéticamente
		{{Key: "$sort", Value: bson.D{{Key: "tech", Value: 1}}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []courses.TechAvailability
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// GetCoursePath obtiene el documento Track completo
func (r *Repository) GetCoursePath(ctx context.Context, tech string, level string, lang string) (courses.CourseByID, error) {
	// 1. SANEAR las variables de entrada
	cleanTech := strings.TrimSpace(tech)
	cleanLevel := strings.TrimSpace(level)
	cleanLang := strings.TrimSpace(lang)

	// 2. Definir el filtro con los valores limpios
	filter := bson.D{
		{"tech", cleanTech},
		{"level", cleanLevel},
		{"lang", cleanLang},
	}

	var track Track

	collection := r.db.Collection("courses")
	// Aquí se utiliza el contexto y el filtro
	err := collection.FindOne(ctx, filter).Decode(&track)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return courses.CourseByID{}, errors.New("no se encontró el path para los filtros especificados")
		}
		return courses.CourseByID{}, err
	}

	return trackToCourse(track), nil
}

// GetCourseContent obtiene la teoría y ejercicios completos de un curso,
// buscando el ÚNICO elemento que coincide con courseID en el campo 'course_ref_id'.
// NOTA: Se ignora 'contentCollection' ya que la colección se obtiene de r.db.
func (r *Repository) GetCourseContent(ctx context.Context, courseID string, contentCollection *mongo.Collection) (courses.Content, error) {

	// Asume que la colección es "courses_content"
	collection := r.db.Collection("courses_content")

	// Definir el filtro
	filter := bson.D{{"course_ref_id", courseID}}

	var content Content

	// Usar FindOne ya que solo se espera UN documento
	err := collection.FindOne(ctx, filter).Decode(&content)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// Este es el error que recibes si no hay coincidencia
			return courses.Content{}, errors.New("contenido del curso no encontrado")
		}
		// Manejar otros errores de la BD
		return courses.Content{}, fmt.Errorf("error al buscar en MongoDB: %w", err)
	}

	// Mapear el contenido de la DB al struct del dominio y retornar
	return contentToCourse(content), nil
}

func (r *Repository) GetLanguages(ctx context.Context) ([]courses.Language, error) {
	// CAMBIO: El nombre exacto en tu captura es "languages"
	collection := r.db.Collection("languages")

	cursor, err := collection.Find(ctx, bson.M{"is_active": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var languages []Language
	if err = cursor.All(ctx, &languages); err != nil {
		return nil, err
	}

	return languageToLanguage(languages), nil
}
