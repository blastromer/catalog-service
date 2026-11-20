package product

import (
	"strconv"

	"micro-services/catalog/internal/platform/apierror"
	"micro-services/catalog/internal/platform/response"

	"github.com/gin-gonic/gin"
)

// Handler exposes product HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler constructs a product handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GetAllProducts handles GET /api/products.
func (h *Handler) GetAllProducts(c *gin.Context) {
	result, params, err := h.service.ListProducts(c.Request.Context(), c.Query("page"), c.Query("limit"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OKPaged(c, result.Items, params.Page, params.Limit, result.Total)
}

// GetProductByID handles GET /api/products/:productId.
func (h *Handler) GetProductByID(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("productId"))
	if err != nil {
		response.Fail(c, apierror.Validation("invalid productId"))
		return
	}

	product, err := h.service.GetProduct(c.Request.Context(), productID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, product)
}
