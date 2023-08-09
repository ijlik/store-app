package service

import (
	"github.com/ijlik/store-app/internal/adapter/repository"
	"github.com/ijlik/store-app/internal/business/domain"
)

func ProductRes(product *repository.Product, store *repository.Store) *domain.Product {
	return &domain.Product{
		ID:          product.ID,
		Name:        product.Name,
		Url:         product.Url,
		Price:       product.Price,
		Description: product.Description,
		Store:       StoreRes(store),
		CreatedAt:   product.CreatedAt,
	}
}

func StoreRes(store *repository.Store) *domain.Store {
	return &domain.Store{
		ID:                   store.ID,
		Name:                 store.Name,
		Url:                  store.Url,
		Address:              store.Address,
		Phone:                store.Phone,
		OperationalTimeStart: store.OperationalTimeStart,
		OperationalTimeEnd:   store.OperationalTimeEnd,
		CreatedAt:            store.CreatedAt,
	}

}
