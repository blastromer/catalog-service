package main

import (
	"log"
	"micro-services/catalog/internal/database"
	"micro-services/catalog/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the MongoDB connection
	if err := database.Init("products"); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	// Initialize the Gin router
	r := gin.Default()

	// Register product-related routes
	routes.RegisterProductRoutes(r)

	// Start the server on port 2001
	if err := r.Run(":2001"); err != nil {
		log.Fatalf("Server startup failed: %v", err)
	}
}
