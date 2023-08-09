package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/ijlik/store-app/internal/adapter/repository"
	"github.com/ijlik/store-app/internal/business/domain"
	errpkg "github.com/ijlik/store-app/pkg/error"
	httppagination "github.com/ijlik/store-app/pkg/http/pagination"
	"sync"
	"sync/atomic"
	"time"
)

func (s *service) CreateStore(ctx context.Context, request *domain.StoreRequest) (*domain.Store, errpkg.ErrorService) {
	store, err := s.repo.CreateStore(ctx, &repository.Store{
		ID:                   uuid.New().String(),
		Name:                 request.Name,
		Url:                  request.Url,
		Address:              request.Address,
		Phone:                request.Phone,
		OperationalTimeStart: request.OperationalTimeStart,
		OperationalTimeEnd:   request.OperationalTimeEnd,
		CreatedAt:            time.Now().UTC(),
	})
	if err != nil {
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}
	return &domain.Store{
		ID:                   store.ID,
		Name:                 store.Name,
		Url:                  store.Url,
		Address:              store.Address,
		Phone:                store.Phone,
		OperationalTimeStart: store.OperationalTimeStart,
		OperationalTimeEnd:   store.OperationalTimeEnd,
		CreatedAt:            store.CreatedAt,
	}, nil
}

func (s *service) GetStoreById(ctx context.Context, id string) (*domain.Store, errpkg.ErrorService) {
	store, err := s.repo.GetStoreById(ctx, id)
	if err != nil {
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}
	if store == nil {
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrNotFound,
			"store not found",
		)
	}
	return &domain.Store{
		ID:                   store.ID,
		Name:                 store.Name,
		Url:                  store.Url,
		Address:              store.Address,
		Phone:                store.Phone,
		OperationalTimeStart: store.OperationalTimeStart,
		OperationalTimeEnd:   store.OperationalTimeEnd,
		CreatedAt:            store.CreatedAt,
	}, nil
}

func (s *service) UpdateStore(ctx context.Context, request *domain.StoreRequest, id string) errpkg.ErrorService {
	err := s.repo.UpdateStore(ctx, &repository.Store{
		ID:                   id,
		Name:                 request.Name,
		Url:                  request.Url,
		Address:              request.Address,
		Phone:                request.Phone,
		OperationalTimeStart: request.OperationalTimeStart,
		OperationalTimeEnd:   request.OperationalTimeEnd,
	})
	if err != nil {
		return errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}
	return nil
}

func (s *service) ShowStoreProducts(ctx context.Context, pagination *httppagination.Pagination, searchAndFilter *domain.SearchAndFilterProduct, id string) errpkg.ErrorService {
	var (
		g           sync.WaitGroup
		int64Atomic atomic.Int64
		arrayAtomic atomic.Value
		errAtomic   atomic.Value
		result      []*domain.Product
	)

	sfp := &repository.SearchFilterPagination{
		Limit:         pagination.Limit,
		Offset:        pagination.Offset,
		Search:        searchAndFilter.Search,
		SortBy:        searchAndFilter.SortBy,
		SortDirection: searchAndFilter.SortDirection,
	}

	store, err := s.repo.GetStoreById(ctx, id)
	if err != nil {
		return errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}
	if store == nil {
		return errpkg.DefaultServiceError(
			errpkg.ErrNotFound,
			"store not found",
		)
	}

	g.Add(1)
	go func() {
		defer g.Done()
		products, err := s.repo.ListProductByStoreId(ctx, sfp, id)
		if err != nil {
			errAtomic.Store(err)
		} else {
			arrayAtomic.Store(products)
		}
	}()

	g.Add(1)
	go func() {
		defer g.Done()
		count, err := s.repo.CountProductByStoreId(ctx, sfp, id)
		if err != nil {
			errAtomic.Store(err)
		} else {
			int64Atomic.Store(count)
		}
	}()
	g.Wait()

	if err, ok := errAtomic.Load().(error); ok {
		return errpkg.DefaultServiceError(errpkg.ErrInternal, err.Error())
	}

	if products, ok := arrayAtomic.Load().([]*repository.Product); !ok {
		return errpkg.DefaultServiceError(errpkg.ErrInternal, "")
	} else {
		for _, product := range products {
			result = append(result, ProductRes(product, store))
		}
	}

	pagination.SetData(result, int64Atomic.Load())
	return nil
}
