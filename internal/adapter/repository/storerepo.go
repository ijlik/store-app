package repository

import (
	"context"
	"database/sql"
	"time"
)

const createStoreQuery = `INSERT INTO stores (name, url, address, phone, operational_time_start, operational_time_end, created_at) VALUES ($1, $2, $3, $4, $5, $6, CURRENT_TIMESTAMP) RETURNING id`

func (r *repo) CreateStore(ctx context.Context, req *Store) (*Store, error) {
	var id string
	if err := r.conn.QueryRowContext(
		ctx,
		createStoreQuery,
		req.RowDataCreate()...,
	).Scan(&id); err != nil {
		return nil, err
	}

	return &Store{
		ID:                   id,
		Name:                 req.Name,
		Url:                  req.Url,
		Address:              req.Address,
		Phone:                req.Phone,
		OperationalTimeStart: req.OperationalTimeStart,
		OperationalTimeEnd:   req.OperationalTimeEnd,
		CreatedAt:            time.Now().UTC(),
		UpdatedAt:            sql.NullTime{},
	}, nil
}

const getStoreByIdQuery = `SELECT id, name, url, address, phone, operational_time_start, operational_time_end, created_at, updated_at FROM stores WHERE id = $1 LIMIT 1`

func (r *repo) GetStoreById(ctx context.Context, id string) (*Store, error) {
	var data Store
	err := r.conn.GetContext(
		ctx,
		&data,
		getStoreByIdQuery,
		id,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &data, nil
}

const updateStoreQuery = `UPDATE stores SET name = $2, url = $3, address = $4, phone = $5, operational_time_start = $6, operational_time_end = $7, updated_at = CURRENT_TIMESTAMP WHERE id = $1`

func (r *repo) UpdateStore(ctx context.Context, req *Store) error {
	if _, err := r.conn.ExecContext(
		ctx,
		updateStoreQuery,
		req.RowDataUpdate()...,
	); err != nil {
		return err
	}

	return nil
}
