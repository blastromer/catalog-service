package response

import (
	"net/http"

	"micro-services/catalog/internal/platform/apierror"

	"github.com/gin-gonic/gin"
)

// Envelope is the standard success payload shape.
type Envelope struct {
	Data       any       `json:"data"`
	Pagination *PageMeta `json:"pagination,omitempty"`
	Meta       *Meta     `json:"meta,omitempty"`
}

// PageMeta describes a paginated collection.
type PageMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"totalPages"`
}

// Meta carries optional response metadata.
type Meta struct {
	Message string `json:"message,omitempty"`
}

// ErrorBody is the standard error payload shape.
type ErrorBody struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail describes a single API error.
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Envelope{Data: data})
}

func OKPaged(c *gin.Context, data any, page, limit int, total int64) {
	totalPages := int64(0)
	if limit > 0 {
		totalPages = (total + int64(limit) - 1) / int64(limit)
	}
	c.JSON(http.StatusOK, Envelope{
		Data: data,
		Pagination: &PageMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

func Fail(c *gin.Context, err error) {
	if apiErr, ok := apierror.As(err); ok {
		status, code := mapKind(apiErr.Kind)
		c.JSON(status, ErrorBody{
			Error: ErrorDetail{Code: code, Message: apiErr.Message},
		})
		return
	}

	c.JSON(http.StatusInternalServerError, ErrorBody{
		Error: ErrorDetail{Code: "internal_error", Message: "unexpected server error"},
	})
}

func mapKind(kind apierror.Kind) (int, string) {
	switch kind {
	case apierror.KindValidation:
		return http.StatusBadRequest, "validation_error"
	case apierror.KindNotFound:
		return http.StatusNotFound, "not_found"
	case apierror.KindConflict:
		return http.StatusConflict, "conflict"
	default:
		return http.StatusInternalServerError, "internal_error"
	}
}
