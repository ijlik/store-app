package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

var countProductsQuery = `SELECT count(*) FROM products`

func (r *repo) CountProduct(ctx context.Context, sfp *SearchFilterPagination) (int64, error) {
	var (
		count  int64
		params []any
		query  = countProductsQuery
	)

	query, params, err := sfp.BuildWhere(query, false, "")
	if err != nil {
		return 0, err
	}

	if err := r.conn.QueryRowContext(
		ctx,
		query,
		params...,
	).Scan(&count); err != nil {
		return count, err
	}

	return count, nil
}

func (r *repo) CountProductByStoreId(ctx context.Context, sfp *SearchFilterPagination, storeId string) (int64, error) {
	var (
		count  int64
		params []any
		query  = countProductsQuery
	)

	query, params, err := sfp.BuildWhere(query, false, fmt.Sprintf("store_id = '%s'", storeId))
	if err != nil {
		return 0, err
	}

	if err := r.conn.QueryRowContext(
		ctx,
		query,
		params...,
	).Scan(&count); err != nil {
		return count, err
	}

	return count, nil
}

var listProductsQuery = `SELECT id, store_id, name, url, price, description, created_at, updated_at FROM products`

func (r *repo) ListProduct(ctx context.Context, sfp *SearchFilterPagination) ([]*Product, error) {
	var (
		data          []*Product
		params        []any
		usePagination bool
		query         = listProductsQuery
	)

	if sfp.Limit != 0 {
		usePagination = true
	}

	query, params, err := sfp.BuildWhere(query, usePagination, "")
	if err != nil {
		return nil, err
	}

	rows, err := r.conn.QueryContext(
		ctx,
		query,
		params...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e Product
		if err := rows.Scan(
			&e.ID,
			&e.StoreID,
			&e.Name,
			&e.Url,
			&e.Price,
			&e.Description,
			&e.CreatedAt,
			&e.UpdatedAt,
		); err != nil {
			return nil, err
		}
		data = append(data, &e)
	}

	return data, nil
}

func (r *repo) ListProductByStoreId(ctx context.Context, sfp *SearchFilterPagination, storeId string) ([]*Product, error) {
	var (
		data          []*Product
		params        []any
		usePagination bool
		query         = listProductsQuery
	)

	if sfp.Limit != 0 {
		usePagination = true
	}

	query, params, err := sfp.BuildWhere(query, usePagination, fmt.Sprintf("store_id = '%s'", storeId))
	if err != nil {
		return nil, err
	}

	rows, err := r.conn.QueryContext(
		ctx,
		query,
		params...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e Product
		if err := rows.Scan(
			&e.ID,
			&e.StoreID,
			&e.Name,
			&e.Url,
			&e.Price,
			&e.Description,
			&e.CreatedAt,
			&e.UpdatedAt,
		); err != nil {
			return nil, err
		}
		data = append(data, &e)
	}

	return data, nil
}

const createProductQuery = `INSERT INTO products (store_id, name, url, price, description, created_at) VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP) RETURNING id`

func (r *repo) CreateProduct(ctx context.Context, req *Product) (*Product, error) {
	var id string
	if err := r.conn.QueryRowContext(
		ctx,
		createProductQuery,
		req.RowDataCreate()...,
	).Scan(&id); err != nil {
		return nil, err
	}

	return &Product{
		ID:          id,
		StoreID:     req.StoreID,
		Name:        req.Name,
		Url:         req.Url,
		Price:       req.Price,
		Description: req.Description,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   sql.NullTime{},
	}, nil
}

const getProductByIdQuery = `SELECT id, store_id, name, url, price, description, created_at, updated_at FROM products WHERE id = $1 LIMIT 1`

func (r *repo) GetProductById(ctx context.Context, id string) (*Product, error) {
	var data Product
	err := r.conn.GetContext(
		ctx,
		&data,
		getProductByIdQuery,
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

const getProductByUrlQuery = `SELECT id, store_id, name, url, price, description, created_at, updated_at FROM products WHERE url = $1 LIMIT 1`

func (r *repo) GetProductByUrl(ctx context.Context, slug string) (*Product, error) {
	var data Product
	err := r.conn.GetContext(
		ctx,
		&data,
		getProductByUrlQuery,
		slug,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &data, nil
}

const updateProductQuery = `UPDATE products SET store_id = $2, name = $3, url = $4, price = $5, description = $6, updated_at = CURRENT_TIMESTAMP WHERE id = $1`

func (r *repo) UpdateProduct(ctx context.Context, req *Product) error {
	if _, err := r.conn.ExecContext(
		ctx,
		updateProductQuery,
		req.RowDataUpdate()...,
	); err != nil {
		return err
	}

	return nil
}

var deleteProductQuery = `DELETE FROM products WHERE id = $1`

func (r *repo) DeleteProduct(ctx context.Context, id string) error {
	if _, err := r.conn.ExecContext(
		ctx,
		deleteProductQuery,
		id,
	); err != nil {
		return err
	}

	return nil
}
