package category

import (
	"strconv"

	"micro-services/catalog/internal/platform/apierror"
	"micro-services/catalog/internal/platform/response"

	"github.com/gin-gonic/gin"
)

// Handler exposes category HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler constructs a category handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GetCategories handles GET /api/categories.
func (h *Handler) GetCategories(c *gin.Context) {
	result, params, err := h.service.ListCategories(c.Request.Context(), c.Query("page"), c.Query("limit"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OKPaged(c, result.Items, params.Page, params.Limit, result.Total)
}

// GetCategoryByCategoryID handles GET /api/categories/:categoryId.
func (h *Handler) GetCategoryByCategoryID(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("categoryId"))
	if err != nil {
		response.Fail(c, apierror.Validation("invalid categoryId"))
		return
	}

	category, err := h.service.GetCategory(c.Request.Context(), categoryID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, category)
}

// GetCategorySmartQuery handles GET /api/categories/query.
func (h *Handler) GetCategorySmartQuery(c *gin.Context) {
	categories, err := h.service.QueryCategories(c.Request.Context(), c.Query("field"), c.Query("value"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, categories)
}
