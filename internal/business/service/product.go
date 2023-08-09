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

func (s *service) ShowProducts(ctx context.Context, pagination *httppagination.Pagination, searchAndFilter *domain.SearchAndFilterProduct) errpkg.ErrorService {
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

	g.Add(1)
	go func() {
		defer g.Done()
		products, err := s.repo.ListProduct(ctx, sfp)
		if err != nil {
			errAtomic.Store(err)
		} else {
			arrayAtomic.Store(products)
		}
	}()

	g.Add(1)
	go func() {
		defer g.Done()
		count, err := s.repo.CountProduct(ctx, sfp)
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
			store, err := s.repo.GetStoreById(ctx, product.StoreID)
			if err == nil && store != nil {
				result = append(result, ProductRes(product, store))
			}
		}
	}

	pagination.SetData(result, int64Atomic.Load())
	return nil
}

func (s *service) CreateProduct(ctx context.Context, request *domain.ProductRequest) (*domain.Product, errpkg.ErrorService) {
	product, err := s.repo.CreateProduct(ctx, &repository.Product{
		ID:          uuid.New().String(),
		Name:        request.Name,
		Url:         request.Url,
		Price:       request.Price,
		StoreID:     request.StoreID,
		Description: request.Description,
		CreatedAt:   time.Now().UTC(),
	})
	if err != nil {
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	store, err := s.repo.GetStoreById(ctx, request.StoreID)
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

	return ProductRes(product, store), nil
}

func (s *service) GetProductByUrl(ctx context.Context, url string) (*domain.Product, errpkg.ErrorService) {
	product, err := s.repo.GetProductByUrl(ctx, url)
	if err != nil {
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}
	if product == nil {
		return nil, errpkg.DefaultServiceError(
			errpkg.ErrNotFound,
			"product not found",
		)
	}

	store, err := s.repo.GetStoreById(ctx, product.StoreID)
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

	return ProductRes(product, store), nil
}

func (s *service) UpdateProduct(ctx context.Context, request *domain.ProductRequest, id string) errpkg.ErrorService {
	product, err := s.repo.GetProductById(ctx, id)
	if err != nil {
		return errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}
	if product == nil {
		return errpkg.DefaultServiceError(
			errpkg.ErrNotFound,
			"product not found",
		)
	}

	err = s.repo.UpdateProduct(ctx, &repository.Product{
		ID:          product.ID,
		Name:        request.Name,
		Url:         product.Url,
		Price:       request.Price,
		StoreID:     request.StoreID,
		Description: request.Description,
	})
	if err != nil {
		return errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	return nil
}

func (s *service) DeleteProduct(ctx context.Context, id string) errpkg.ErrorService {
	product, err := s.repo.GetProductById(ctx, id)
	if err != nil {
		return errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}
	if product == nil {
		return errpkg.DefaultServiceError(
			errpkg.ErrNotFound,
			"product not found",
		)
	}

	err = s.repo.DeleteProduct(ctx, id)
	if err != nil {
		return errpkg.DefaultServiceError(
			errpkg.ErrInternal,
			err.Error(),
		)
	}

	return nil
}
