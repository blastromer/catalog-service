package main

import (
	"context"
	"log"
	"time"

	"micro-services/catalog/internal/category"
	"micro-services/catalog/internal/config"
	"micro-services/catalog/internal/database"
	"micro-services/catalog/internal/platform/server"
	"micro-services/catalog/internal/routes"
)

func main() {
	cfg, err := config.Load("category-service", ":2002", "categories")
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

	repo := category.NewRepository(store.Collection)
	service := category.NewService(repo, cfg.QueryTimeout, cfg.DefaultLimit, cfg.MaxLimit)
	handler := category.NewHandler(service)
	routes.RegisterCategoryRoutes(router, handler)

	if err := server.Run(deps, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
