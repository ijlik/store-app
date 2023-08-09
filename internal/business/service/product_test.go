package service

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ijlik/store-app/internal/adapter/repository"
	"github.com/ijlik/store-app/internal/business/domain"
	httppagination "github.com/ijlik/store-app/pkg/http/pagination"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	mocktest "github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestShowProducts(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := repository.NewStoreRepo(dbx)
	mockProductService := &ProductServiceMock{Mock: mocktest.Mock{}, repo: repo}
	svc := NewProductService(mockProductService, repo)

	expectedProductData := []*repository.Product{
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
		AddRow(expectedProductData[0].ID, expectedProductData[0].StoreID, expectedProductData[0].Name, expectedProductData[0].Url, expectedProductData[0].Price, expectedProductData[0].Description, expectedProductData[0].CreatedAt, expectedProductData[0].UpdatedAt))

	expectedCount := int64(10)
	countProductsQueryMock := "SELECT count\\(\\*\\) FROM products"
	mock.ExpectQuery(countProductsQueryMock).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

	expectedStoreData := &repository.Store{
		ID:                   "test_store_id",
		Name:                 "test_store_name",
		Url:                  "test_store_url",
		Address:              "test_store_address",
		Phone:                "test_store_phone",
		OperationalTimeStart: 8,
		OperationalTimeEnd:   16,
		CreatedAt:            time.Now(),
		UpdatedAt: sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		},
	}

	getStoreByIdQueryMock := "SELECT id, name, url, address, phone, operational_time_start, operational_time_end, created_at, updated_at FROM stores WHERE id = \\$1 LIMIT 1"
	rows := sqlmock.NewRows([]string{"id", "name", "url", "address", "phone", "operational_time_start", "operational_time_end", "created_at", "updated_at"}).
		AddRow(expectedStoreData.ID, expectedStoreData.Name, expectedStoreData.Url, expectedStoreData.Address, expectedStoreData.Phone, expectedStoreData.OperationalTimeStart, expectedStoreData.OperationalTimeEnd, expectedStoreData.CreatedAt, nil)
	mock.ExpectQuery(getStoreByIdQueryMock).WithArgs(expectedStoreData.ID).WillReturnRows(rows)

	pagination := httppagination.Pagination{
		Limit:  10,
		Offset: 0,
	}

	searchAndFilter := &domain.SearchAndFilterProduct{
		Limit:         10,
		Page:          0,
		Search:        "",
		SortBy:        "created_at",
		SortDirection: "DESC",
	}

	mockProductService.Mock.On("MockShowProducts", &pagination, searchAndFilter).Return(expectedProductData, expectedCount, nil)
	products, count, err := svc.productService.MockShowProducts(context.Background(), &pagination, searchAndFilter)
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	assert.Equal(t, []*domain.Product{ProductRes(expectedProductData[0], expectedStoreData)}, products)
}

func TestCreateProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := repository.NewStoreRepo(dbx)
	mockProductService := &ProductServiceMock{Mock: mocktest.Mock{}, repo: repo}
	svc := NewProductService(mockProductService, repo)

	ctx := context.Background()
	request := &domain.ProductRequest{
		StoreID:     "test_store_id",
		Name:        "test_product_name",
		Url:         "test_product_url",
		Price:       100,
		Description: "test_product_description",
	}

	expectedProductData := &repository.Product{
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
		WithArgs(expectedProductData.StoreID, expectedProductData.Name, expectedProductData.Url, expectedProductData.Price, expectedProductData.Description).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedProductData.ID))

	expectedStoreData := &repository.Store{
		ID:                   "test_store_id",
		Name:                 "test_store_name",
		Url:                  "test_store_url",
		Address:              "test_store_address",
		Phone:                "test_store_phone",
		OperationalTimeStart: 8,
		OperationalTimeEnd:   16,
		CreatedAt:            time.Now(),
		UpdatedAt: sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		},
	}

	getStoreByIdQueryMock := "SELECT id, name, url, address, phone, operational_time_start, operational_time_end, created_at, updated_at FROM stores WHERE id = \\$1 LIMIT 1"
	rows := sqlmock.NewRows([]string{"id", "name", "url", "address", "phone", "operational_time_start", "operational_time_end", "created_at", "updated_at"}).
		AddRow(expectedStoreData.ID, expectedStoreData.Name, expectedStoreData.Url, expectedStoreData.Address, expectedStoreData.Phone, expectedStoreData.OperationalTimeStart, expectedStoreData.OperationalTimeEnd, expectedStoreData.CreatedAt, nil)
	mock.ExpectQuery(getStoreByIdQueryMock).WithArgs(expectedStoreData.ID).WillReturnRows(rows)

	mockProductService.Mock.On("MockCreateProduct", request).Return(ProductRes(expectedProductData, expectedStoreData), nil)
	product, err := svc.productService.MockCreateProduct(ctx, request)

	assert.NoError(t, err)
	expectedResponse := ProductRes(expectedProductData, expectedStoreData)
	assert.Equal(t, expectedResponse.ID, product.ID)
	assert.Equal(t, expectedResponse.Name, product.Name)
	assert.Equal(t, expectedResponse.Url, product.Url)
	assert.Equal(t, expectedResponse.Price, product.Price)
	assert.Equal(t, expectedResponse.Description, product.Description)
	assert.Equal(t, expectedResponse.Store.ID, product.Store.ID)
	assert.Equal(t, expectedResponse.Store.Name, product.Store.Name)
	assert.Equal(t, expectedResponse.Store.Url, product.Store.Url)
	assert.Equal(t, expectedResponse.Store.Address, product.Store.Address)
	assert.Equal(t, expectedResponse.Store.Phone, product.Store.Phone)
	assert.Equal(t, expectedResponse.Store.OperationalTimeStart, product.Store.OperationalTimeStart)
	assert.Equal(t, expectedResponse.Store.OperationalTimeEnd, product.Store.OperationalTimeEnd)

}

