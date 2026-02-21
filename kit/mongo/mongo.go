package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client *mongo.Client
)

func ConnectMongo(conexion string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(conexion)

	var err error
	Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error conectando a MongoDB: %w", err)
	}

	if err := Client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("error al hacer ping a MongoDB: %w", err)
	}

	db := Client.Database("itudy")

	return db, nil
}

func CloseMongo() {
	if Client != nil {
		if err := Client.Disconnect(context.Background()); err != nil {
			fmt.Printf("❌ Error cerrando MongoDB: %v\n", err)
		} else {
			fmt.Println("✅ Conexión a MongoDB cerrada")
		}
	}
}
