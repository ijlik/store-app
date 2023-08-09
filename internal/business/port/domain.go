package port

import (
	"context"
	"github.com/ijlik/store-app/internal/business/domain"
	errpkg "github.com/ijlik/store-app/pkg/error"
	httppagination "github.com/ijlik/store-app/pkg/http/pagination"
)

type StoreDomainService interface {
	CreateStore(ctx context.Context, request *domain.StoreRequest) (*domain.Store, errpkg.ErrorService)
	GetStoreById(ctx context.Context, id string) (*domain.Store, errpkg.ErrorService)
	UpdateStore(ctx context.Context, request *domain.StoreRequest, id string) errpkg.ErrorService
	ShowStoreProducts(ctx context.Context, pagination *httppagination.Pagination, searchAndFilter *domain.SearchAndFilterProduct, id string) errpkg.ErrorService

	ShowProducts(ctx context.Context, pagination *httppagination.Pagination, searchAndFilter *domain.SearchAndFilterProduct) errpkg.ErrorService
	CreateProduct(ctx context.Context, request *domain.ProductRequest) (*domain.Product, errpkg.ErrorService)
	GetProductByUrl(ctx context.Context, url string) (*domain.Product, errpkg.ErrorService)
	UpdateProduct(ctx context.Context, request *domain.ProductRequest, id string) errpkg.ErrorService
	DeleteProduct(ctx context.Context, id string) errpkg.ErrorService
}
