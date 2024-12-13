package main

import (
	"log"
	"micro-services/catalog/internal/database"
	"micro-services/catalog/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the MongoDB connection
	if err := database.Init("categories"); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	// Initialize the Gin router
	r := gin.Default()

	// Register product-related routes
	routes.RegisterCategoryRoutes(r)

	// Start the server on port 8082
	if err := r.Run(":2002"); err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}
}
