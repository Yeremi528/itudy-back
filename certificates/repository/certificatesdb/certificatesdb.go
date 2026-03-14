package certificatesdb

import (
	"context"
	"time"

	"github.com/Yeremi528/itudy-back/certificates"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(ctx context.Context, cert certificates.Certificate) error {
	if cert.CreatedAt.IsZero() {
		cert.CreatedAt = time.Now().UTC()
	}
	_, err := r.db.Collection("certificates").InsertOne(ctx, cert)
	return err
}

func (r *Repository) GetByUserID(ctx context.Context, userID string) ([]certificates.Certificate, error) {
	filter := bson.M{"user_id": userID}
	opts := options.Find().SetSort(bson.M{"approved_at": -1}) // más reciente primero

	cursor, err := r.db.Collection("certificates").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var certs []certificates.Certificate
	if err := cursor.All(ctx, &certs); err != nil {
		return nil, err
	}

	if certs == nil {
		certs = []certificates.Certificate{}
	}

	return certs, nil
}
