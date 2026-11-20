package attribute

import (
	"strconv"

	"micro-services/catalog/internal/platform/apierror"
	"micro-services/catalog/internal/platform/response"

	"github.com/gin-gonic/gin"
)

// Handler exposes attribute HTTP endpoints.
type Handler struct {
	service *Service
}

// NewHandler constructs an attribute handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GetAttributes handles GET /api/attributes.
func (h *Handler) GetAttributes(c *gin.Context) {
	result, params, err := h.service.ListAttributes(c.Request.Context(), c.Query("page"), c.Query("limit"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OKPaged(c, result.Items, params.Page, params.Limit, result.Total)
}

// GetAttributeByAttributeID handles GET /api/attributes/:attributeId.
func (h *Handler) GetAttributeByAttributeID(c *gin.Context) {
	attributeID, err := strconv.Atoi(c.Param("attributeId"))
	if err != nil {
		response.Fail(c, apierror.Validation("invalid attributeId"))
		return
	}

	attribute, err := h.service.GetAttribute(c.Request.Context(), attributeID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, attribute)
}
