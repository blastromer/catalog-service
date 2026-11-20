package main

import (
	"context"
	"log"
	"time"

	"micro-services/catalog/internal/config"
	"micro-services/catalog/internal/database"
	"micro-services/catalog/internal/platform/server"
	"micro-services/catalog/internal/product"
	"micro-services/catalog/internal/routes"
)

func main() {
	cfg, err := config.Load("product-service", ":2001", "products")
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	store, err := database.Connect(ctx, cfg.MongoURI, cfg.MongoDatabase, cfg.CollectionName)
	if err != nil {
		log.Fatalf("database error: %v", err)
	}

	deps := server.Dependencies{Config: cfg, Store: store}
	router := server.NewRouter(deps)

	repo := product.NewRepository(store.Collection)
	service := product.NewService(repo, cfg.QueryTimeout, cfg.DefaultLimit, cfg.MaxLimit)
	handler := product.NewHandler(service)
	routes.RegisterProductRoutes(router, handler)

	if err := server.Run(deps, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
