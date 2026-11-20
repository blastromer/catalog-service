package product

import (
	"context"
	"errors"
	"testing"
	"time"

	"micro-services/catalog/internal/platform/apierror"

	"go.mongodb.org/mongo-driver/mongo"
)

type stubRepo struct {
	listResult ListResult
	listErr    error
	product    *Product
	getErr     error
}

func (s *stubRepo) List(ctx context.Context, skip, limit int64) (ListResult, error) {
	return s.listResult, s.listErr
}

func (s *stubRepo) GetByProductID(ctx context.Context, productID int) (*Product, error) {
	return s.product, s.getErr
}

func TestListProducts_Validation(t *testing.T) {
	svc := NewService(&stubRepo{}, time.Second, 50, 200)
	_, _, err := svc.ListProducts(context.Background(), "0", "10")
	assertValidation(t, err)
}

func TestListProducts_Success(t *testing.T) {
	svc := NewService(&stubRepo{
		listResult: ListResult{
			Items: []Product{{ProductID: 1, Title: "Widget"}},
			Total: 1,
		},
	}, time.Second, 50, 200)

	result, params, err := svc.ListProducts(context.Background(), "1", "10")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if params.Page != 1 || params.Limit != 10 {
		t.Fatalf("unexpected pagination: %+v", params)
	}
	if len(result.Items) != 1 || result.Total != 1 {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestGetProduct_NotFound(t *testing.T) {
	svc := NewService(&stubRepo{getErr: mongo.ErrNoDocuments}, time.Second, 50, 200)
	_, err := svc.GetProduct(context.Background(), 99)
	apiErr, ok := apierror.As(err)
	if !ok || apiErr.Kind != apierror.KindNotFound {
		t.Fatalf("expected not_found, got %v", err)
	}
}

func TestGetProduct_InvalidID(t *testing.T) {
	svc := NewService(&stubRepo{}, time.Second, 50, 200)
	_, err := svc.GetProduct(context.Background(), 0)
	assertValidation(t, err)
}

func TestGetProduct_InternalError(t *testing.T) {
	svc := NewService(&stubRepo{getErr: errors.New("boom")}, time.Second, 50, 200)
	_, err := svc.GetProduct(context.Background(), 1)
	apiErr, ok := apierror.As(err)
	if !ok || apiErr.Kind != apierror.KindInternal {
		t.Fatalf("expected internal, got %v", err)
	}
}

func assertValidation(t *testing.T, err error) {
	t.Helper()
	apiErr, ok := apierror.As(err)
	if !ok || apiErr.Kind != apierror.KindValidation {
		t.Fatalf("expected validation error, got %v", err)
	}
}
