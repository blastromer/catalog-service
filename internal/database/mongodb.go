package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Collection *mongo.Collection

// Init initializes the MongoDB client and connects to the database
func Init(dbName string) error {
	// Load the .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		return err
	}

	// Retrieve MongoDB connection details from environment variables
	mongoURI 		:= string(os.Getenv("MONGO_URI"))
	databaseName 	:= string(os.Getenv("MONGO_DATABASE"))
	collectionName 	:= dbName

	// Create a new MongoDB client
	client, err 	:= mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
		return err
	}

	// Connect to MongoDB
	ctx, cancel 	:= context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
		return err
	}

	// Assign the collection
	Collection = client.Database(databaseName).Collection(collectionName)
	return nil
}
