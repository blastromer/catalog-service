package product

import (
	"context"
	"micro-services/catalog/internal/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FetchAllProducts retrieves products with pagination
func FetchAllProducts(page, limit int) ([]Product, error) {
	skip := (page - 1) * limit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := database.Collection.Find(ctx, bson.M{}, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit)))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

// FetchProductByID retrieves a product by its ID
func FetchProductByID(productID int) (*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var product Product
	err := database.Collection.FindOne(ctx, bson.M{"productId": productID}).Decode(&product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}
