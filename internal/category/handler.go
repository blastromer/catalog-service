package category

import (
	"context" // Import context package
	"log"     // Import log package
	"net/http"
	"strconv"
	"strings" // Add this line
	"time"    // Import time package

	"micro-services/catalog/internal/database"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options" // Import options package
)

func GetCategories(c *gin.Context) {
	// Retrieve the page number and limit from query parameters (default limit = 100)
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "100")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit number"})
		return
	}

	// Skip documents based on page number and limit
	skip := (page - 1) * limit

	// Create context with timeout to avoid long requests
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Fetch products from the database, using skip and limit for pagination
	var categories []Category
	cursor, err := database.Collection.Find(ctx, bson.M{}, options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))) // Ensure options is used here
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Printf("Failed to close cursor: %v", err)
		}
	}()

	// Decode the results into the products slice
	if err := cursor.All(ctx, &categories); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode products"})
		return
	}

	// Handle case where no products are found
	if len(categories) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No products found"})
		return
	}

	// Return the fetched products
	c.JSON(http.StatusOK, categories)
}

// GetProductByID handles fetching a product by ID
func GetCategoryByCategoryId(c *gin.Context) {
	categoryIdStr := c.Param("categoryId")
	categoryId, err := strconv.Atoi(categoryIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var category Category
	err = database.Collection.FindOne(c, bson.M{"categoryId": categoryId}).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
		}
		return
	}

	c.JSON(http.StatusOK, category)
}

func GetCategorySmartQuery(c *gin.Context) {
	queryField := c.DefaultQuery("field", "")
	queryValue := c.DefaultQuery("value", "")

	// Log the input parameters
	log.Printf("Query parameters - field: %s, value: %s", queryField, queryValue)

	// Start with an empty filter
	filter := bson.M{}

	if queryField != "" && queryValue != "" {
		// Handle nested fields (like meta.slug)
		if strings.Contains(queryField, ".") {
			// For nested fields, we don't try to convert to int
			filter[queryField] = queryValue
		} else {
			// For non-nested fields, try to convert to int if possible
			if intValue, err := strconv.Atoi(queryValue); err == nil {
				filter[queryField] = intValue
			} else {
				filter[queryField] = queryValue
			}
		}
	}

	// Add debug logging
	log.Printf("Searching with filter: %+v", filter)
	var categories []Category
	cursor, err := database.Collection.Find(c, filter)
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}
	defer cursor.Close(c)

	for cursor.Next(c) {
		var category Category
		if err := cursor.Decode(&category); err != nil {
			log.Printf("Error decoding category: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode category"})
			return
		}
		categories = append(categories, category)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor iteration error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cursor iteration error"})
		return
	}

	if len(categories) == 0 {
		log.Printf("No documents found for filter: %+v", filter)
		c.JSON(http.StatusNotFound, gin.H{"error": "No categories found"})
		return
	}

	c.JSON(http.StatusOK, categories)
}
