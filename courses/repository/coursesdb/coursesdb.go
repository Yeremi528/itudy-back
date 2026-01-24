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

// GetAllAvailableTechs utiliza 'Distinct' para encontrar todos los lenguajes/tecnologías únicas
// disponibles en la colección 'tracks'.
func (r *Repository) GetAllAvailableTechs(ctx context.Context) ([]string, error) {
	// El segundo parámetro de Distinct es el nombre del campo que queremos obtener
	collection := r.db.Collection("tracks")

	values, err := collection.Distinct(ctx, "tech", bson.D{})
	if err != nil {
		return nil, err
	}

	// Convertir el resultado de []interface{} a []string
	techs := make([]string, len(values))
	for i, v := range values {
		// Asumimos que los valores del campo "tech" son strings
		if s, ok := v.(string); ok {
			techs[i] = s
		}
	}

	return techs, nil
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

	// Loguear el valor exacto que se está buscando para debugging
	fmt.Println("🔎 Buscando contenido para course_ref_id:", courseID)

	var content Content
	// fmt.Println(collection, "que tiene la coleccion") // Este log no es necesario

	// Usar FindOne ya que solo se espera UN documento
	err := collection.FindOne(ctx, filter).Decode(&content)

	// Loguear el resultado o el error
	fmt.Println("Contenido obtenido:", content)

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
