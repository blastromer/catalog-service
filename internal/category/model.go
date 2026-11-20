package category

import "go.mongodb.org/mongo-driver/bson/primitive"

// Meta holds SEO and routing metadata for a category.
type Meta struct {
	Title       string `bson:"title" json:"title"`
	Description string `bson:"desciption" json:"description"` // legacy field spelling preserved for compatibility
	Slug        string `bson:"slug" json:"slug"`
}

// Category is the catalog category document.
type Category struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Domain     string             `bson:"domain" json:"domain"`
	Name       string             `bson:"name" json:"name"`
	Meta       Meta               `bson:"meta" json:"meta"`
	Sort       int                `bson:"sort" json:"sort"`
	Parent     int                `bson:"parent" json:"parent"`
	Content    string             `bson:"content" json:"content"`
	CategoryID int                `bson:"categoryId" json:"categoryId"`
	Image      string             `bson:"image" json:"image"`
	Brands     []int              `bson:"brands" json:"brands"`
}

// ListResult is a paginated category collection.
type ListResult struct {
	Items []Category
	Total int64
}

// QueryFilter is a validated field/value search.
type QueryFilter struct {
	Field string
	Value any
}
