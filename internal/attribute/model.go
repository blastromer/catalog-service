package attribute

import "go.mongodb.org/mongo-driver/bson/primitive"

// Value represents the nested "values" field in the document
type Value struct {
    ID        		int    				`bson:"id" json:"id"`
    Name      		string 				`bson:"name" json:"name"`
    SortOrder 		int    				`bson:"sortOrder" json:"sortOrder"`
}

// Attribute represents the root-level structure of the document
type Attribute struct {
    ID             	int                	`bson:"_id" json:"id"`
    Name           	string             	`bson:"name" json:"name"`
    SortOrder      	int                	`bson:"sortOrder" json:"sortOrder"`
    Values         	[]Value            	`bson:"values" json:"values"`
    AttributeValue 	string             	`bson:"attributeValue" json:"attributeValue"`
    Note           	string             	`bson:"note" json:"note"`
    CreatedAt      	primitive.DateTime 	`bson:"created_at" json:"createdAt"`
    UpdatedAt      	primitive.DateTime 	`bson:"updated_at" json:"updatedAt"`
}