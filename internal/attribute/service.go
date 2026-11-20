package attribute

import (
	"context"
	"errors"
	"time"

	"micro-services/catalog/internal/platform/apierror"
	"micro-services/catalog/internal/platform/pagination"

	"go.mongodb.org/mongo-driver/mongo"
)

// Service contains attribute business logic.
type Service struct {
	repo         Repository
	queryTimeout time.Duration
	defaultLimit int64
	maxLimit     int64
}

// NewService constructs an attribute service.
func NewService(repo Repository, queryTimeout time.Duration, defaultLimit, maxLimit int64) *Service {
	return &Service{
		repo:         repo,
		queryTimeout: queryTimeout,
		defaultLimit: defaultLimit,
		maxLimit:     maxLimit,
	}
}

// ListAttributes returns a paginated attribute list.
func (s *Service) ListAttributes(ctx context.Context, pageRaw, limitRaw string) (ListResult, pagination.Params, error) {
	params, err := pagination.Parse(pageRaw, limitRaw, s.defaultLimit, s.maxLimit)
	if err != nil {
		return ListResult{}, pagination.Params{}, err
	}

	queryCtx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	result, err := s.repo.List(queryCtx, params.Skip, int64(params.Limit))
	if err != nil {
		return ListResult{}, pagination.Params{}, apierror.Internal("failed to list attributes", err)
	}
	if result.Items == nil {
		result.Items = []Attribute{}
	}
	return result, params, nil
}

// GetAttribute returns a single attribute by ID.
func (s *Service) GetAttribute(ctx context.Context, attributeID int) (*Attribute, error) {
	if attributeID <= 0 {
		return nil, apierror.Validation("attributeId must be a positive integer")
	}

	queryCtx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	attribute, err := s.repo.GetByID(queryCtx, attributeID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apierror.NotFound("attribute not found")
		}
		return nil, apierror.Internal("failed to fetch attribute", err)
	}
	return attribute, nil
}
