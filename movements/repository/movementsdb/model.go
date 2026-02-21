package movementsdb

import (
	"time"

	"github.com/Yeremi528/itudy-back/movements"
)

type Movement struct {
	ID            string    `bson:"_id,omitempty"` // Mongo usa _id por defecto
	Email         string    `bson:"email"`
	Amount        float64   `bson:"amount"`
	TestID        string    `bson:"test_id"`
	State         int       `bson:"state"`
	TransactionID string    `bson:"transaction_id"`
	CreatedAt     time.Time `bson:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at"`
}

func toCore(movementsdb []Movement) []movements.Movement {
	var movementsCore []movements.Movement
	for _, m := range movementsdb {
		movementsCore = append(movementsCore, movements.Movement{
			ID:            m.ID,
			Email:         m.Email,
			Amount:        m.Amount,
			TestID:        m.TestID,
			State:         m.State,
			TransactionID: m.TransactionID,
			CreatedAt:     m.CreatedAt,
			UpdatedAt:     m.UpdatedAt,
		})
	}

	return movementsCore
}
