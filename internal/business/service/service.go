package service

import (
	configdata "github.com/ijlik/store-app/pkg/config/data"
	// business package
	"github.com/ijlik/store-app/internal/adapter/repository"
	"github.com/ijlik/store-app/internal/business/port"
)

type service struct {
	repo   repository.StoreRepository
	config configdata.Config
}

func NewStoreService(
	repo repository.StoreRepository,
	config configdata.Config,
) port.StoreDomainService {
	return &service{
		repo,
		config,
	}
}
