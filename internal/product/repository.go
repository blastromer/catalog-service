package product

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository defines persistence operations for products.
type Repository interface {
	List(ctx context.Context, skip, limit int64) (ListResult, error)
	GetByProductID(ctx context.Context, productID int) (*Product, error)
}

type mongoRepository struct {
	collection *mongo.Collection
}

// NewRepository creates a MongoDB-backed product repository.
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
		SetSort(bson.D{{Key: "productId", Value: 1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return ListResult{}, err
	}
	defer cursor.Close(ctx)

	items := make([]Product, 0)
	if err := cursor.All(ctx, &items); err != nil {
		return ListResult{}, err
	}

	return ListResult{Items: items, Total: total}, nil
}

func (r *mongoRepository) GetByProductID(ctx context.Context, productID int) (*Product, error) {
	var product Product
	err := r.collection.FindOne(ctx, bson.M{"productId": productID}).Decode(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}
