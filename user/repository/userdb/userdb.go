package userdb

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Yeremi528/itudy-back/user"
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

func (r *Repository) CreateUser(ctx context.Context, u user.User) error {
	collection := r.db.Collection("users")

	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now().UTC()
	}

	_, err := collection.InsertOne(ctx, u)
	if err != nil {
		return err
	}

	return nil
}
func (r *Repository) UpdateUser(ctx context.Context, u user.User) error {
	collection := r.db.Collection("users")
	// 1. CORRECCIÓN DEL ID: Convertir string a ObjectID

	filter := bson.M{"_id": u.ID}

	// 2. CORRECCIÓN DE DATOS VACÍOS: Solo agregamos lo que tenga valor
	updateFields := bson.M{}

	// Actualizamos CoursesInfo si tiene contenido
	if u.CoursesInfo.Active != "" || len(u.CoursesInfo.List) > 0 {
		updateFields["courses_info"] = u.CoursesInfo
	}

	// Agrega aquí el resto de validaciones si quieres actualizar otros campos a la vez
	// (Ejemplo: si quisieras actualizar Stats también)
	if u.Stats != (user.Stats{}) {
		updateFields["stats"] = u.Stats
	}
	// ... repetir para otros campos (Name, Email) si fuera necesario

	// Si no hay nada que actualizar, salimos
	if len(updateFields) == 0 {
		return nil
	}

	update := bson.M{"$set": updateFields}

	// 3. Ejecutar Update
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	// Opcional: Verificar si realmente encontró el documento
	if result.MatchedCount == 0 {
		return fmt.Errorf("no se encontró el usuario con id %s", u.ID)
	}

	return nil
}
func (r *Repository) GetUser(ctx context.Context, idOrEmail string) (user.User, error) {
	collection := r.db.Collection("users")

	var filter bson.M

	if strings.Contains(idOrEmail, "@") {
		filter = bson.M{"email": idOrEmail}
	} else {
		filter = bson.M{"id": idOrEmail}
	}
	var u user.User
	err := collection.FindOne(ctx, filter).Decode(&u)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user.User{}, nil
		}
		return user.User{}, err
	}

	return u, nil
}
