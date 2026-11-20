package category

import (
	"context"
	"errors"
	"strconv"
	"time"

	"micro-services/catalog/internal/platform/apierror"
	"micro-services/catalog/internal/platform/pagination"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// allowedQueryFields restricts dynamic filters to known safe indexes.
var allowedQueryFields = map[string]func(string) (any, error){
	"categoryId": parseIntValue,
	"parent":     parseIntValue,
	"sort":       parseIntValue,
	"name":       parseStringValue,
	"domain":     parseStringValue,
	"meta.slug":  parseStringValue,
}

// Service contains category business logic.
type Service struct {
	repo         Repository
	queryTimeout time.Duration
	defaultLimit int64
	maxLimit     int64
}

// NewService constructs a category service.
func NewService(repo Repository, queryTimeout time.Duration, defaultLimit, maxLimit int64) *Service {
	return &Service{
		repo:         repo,
		queryTimeout: queryTimeout,
		defaultLimit: defaultLimit,
		maxLimit:     maxLimit,
	}
}

// ListCategories returns a paginated category list.
func (s *Service) ListCategories(ctx context.Context, pageRaw, limitRaw string) (ListResult, pagination.Params, error) {
	params, err := pagination.Parse(pageRaw, limitRaw, s.defaultLimit, s.maxLimit)
	if err != nil {
		return ListResult{}, pagination.Params{}, err
	}

	queryCtx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	result, err := s.repo.List(queryCtx, params.Skip, int64(params.Limit))
	if err != nil {
		return ListResult{}, pagination.Params{}, apierror.Internal("failed to list categories", err)
	}
	if result.Items == nil {
		result.Items = []Category{}
	}
	return result, params, nil
}

// GetCategory returns a single category by business ID.
func (s *Service) GetCategory(ctx context.Context, categoryID int) (*Category, error) {
	if categoryID <= 0 {
		return nil, apierror.Validation("categoryId must be a positive integer")
	}

	queryCtx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	category, err := s.repo.GetByCategoryID(queryCtx, categoryID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apierror.NotFound("category not found")
		}
		return nil, apierror.Internal("failed to fetch category", err)
	}
	return category, nil
}

// QueryCategories runs a whitelisted field filter against categories.
func (s *Service) QueryCategories(ctx context.Context, field, value string) ([]Category, error) {
	if field == "" || value == "" {
		return nil, apierror.Validation("both field and value query parameters are required")
	}

	parser, ok := allowedQueryFields[field]
	if !ok {
		return nil, apierror.Validation("unsupported query field")
	}

	parsedValue, err := parser(value)
	if err != nil {
		return nil, apierror.Validation("invalid value for query field")
	}

	queryCtx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	categories, err := s.repo.FindByFilter(queryCtx, bson.M{field: parsedValue})
	if err != nil {
		return nil, apierror.Internal("failed to query categories", err)
	}
	if len(categories) == 0 {
		return nil, apierror.NotFound("no categories found for the given filter")
	}
	return categories, nil
}

func parseIntValue(raw string) (any, error) {
	return strconv.Atoi(raw)
}

func parseStringValue(raw string) (any, error) {
	return raw, nil
}