func TestGetProductByUrl(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := repository.NewStoreRepo(dbx)
	mockProductService := &ProductServiceMock{Mock: mocktest.Mock{}, repo: repo}
	svc := NewProductService(mockProductService, repo)

	ctx := context.Background()
	url := "test_product_url"
	expectedProductData := &repository.Product{
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
	mock.ExpectQuery(getProductByUrlQueryMock).WithArgs(expectedProductData.Url).WillReturnRows(sqlmock.NewRows([]string{"id", "store_id", "name", "url", "price", "description", "created_at", "updated_at"}).
		AddRow(expectedProductData.ID, expectedProductData.StoreID, expectedProductData.Name, expectedProductData.Url, expectedProductData.Price, expectedProductData.Description, expectedProductData.CreatedAt, expectedProductData.UpdatedAt))

	expectedStoreData := &repository.Store{
		ID:                   "test_store_id",
		Name:                 "test_store_name",
		Url:                  "test_store_url",
		Address:              "test_store_address",
		Phone:                "test_store_phone",
		OperationalTimeStart: 8,
		OperationalTimeEnd:   16,
		CreatedAt:            time.Now(),
		UpdatedAt: sql.NullTime{
			Time:  time.Time{},
			Valid: false,
		},
	}

	getStoreByIdQueryMock := "SELECT id, name, url, address, phone, operational_time_start, operational_time_end, created_at, updated_at FROM stores WHERE id = \\$1 LIMIT 1"
	rows := sqlmock.NewRows([]string{"id", "name", "url", "address", "phone", "operational_time_start", "operational_time_end", "created_at", "updated_at"}).
		AddRow(expectedStoreData.ID, expectedStoreData.Name, expectedStoreData.Url, expectedStoreData.Address, expectedStoreData.Phone, expectedStoreData.OperationalTimeStart, expectedStoreData.OperationalTimeEnd, expectedStoreData.CreatedAt, nil)
	mock.ExpectQuery(getStoreByIdQueryMock).WithArgs(expectedStoreData.ID).WillReturnRows(rows)

	mockProductService.Mock.On("MockGetProductByUrl", url).Return(ProductRes(expectedProductData, expectedStoreData), nil)
	product, err := svc.productService.MockGetProductByUrl(ctx, url)

	assert.NoError(t, err)
	expectedResponse := ProductRes(expectedProductData, expectedStoreData)
	assert.Equal(t, expectedResponse.ID, product.ID)
	assert.Equal(t, expectedResponse.Name, product.Name)
	assert.Equal(t, expectedResponse.Url, product.Url)
	assert.Equal(t, expectedResponse.Price, product.Price)
	assert.Equal(t, expectedResponse.Description, product.Description)
	assert.Equal(t, expectedResponse.Store.ID, product.Store.ID)
	assert.Equal(t, expectedResponse.Store.Name, product.Store.Name)
	assert.Equal(t, expectedResponse.Store.Url, product.Store.Url)
	assert.Equal(t, expectedResponse.Store.Address, product.Store.Address)
	assert.Equal(t, expectedResponse.Store.Phone, product.Store.Phone)
	assert.Equal(t, expectedResponse.Store.OperationalTimeStart, product.Store.OperationalTimeStart)
	assert.Equal(t, expectedResponse.Store.OperationalTimeEnd, product.Store.OperationalTimeEnd)
}

func TestUpdateProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := repository.NewStoreRepo(dbx)
	mockProductService := &ProductServiceMock{Mock: mocktest.Mock{}, repo: repo}
	svc := NewProductService(mockProductService, repo)

	ctx := context.Background()
	productId := "test_product_id"
	request := &domain.ProductRequest{
		StoreID:     "test_store_id",
		Name:        "test_product_name",
		Url:         "test_product_url",
		Price:       100,
		Description: "test_product_description",
	}

	expectedProductData := &repository.Product{
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
		WithArgs(expectedProductData.ID, expectedProductData.StoreID, expectedProductData.Name, expectedProductData.Url, expectedProductData.Price, expectedProductData.Description).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mockProductService.Mock.On("MockUpdateProduct", request, productId).Return(nil)
	err = svc.productService.MockUpdateProduct(ctx, request, productId)
	assert.NoError(t, err)
}

func TestDeleteProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := repository.NewStoreRepo(dbx)
	mockProductService := &ProductServiceMock{Mock: mocktest.Mock{}, repo: repo}
	svc := NewProductService(mockProductService, repo)

	ctx := context.Background()
	productId := "test_product_id"

	deleteProductByIdQueryMock := "DELETE FROM products WHERE id = \\$1"
	mock.ExpectExec(deleteProductByIdQueryMock).
		WithArgs(productId).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mockProductService.Mock.On("MockDeleteProduct", productId).Return(nil)
	err = svc.productService.MockDeleteProduct(ctx, productId)
	assert.NoError(t, err)
}
