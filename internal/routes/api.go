package routes

import (
	"micro-services/catalog/internal/attribute"
	"micro-services/catalog/internal/category"
	"micro-services/catalog/internal/product"

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
	r.GET("/api/categories/query", category.GetCategorySmartQuery)
}

// RegisterAttributeRoutes sets up the routes for product-attribute endpoints
func RegisterAttributeRoutes(r *gin.Engine) {
	r.GET("/api/attributes", attribute.GetAttributes)
	r.GET("/api/attributes/:attributeId", attribute.GetAttributeByAttributeId)
	// r.GET("/api/categories/query", category.GetCategorySmartQuery)
}