package repository

import (
	"github.com/jmoiron/sqlx"
)

type repo struct {
	conn *sqlx.DB
}

func NewStoreRepo(db *sqlx.DB) StoreRepository {
	return &repo{db}
}
