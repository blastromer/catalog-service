package attribute

import (
	// "context" // Import context package
	// "log"     // Import log package
	"net/http"
	"strconv"
	// "strings" // Add this line
	// "time"    // Import time package

	"micro-services/catalog/internal/database"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options" // Import options package
)

func GetAttributes(c *gin.Context) {
	// Create a find options object with a limit of 10
	findOptions := options.Find().SetLimit(10)

	// Slice to hold the results
	var attributes []Attribute

	// Perform the query
	cursor, err := database.Collection.Find(c, bson.M{}, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attributes"})
		return
	}
	defer cursor.Close(c)

	// Iterate through the cursor and decode the documents
	for cursor.Next(c) {
		var attribute Attribute
		if err := cursor.Decode(&attribute); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode attribute"})
			return
		}
		attributes = append(attributes, attribute)
	}

	// Check for cursor errors
	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor error"})
		return
	}

	// Return the result as JSON
	c.JSON(http.StatusOK, attributes)
}

func GetAttributeByAttributeId(c *gin.Context) {
	attributeIdStr := c.Param("attributeId")
	attributeId, err := strconv.Atoi(attributeIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var attribute Attribute
	err = database.Collection.FindOne(c, bson.M{"_id": attributeId}).Decode(&attribute)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
		}
		return
	}

	c.JSON(http.StatusOK, attribute)
}