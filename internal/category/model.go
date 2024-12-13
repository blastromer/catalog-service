package category

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)
// Meta represents the category schema for MongoDB
type Meta struct {
	Title		string				`bson:"title"			json:"title"`
	Description	string				`bson:"desciption"		json:"description"`
	Slug		string				`bson:"slug"			json:"slug"`
}

// Product represents the product schema for MongoDB
type Category struct {
	ID          primitive.ObjectID 	`bson:"_id,omitempty" 	json:"id"`
	Domain   	string              `bson:"domain" 			json:"domain"`
	Name       	string             	`bson:"name" 			json:"name"`
	Meta		Meta				`bson:"meta"			json:"meta"`
	Sort		int					`bson:"sort"			json:"sort"`
	Parent		int					`bson:"parent"			json:"parent"`
	Content		string				`bson:"content"			json:"content"`
	CategoryId	int					`bson:"categoryId"		json:"categoryId"`
	Image		string				`bson:"image"			json:"image"`
	Brands		[]int				`bson:"brands"			json:"brands"`
}
