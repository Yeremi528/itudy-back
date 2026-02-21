package payindb

import (
	"context"

	"github.com/Yeremi528/itudy-back/payin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	collectionPayin = "payins"
	collectionAudit = "audits"
)

type Storer struct {
	db *mongo.Database
}

func New(db *mongo.Database) *Storer {
	return &Storer{
		db: db,
	}
}

func (s *Storer) QueryByEmail(ctx context.Context, email string) (payin.Pay, error) {
	coll := s.db.Collection(collectionPayin)

	var dbbalance dbBalance
	filter := bson.M{"email": email}

	err := coll.FindOne(ctx, filter).Decode(&dbbalance)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return payin.Pay{}, nil
		}
		return payin.Pay{}, err
	}

	return toCoreBalance(dbbalance), nil
}

func (s *Storer) Update(ctx context.Context, pay payin.Pay) error {
	coll := s.db.Collection(collectionPayin)

	filter := bson.M{"email": pay.Email}
	update := bson.M{
		"$set": bson.M{
			"balance": pay.Amount,
		},
	}

	_, err := coll.UpdateOne(ctx, filter, update)
	return err
}

func (s *Storer) Audit(ctx context.Context, mercadoPago payin.MercadoPago) error {
	coll := s.db.Collection(collectionAudit)

	doc := bson.M{
		"id_mercado_pago": mercadoPago.ID,
		"topic":           mercadoPago.Topic,
		"state":           mercadoPago.State,
		"created_at":      bson.M{"$currentDate": true}, // Recomendado para auditoría
	}

	_, err := coll.InsertOne(ctx, doc)
	return err
}
