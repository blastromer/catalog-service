package product

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// Dimensions represents the dimensions of the image
type Dimensions struct {
    Width  int `bson:"width" json:"width"`
    Height int `bson:"height" json:"height"`
}

// Image represents the image schema
type Image struct {
    ID         string             `bson:"id" json:"id"`
    Name       string             `bson:"name" json:"name"`
    Sort       int                `bson:"sort" json:"sort"`
    Dimensions Dimensions         `bson:"dimensions" json:"dimensions"`
    Path       string             `bson:"path" json:"path"`
    FileSize   string             `bson:"file_size" json:"file_size"`
    UploadedAt primitive.DateTime `bson:"uploaded_at,omitempty" json:"uploaded_at,omitempty"`
}

// Product represents the product schema for MongoDB
type Product struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    ProductID   int                `bson:"productId" json:"productId"`
    Title       string             `bson:"title" json:"title"`
    Description string             `bson:"description" json:"description"`
    UPC         string             `bson:"upc" json:"upc"`
    MPN         string             `bson:"mpn" json:"mpn"`
    Brand       string             `bson:"brand" json:"brand"`
    Images      []Image            `bson:"images" json:"images"`
    CreatedAt   primitive.DateTime `bson:"created_at" json:"created_at"`
    UpdatedAt   primitive.DateTime `bson:"updated_at" json:"updated_at"`
}
