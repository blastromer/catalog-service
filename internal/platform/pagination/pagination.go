package pagination

import (
	"strconv"

	"micro-services/catalog/internal/platform/apierror"
)

// Params holds normalized pagination input.
type Params struct {
	Page  int
	Limit int
	Skip  int64
}

// Parse reads page/limit query values and applies service defaults and caps.
func Parse(pageRaw, limitRaw string, defaultLimit, maxLimit int64) (Params, error) {
	page := 1
	if pageRaw != "" {
		parsed, err := strconv.Atoi(pageRaw)
		if err != nil || parsed <= 0 {
			return Params{}, apierror.Validation("page must be a positive integer")
		}
		page = parsed
	}

	limit := int(defaultLimit)
	if limitRaw != "" {
		parsed, err := strconv.Atoi(limitRaw)
		if err != nil || parsed <= 0 {
			return Params{}, apierror.Validation("limit must be a positive integer")
		}
		limit = parsed
	}

	if int64(limit) > maxLimit {
		return Params{}, apierror.Validation("limit exceeds maximum allowed value")
	}

	return Params{
		Page:  page,
		Limit: limit,
		Skip:  int64((page - 1) * limit),
	}, nil
}
