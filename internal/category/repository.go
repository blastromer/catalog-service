package category

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository defines persistence operations for categories.
type Repository interface {
	List(ctx context.Context, skip, limit int64) (ListResult, error)
	GetByCategoryID(ctx context.Context, categoryID int) (*Category, error)
	FindByFilter(ctx context.Context, filter bson.M) ([]Category, error)
}

type mongoRepository struct {
	collection *mongo.Collection
}

// NewRepository creates a MongoDB-backed category repository.
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
		SetSort(bson.D{{Key: "sort", Value: 1}, {Key: "categoryId", Value: 1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return ListResult{}, err
	}
	defer cursor.Close(ctx)

	items := make([]Category, 0)
	if err := cursor.All(ctx, &items); err != nil {
		return ListResult{}, err
	}

	return ListResult{Items: items, Total: total}, nil
}

func (r *mongoRepository) GetByCategoryID(ctx context.Context, categoryID int) (*Category, error) {
	var category Category
	err := r.collection.FindOne(ctx, bson.M{"categoryId": categoryID}).Decode(&category)
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *mongoRepository) FindByFilter(ctx context.Context, filter bson.M) ([]Category, error) {
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetLimit(100))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	items := make([]Category, 0)
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}
