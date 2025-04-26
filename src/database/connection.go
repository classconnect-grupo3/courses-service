package database

import (
	"courses-service/src/config"
	"log"

	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDBClient(config *config.Config) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(config.DBURI).SetAuth(options.Credential{
		Username: config.DBUsername,
		Password: config.DBPassword,
	})

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return nil, err
	}

	log.Println("Connected to database")

	return client, nil
}
