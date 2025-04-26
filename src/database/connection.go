package database

import (
	"courses-service/src/config"

	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDBClient(config *config.Config) (*mongo.Client, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.DBURI))
	if err != nil {
		return nil, err
	}

	return client, nil
}
