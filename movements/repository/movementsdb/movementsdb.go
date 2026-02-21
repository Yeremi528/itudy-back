package movementsdb

import (
	"context"

	"github.com/Yeremi528/itudy-back/movements"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionMovements = "movements"

type Repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Query(ctx context.Context, email string) ([]movements.Movement, error) {
	coll := r.db.Collection(collectionMovements)

	// Filtramos por email
	filter := bson.M{"email": email}

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var movementsdb []Movement
	if err := cursor.All(ctx, &movementsdb); err != nil {
		return nil, err
	}

	return toCore(movementsdb), nil
}

func (r *Repository) Insert(ctx context.Context, movement movements.Movement) error {
	coll := r.db.Collection(collectionMovements)

	doc := bson.M{
		"email":          movement.Email,
		"amount":         movement.Amount,
		"state":          movement.State,
		"testID":         movement.TestID,
		"transaction_id": movement.TransactionID,
		"created_at":     movement.CreatedAt,
		"updated_at":     movement.UpdatedAt,
	}

	_, err := coll.InsertOne(ctx, doc)
	return err
}

func (r *Repository) Update(ctx context.Context, movement movements.Movement) error {
	coll := r.db.Collection(collectionMovements)

	filter := bson.M{"transaction_id": movement.TransactionID}

	update := bson.M{
		"$set": bson.M{
			"state":      movement.State,
			"updated_at": movement.UpdatedAt,
		},
	}

	_, err := coll.UpdateOne(ctx, filter, update)
	return err
}
