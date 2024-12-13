package routes

import (
	"micro-services/catalog/internal/product"
	"micro-services/catalog/internal/category"

	"github.com/gin-gonic/gin"
)

// RegisterProductRoutes sets up the routes for product-related endpoints
func RegisterProductRoutes(r *gin.Engine) {
	r.GET("/api/products", product.GetAllProducts)
	r.GET("/api/products/:productId", product.GetProductByID)
}

// RegisterCategoryRoutes sets up the routes for product-related endpoints
func RegisterCategoryRoutes(r *gin.Engine) {
	r.GET("/api/categories", category.GetCategories)
	r.GET("/api/categories/:categoryId", category.GetCategoryByCategoryId)
}