package attribute

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository defines persistence operations for attributes.
type Repository interface {
	List(ctx context.Context, skip, limit int64) (ListResult, error)
	GetByID(ctx context.Context, attributeID int) (*Attribute, error)
}

type mongoRepository struct {
	collection *mongo.Collection
}

// NewRepository creates a MongoDB-backed attribute repository.
func NewRepository(collection *mongo.Collection) Repository {
	return &mongoRepository{collection: collection}
}

func (r *mongoRepository) List(ctx context.Context, skip, limit int64) (ListResult, error) {
	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return ListResult{}, err
	}

	opts := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "sortOrder", Value: 1}, {Key: "_id", Value: 1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return ListResult{}, err
	}
	defer cursor.Close(ctx)

	items := make([]Attribute, 0)
	if err := cursor.All(ctx, &items); err != nil {
		return ListResult{}, err
	}

	return ListResult{Items: items, Total: total}, nil
}

func (r *mongoRepository) GetByID(ctx context.Context, attributeID int) (*Attribute, error) {
	var attribute Attribute
	err := r.collection.FindOne(ctx, bson.M{"_id": attributeID}).Decode(&attribute)
	if err != nil {
		return nil, err
	}
	return &attribute, nil
}
