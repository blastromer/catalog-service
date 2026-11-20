package routes

import (
	"micro-services/catalog/internal/attribute"
	"micro-services/catalog/internal/category"
	"micro-services/catalog/internal/product"

	"github.com/gin-gonic/gin"
)

// RegisterProductRoutes sets up product endpoints.
func RegisterProductRoutes(r *gin.Engine, handler *product.Handler) {
	api := r.Group("/api")
	{
		api.GET("/products", handler.GetAllProducts)
		api.GET("/products/:productId", handler.GetProductByID)
	}
}

// RegisterCategoryRoutes sets up category endpoints.
// Static paths are registered before parameterized paths so /query is reachable.
func RegisterCategoryRoutes(r *gin.Engine, handler *category.Handler) {
	api := r.Group("/api")
	{
		api.GET("/categories", handler.GetCategories)
		api.GET("/categories/query", handler.GetCategorySmartQuery)
		api.GET("/categories/:categoryId", handler.GetCategoryByCategoryID)
	}
}

// RegisterAttributeRoutes sets up attribute endpoints.
func RegisterAttributeRoutes(r *gin.Engine, handler *attribute.Handler) {
	api := r.Group("/api")
	{
		api.GET("/attributes", handler.GetAttributes)
		api.GET("/attributes/:attributeId", handler.GetAttributeByAttributeID)
	}
}
