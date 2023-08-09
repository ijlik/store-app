package service

import (
	"context"
	"github.com/ijlik/store-app/internal/adapter/repository"
	"github.com/ijlik/store-app/internal/business/domain"
	errpkg "github.com/ijlik/store-app/pkg/error"
	httppagination "github.com/ijlik/store-app/pkg/http/pagination"
	"github.com/stretchr/testify/mock"
)

type StoreServiceMock struct {
	Mock mock.Mock
	repo repository.StoreRepository
}

type StoreService interface {
	MockCreateStore(ctx context.Context, request *domain.StoreRequest) (*domain.Store, errpkg.ErrorService)
	MockGetStoreById(ctx context.Context, id string) (*domain.Store, errpkg.ErrorService)
	MockUpdateStore(ctx context.Context, request *domain.StoreRequest, id string) errpkg.ErrorService
	MockShowStoreProducts(ctx context.Context, pagination *httppagination.Pagination, searchAndFilter *domain.SearchAndFilterProduct, storeId string) ([]*domain.Product, int64, errpkg.ErrorService)
}

func NewStoreMockService(storeService StoreService, repo repository.StoreRepository) *MockStoreService {
	return &MockStoreService{
		storeService: storeService,
		repo:         repo,
	}
}

type MockStoreService struct {
	storeService StoreService
	repo         repository.StoreRepository
}

func (s *StoreServiceMock) MockCreateStore(ctx context.Context, request *domain.StoreRequest) (*domain.Store, errpkg.ErrorService) {
	_ = s.Mock.Called(request)
	if request == nil {
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			"empty request",
		)
	}
	storeReq := &repository.Store{
		Name:                 request.Name,
		Url:                  request.Url,
		Address:              request.Address,
		Phone:                request.Phone,
		OperationalTimeStart: request.OperationalTimeStart,
		OperationalTimeEnd:   request.OperationalTimeEnd,
	}

	store, err := s.repo.CreateStore(ctx, storeReq)
	if err != nil {
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	return StoreRes(store), nil
}

func (s *StoreServiceMock) MockGetStoreById(ctx context.Context, id string) (*domain.Store, errpkg.ErrorService) {
	_ = s.Mock.Called(id)
	if id == "" {
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			"empty request",
		)
	}

	store, err := s.repo.GetStoreById(ctx, id)
	if err != nil {
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	return StoreRes(store), nil
}

func (s *StoreServiceMock) MockUpdateStore(ctx context.Context, request *domain.StoreRequest, id string) errpkg.ErrorService {
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
			"empty request",
		)
	}
	storeReq := &repository.Store{
		ID:                   id,
		Name:                 request.Name,
		Url:                  request.Url,
		Address:              request.Address,
		Phone:                request.Phone,
		OperationalTimeStart: request.OperationalTimeStart,
		OperationalTimeEnd:   request.OperationalTimeEnd,
	}

	err := s.repo.UpdateStore(ctx, storeReq)
	if err != nil {
		return errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	return nil
}

func (s *StoreServiceMock) MockShowStoreProducts(ctx context.Context, pagination *httppagination.Pagination, searchAndFilter *domain.SearchAndFilterProduct, storeId string) ([]*domain.Product, int64, errpkg.ErrorService) {
	_ = s.Mock.Called(pagination, searchAndFilter, storeId)
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
	if storeId == "" {
		return nil, 0, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			"empty id",
		)
	}

	sfp := &repository.SearchFilterPagination{
		Limit:         pagination.Limit,
		Offset:        pagination.Offset,
		Search:        searchAndFilter.Search,
		SortBy:        searchAndFilter.SortBy,
		SortDirection: searchAndFilter.SortDirection,
	}

	products, err := s.repo.ListProductByStoreId(ctx, sfp, storeId)
	if err != nil {
		return nil, 0, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	count, err := s.repo.CountProductByStoreId(ctx, sfp, storeId)
	if err != nil {
		return nil, 0, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	store, err := s.repo.GetStoreById(ctx, storeId)
	if err != nil {
		return nil, 0, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	var result []*domain.Product

	for _, product := range products {
		if err == nil && store != nil {
			result = append(result, ProductRes(product, store))
		}
	}

	return result, count, nil
}
