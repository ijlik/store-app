package repository

import "context"

type StoreRepository interface {
	StoreRepo
	ProductRepo
}

type StoreRepo interface {
	CreateStore(ctx context.Context, req *Store) (*Store, error)
	GetStoreById(ctx context.Context, id string) (*Store, error)
	UpdateStore(ctx context.Context, req *Store) error
}

type ProductRepo interface {
	CountProduct(ctx context.Context, sfp *SearchFilterPagination) (int64, error)
	CountProductByStoreId(ctx context.Context, sfp *SearchFilterPagination, storeId string) (int64, error)
	ListProduct(ctx context.Context, sfp *SearchFilterPagination) ([]*Product, error)
	ListProductByStoreId(ctx context.Context, sfp *SearchFilterPagination, storeId string) ([]*Product, error)
	CreateProduct(ctx context.Context, req *Product) (*Product, error)
	GetProductById(ctx context.Context, id string) (*Product, error)
	GetProductByUrl(ctx context.Context, url string) (*Product, error)
	UpdateProduct(ctx context.Context, req *Product) error
	DeleteProduct(ctx context.Context, id string) error
}
