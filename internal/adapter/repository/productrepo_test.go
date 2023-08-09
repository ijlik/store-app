package repository

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCountProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewStoreRepo(dbx)

	expectedCount := int64(10)
	countProductsQueryMock := "SELECT count\\(\\*\\) FROM products"
	mock.ExpectQuery(countProductsQueryMock).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

	sfp := &SearchFilterPagination{
		Limit:         10,
		Offset:        0,
		Search:        "",
		SortBy:        "created_at",
		SortDirection: "DESC",
	}

	ctx := context.Background()
	result, err := repo.CountProduct(ctx, sfp)
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, result)
}

func TestCountProductByStoreId(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewStoreRepo(dbx)

	storeId := "test_store_id"
	expectedCount := int64(8)
	countProductsQueryMock := "SELECT count\\(\\*\\) FROM products"
	mock.ExpectQuery(countProductsQueryMock).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

	sfp := &SearchFilterPagination{
		Limit:         10,
		Offset:        0,
		Search:        "",
		SortBy:        "created_at",
		SortDirection: "DESC",
	}

	ctx := context.Background()
	result, err := repo.CountProductByStoreId(ctx, sfp, storeId)
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, result)
}

func TestListProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewStoreRepo(dbx)

	expectedData := []*Product{
		{
			ID:          "test_product_id",
			StoreID:     "test_store_id",
			Name:        "test_product_name",
			Url:         "test_product_url",
			Price:       100,
			Description: "test_product_description",
			CreatedAt:   time.Now(),
			UpdatedAt: sql.NullTime{
				Time:  time.Time{},
				Valid: false,
			},
		},
	}
	listProductsQueryMock := "SELECT id, store_id, name, url, price, description, created_at, updated_at FROM products"
	mock.ExpectQuery(listProductsQueryMock).WillReturnRows(sqlmock.NewRows([]string{"id", "store_id", "name", "url", "price", "description", "created_at", "updated_at"}).
		AddRow(expectedData[0].ID, expectedData[0].StoreID, expectedData[0].Name, expectedData[0].Url, expectedData[0].Price, expectedData[0].Description, expectedData[0].CreatedAt, expectedData[0].UpdatedAt))

	sfp := &SearchFilterPagination{
		Limit:         10,
		Offset:        0,
		Search:        "",
		SortBy:        "created_at",
		SortDirection: "DESC",
	}

	ctx := context.Background()
	result, err := repo.ListProduct(ctx, sfp)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, result)
}

func TestListProductByStoreId(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewStoreRepo(dbx)

	storeId := "test_store_id"
	expectedData := []*Product{
		{
			ID:          "test_product_id",
			StoreID:     storeId,
			Name:        "test_product_name",
			Url:         "test_product_url",
			Price:       100,
			Description: "test_product_description",
			CreatedAt:   time.Now(),
			UpdatedAt: sql.NullTime{
				Time:  time.Time{},
				Valid: false,
			},
		},
	}
	listProductsQueryMock := "SELECT id, store_id, name, url, price, description, created_at, updated_at FROM products"
	mock.ExpectQuery(listProductsQueryMock).WillReturnRows(sqlmock.NewRows([]string{"id", "store_id", "name", "url", "price", "description", "created_at", "updated_at"}).
		AddRow(expectedData[0].ID, expectedData[0].StoreID, expectedData[0].Name, expectedData[0].Url, expectedData[0].Price, expectedData[0].Description, expectedData[0].CreatedAt, expectedData[0].UpdatedAt))

	sfp := &SearchFilterPagination{
		Limit:         10,
		Offset:        0,
		Search:        "",
		SortBy:        "created_at",
		SortDirection: "DESC",
	}

	ctx := context.Background()
	result, err := repo.ListProductByStoreId(ctx, sfp, storeId)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, result)
}

func TestCreateProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewStoreRepo(dbx)

	expectedData := &Product{
		ID:          "test_product_id",
		StoreID:     "test_store_id",
		Name:        "test_product_name",
		Url:         "test_product_url",
		Price:       100,
		Description: "test_product_description",
		CreatedAt:   time.Now(),
		UpdatedAt: sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		},
	}

	createProductQueryMock := "INSERT INTO products \\(store_id, name, url, price, description, created_at\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5, CURRENT_TIMESTAMP\\) RETURNING id"
	mock.ExpectQuery(createProductQueryMock).
		WithArgs(expectedData.StoreID, expectedData.Name, expectedData.Url, expectedData.Price, expectedData.Description).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedData.ID))

	ctx := context.Background()
	result, err := repo.CreateProduct(ctx, expectedData)
	assert.NoError(t, err)
	assert.Equal(t, expectedData.ID, result.ID)
}

func TestGetProductById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewStoreRepo(dbx)

	expectedData := &Product{
		ID:          "test_product_id",
		StoreID:     "test_store_id",
		Name:        "test_product_name",
		Url:         "test_product_url",
		Price:       100,
		Description: "test_product_description",
		CreatedAt:   time.Now(),
		UpdatedAt: sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		},
	}
	getProductByIdQueryMock := "SELECT id, store_id, name, url, price, description, created_at, updated_at FROM products WHERE id = \\$1 LIMIT 1"
	mock.ExpectQuery(getProductByIdQueryMock).WithArgs(expectedData.ID).WillReturnRows(sqlmock.NewRows([]string{"id", "store_id", "name", "url", "price", "description", "created_at", "updated_at"}).
		AddRow(expectedData.ID, expectedData.StoreID, expectedData.Name, expectedData.Url, expectedData.Price, expectedData.Description, expectedData.CreatedAt, expectedData.UpdatedAt))

	ctx := context.Background()
	result, err := repo.GetProductById(ctx, expectedData.ID)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, result)
}

func TestGetProductByUrl(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewStoreRepo(dbx)

	expectedData := &Product{
		ID:          "test_product_id",
		StoreID:     "test_store_id",
		Name:        "test_product_name",
		Url:         "test_product_url",
		Price:       100,
		Description: "test_product_description",
		CreatedAt:   time.Now(),
		UpdatedAt: sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		},
	}
	getProductByUrlQueryMock := "SELECT id, store_id, name, url, price, description, created_at, updated_at FROM products WHERE url = \\$1 LIMIT 1"
	mock.ExpectQuery(getProductByUrlQueryMock).WithArgs(expectedData.Url).WillReturnRows(sqlmock.NewRows([]string{"id", "store_id", "name", "url", "price", "description", "created_at", "updated_at"}).
		AddRow(expectedData.ID, expectedData.StoreID, expectedData.Name, expectedData.Url, expectedData.Price, expectedData.Description, expectedData.CreatedAt, expectedData.UpdatedAt))

	ctx := context.Background()
	result, err := repo.GetProductByUrl(ctx, expectedData.Url)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, result)
}

func TestUpdateProductById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewStoreRepo(dbx)

	expectedData := &Product{
		ID:          "test_product_id",
		StoreID:     "test_store_id",
		Name:        "test_product_name",
		Url:         "test_product_url",
		Price:       100,
		Description: "test_product_description",
		CreatedAt:   time.Now(),
		UpdatedAt: sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		},
	}

	updateProductByIdQueryMock := "UPDATE products SET store_id = \\$2, name = \\$3, url = \\$4, price = \\$5, description = \\$6, updated_at = CURRENT_TIMESTAMP WHERE id = \\$1"
	mock.ExpectExec(updateProductByIdQueryMock).
		WithArgs(expectedData.ID, expectedData.StoreID, expectedData.Name, expectedData.Url, expectedData.Price, expectedData.Description).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	err = repo.UpdateProduct(ctx, expectedData)
	assert.NoError(t, err)
}

func TestDeleteProductById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewStoreRepo(dbx)

	expectedData := &Product{
		ID:          "test_product_id",
		StoreID:     "test_store_id",
		Name:        "test_product_name",
		Url:         "test_product_url",
		Price:       100,
		Description: "test_product_description",
		CreatedAt:   time.Now(),
		UpdatedAt: sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		},
	}

	deleteProductByIdQueryMock := "DELETE FROM products WHERE id = \\$1"
	mock.ExpectExec(deleteProductByIdQueryMock).
		WithArgs(expectedData.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	err = repo.DeleteProduct(ctx, expectedData.ID)
	assert.NoError(t, err)
}
