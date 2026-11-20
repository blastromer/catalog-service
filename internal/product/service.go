package product

import (
	"context"
	"errors"
	"time"

	"micro-services/catalog/internal/platform/apierror"
	"micro-services/catalog/internal/platform/pagination"

	"go.mongodb.org/mongo-driver/mongo"
)

// Service contains product business logic.
type Service struct {
	repo         Repository
	queryTimeout time.Duration
	defaultLimit int64
	maxLimit     int64
}

// NewService constructs a product service.
func NewService(repo Repository, queryTimeout time.Duration, defaultLimit, maxLimit int64) *Service {
	return &Service{
		repo:         repo,
		queryTimeout: queryTimeout,
		defaultLimit: defaultLimit,
		maxLimit:     maxLimit,
	}
}

// ListProducts returns a paginated product list.
func (s *Service) ListProducts(ctx context.Context, pageRaw, limitRaw string) (ListResult, pagination.Params, error) {
	params, err := pagination.Parse(pageRaw, limitRaw, s.defaultLimit, s.maxLimit)
	if err != nil {
		return ListResult{}, pagination.Params{}, err
	}

	queryCtx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	result, err := s.repo.List(queryCtx, params.Skip, int64(params.Limit))
	if err != nil {
		return ListResult{}, pagination.Params{}, apierror.Internal("failed to list products", err)
	}
	if result.Items == nil {
		result.Items = []Product{}
	}
	return result, params, nil
}

// GetProduct returns a single product by business ID.
func (s *Service) GetProduct(ctx context.Context, productID int) (*Product, error) {
	if productID <= 0 {
		return nil, apierror.Validation("productId must be a positive integer")
	}

	queryCtx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	product, err := s.repo.GetByProductID(queryCtx, productID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apierror.NotFound("product not found")
		}
		return nil, apierror.Internal("failed to fetch product", err)
	}
	return product, nil
}
