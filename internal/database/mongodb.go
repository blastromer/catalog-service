package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Store wraps a MongoDB client and the active collection for a service.
type Store struct {
	Client     *mongo.Client
	Database   *mongo.Database
	Collection *mongo.Collection
}

// Connect establishes a MongoDB connection, pings the server, and returns a Store.
func Connect(ctx context.Context, uri, databaseName, collectionName string) (*Store, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("connect mongo: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, fmt.Errorf("ping mongo: %w", err)
	}

	db := client.Database(databaseName)
	return &Store{
		Client:     client,
		Database:   db,
		Collection: db.Collection(collectionName),
	}, nil
}

// Close disconnects the MongoDB client.
func (s *Store) Close(ctx context.Context) error {
	if s == nil || s.Client == nil {
		return nil
	}
	return s.Client.Disconnect(ctx)
}

// Ping verifies the database connection is still healthy.
func (s *Store) Ping(ctx context.Context) error {
	if s == nil || s.Client == nil {
		return fmt.Errorf("mongo client is not initialized")
	}
	return s.Client.Ping(ctx, readpref.Primary())
}
