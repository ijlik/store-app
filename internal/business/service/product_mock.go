package service

import (
	"context"
	"database/sql"
	"github.com/ijlik/store-app/internal/adapter/repository"
	"github.com/ijlik/store-app/internal/business/domain"
	errpkg "github.com/ijlik/store-app/pkg/error"
	httppagination "github.com/ijlik/store-app/pkg/http/pagination"
	"github.com/stretchr/testify/mock"
	"time"
)

type ProductServiceMock struct {
	Mock mock.Mock
	repo repository.StoreRepository
}

type ProductService interface {
	MockShowProducts(ctx context.Context, pagination *httppagination.Pagination, searchAndFilter *domain.SearchAndFilterProduct) ([]*domain.Product, int64, errpkg.ErrorService)
	MockCreateProduct(ctx context.Context, request *domain.ProductRequest) (*domain.Product, errpkg.ErrorService)
	MockGetProductByUrl(ctx context.Context, url string) (*domain.Product, errpkg.ErrorService)
	MockUpdateProduct(ctx context.Context, request *domain.ProductRequest, id string) errpkg.ErrorService
	MockDeleteProduct(ctx context.Context, id string) errpkg.ErrorService
}

func NewProductService(productService ProductService, repo repository.StoreRepository) *MockProductService {
	return &MockProductService{
		productService: productService,
		repo:           repo,
	}
}

type MockProductService struct {
	productService ProductService
	repo           repository.StoreRepository
}

func (s *ProductServiceMock) MockShowProducts(ctx context.Context, pagination *httppagination.Pagination, searchAndFilter *domain.SearchAndFilterProduct) ([]*domain.Product, int64, errpkg.ErrorService) {
	_ = s.Mock.Called(pagination, searchAndFilter)
	if pagination == nil {
		return nil, 0, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			"empty pagination",
		)
	}
	if searchAndFilter == nil {
		return nil, 0, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			"empty searchAndFilter",
		)
	}
	sfp := &repository.SearchFilterPagination{
		Limit:         pagination.Limit,
		Offset:        pagination.Offset,
		Search:        searchAndFilter.Search,
		SortBy:        searchAndFilter.SortBy,
		SortDirection: searchAndFilter.SortDirection,
	}

	products, err := s.repo.ListProduct(ctx, sfp)
	if err != nil {
		return nil, 0, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	count, err := s.repo.CountProduct(ctx, sfp)
	if err != nil {
		return nil, 0, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	var result []*domain.Product

	for _, product := range products {
		store, err := s.repo.GetStoreById(ctx, product.StoreID)
		if err == nil && store != nil {
			result = append(result, ProductRes(product, store))
		}
	}

	return result, count, nil
}

func (s *ProductServiceMock) MockCreateProduct(ctx context.Context, request *domain.ProductRequest) (*domain.Product, errpkg.ErrorService) {
	_ = s.Mock.Called(request)
	if request == nil {
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			"empty request",
		)
	}
	productReq := &repository.Product{
		Name:        request.Name,
		Description: request.Description,
		Url:         request.Url,
		Price:       request.Price,
		StoreID:     request.StoreID,
		CreatedAt:   time.Now(),
		UpdatedAt: sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		},
	}

	product, err := s.repo.CreateProduct(ctx, productReq)
	if err != nil {
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	store, err := s.repo.GetStoreById(ctx, product.StoreID)
	if err != nil {
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	return ProductRes(product, store), nil
}

func (s *ProductServiceMock) MockGetProductByUrl(ctx context.Context, url string) (*domain.Product, errpkg.ErrorService) {
	_ = s.Mock.Called(url)
	if url == "" {
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			"empty url",
		)
	}

	product, err := s.repo.GetProductByUrl(ctx, url)
	if err != nil {
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	store, err := s.repo.GetStoreById(ctx, product.StoreID)
	if err != nil {
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	return ProductRes(product, store), nil
}

func (s *ProductServiceMock) MockUpdateProduct(ctx context.Context, request *domain.ProductRequest, id string) errpkg.ErrorService {
	_ = s.Mock.Called(request, id)
	if request == nil {
		return errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			"empty request",
		)
	}
	if id == "" {
		return errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			"empty id",
		)
	}
	productReq := &repository.Product{
		ID:          id,
		Name:        request.Name,
		Description: request.Description,
		Url:         request.Url,
		Price:       request.Price,
		StoreID:     request.StoreID,
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}

	err := s.repo.UpdateProduct(ctx, productReq)
	if err != nil {
		return errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	return nil
}

func (s *ProductServiceMock) MockDeleteProduct(ctx context.Context, id string) errpkg.ErrorService {
	_ = s.Mock.Called(id)
	if id == "" {
		return errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			"empty id",
		)
	}

	err := s.repo.DeleteProduct(ctx, id)
	if err != nil {
		return errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	return nil
}
