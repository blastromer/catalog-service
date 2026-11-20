package category

import (
	"context"
	"testing"
	"time"

	"micro-services/catalog/internal/platform/apierror"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type stubRepo struct {
	listResult ListResult
	listErr    error
	category   *Category
	getErr     error
	queryItems []Category
	queryErr   error
	lastFilter bson.M
}

func (s *stubRepo) List(ctx context.Context, skip, limit int64) (ListResult, error) {
	return s.listResult, s.listErr
}

func (s *stubRepo) GetByCategoryID(ctx context.Context, categoryID int) (*Category, error) {
	return s.category, s.getErr
}

func (s *stubRepo) FindByFilter(ctx context.Context, filter bson.M) ([]Category, error) {
	s.lastFilter = filter
	return s.queryItems, s.queryErr
}

func TestQueryCategories_RejectsUnknownField(t *testing.T) {
	svc := NewService(&stubRepo{}, time.Second, 50, 200)
	_, err := svc.QueryCategories(context.Background(), "$where", "1")
	apiErr, ok := apierror.As(err)
	if !ok || apiErr.Kind != apierror.KindValidation {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestQueryCategories_ParsesIntegerFields(t *testing.T) {
	repo := &stubRepo{queryItems: []Category{{CategoryID: 12, Name: "Brakes"}}}
	svc := NewService(repo, time.Second, 50, 200)

	items, err := svc.QueryCategories(context.Background(), "categoryId", "12")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if repo.lastFilter["categoryId"] != 12 {
		t.Fatalf("expected int filter value, got %#v", repo.lastFilter["categoryId"])
	}
}

func TestQueryCategories_RequiresParams(t *testing.T) {
	svc := NewService(&stubRepo{}, time.Second, 50, 200)
	_, err := svc.QueryCategories(context.Background(), "", "")
	apiErr, ok := apierror.As(err)
	if !ok || apiErr.Kind != apierror.KindValidation {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestGetCategory_NotFound(t *testing.T) {
	svc := NewService(&stubRepo{getErr: mongo.ErrNoDocuments}, time.Second, 50, 200)
	_, err := svc.GetCategory(context.Background(), 5)
	apiErr, ok := apierror.As(err)
	if !ok || apiErr.Kind != apierror.KindNotFound {
		t.Fatalf("expected not_found, got %v", err)
	}
}
