package product

import (
	"context" // Import context package
	"log"     // Import log package
	"net/http"
	"strconv"
	"time" // Import time package

	"micro-services/catalog/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options" // Import options package
	"github.com/gin-gonic/gin"
)

func GetAllProducts(c *gin.Context) {
	// Retrieve the page number and limit from query parameters (default limit = 100)
	pageStr		:= c.DefaultQuery("page", "1")
	limitStr	:= c.DefaultQuery("limit", "100")

	page, err 	:= strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	limit, err 	:= strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit number"})
		return
	}

	// Skip documents based on page number and limit
	skip 		:= (page - 1) * limit

	// Create context with timeout to avoid long requests
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Fetch products from the database, using skip and limit for pagination
	var products []Product
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
	if err := cursor.All(ctx, &products); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode products"})
		return
	}

	// Handle case where no products are found
	if len(products) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No products found"})
		return
	}

	// Return the fetched products
	c.JSON(http.StatusOK, products)
}

// GetProductByID handles fetching a product by ID
func GetProductByID(c *gin.Context) {
    productIDStr := c.Param("productId")
    productID, err := strconv.Atoi(productIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
        return
    }

    var product Product
    err = database.Collection.FindOne(c, bson.M{"productId": productID}).Decode(&product)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
        }
        return
    }

    c.JSON(http.StatusOK, product)
}

