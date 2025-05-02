package testutil

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DBSetup holds the MongoDB client and database name for testing
type DBSetup struct {
	Client *mongo.Client
	DBName string
}

// SetupTestDB initializes a MongoDB client for testing
func SetupTestDB() *DBSetup {
	ctx := context.Background()
	uri := os.Getenv("DB_URI")
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	dbName := os.Getenv("DB_NAME")
	log.Printf("Initialized test database %s running on %s", dbName, uri)

	return &DBSetup{
		Client: client,
		DBName: dbName,
	}
}

// CleanupTestDB disconnects from the MongoDB client
func CleanupTestDB(client *mongo.Client) {
	if err := client.Disconnect(context.Background()); err != nil {
		log.Printf("Error disconnecting from database: %v", err)
	}
}

// CleanupCollection drops all documents from a collection
func (db *DBSetup) CleanupCollection(collection string) {
	coll := db.Client.Database(db.DBName).Collection(collection)
	_, err := coll.DeleteMany(context.Background(), bson.M{})
	if err != nil {
		log.Printf("Error cleaning up collection %s: %v", collection, err)
	}
}
