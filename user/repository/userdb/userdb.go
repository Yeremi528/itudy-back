package userdb

import (
	"context"
	"errors"
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
