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

func TestCreateStore(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := repository.NewStoreRepo(dbx)
	mockStoreService := &StoreServiceMock{Mock: mocktest.Mock{}, repo: repo}
	svc := NewStoreMockService(mockStoreService, repo)

	ctx := context.Background()
	request := &domain.StoreRequest{
		Name:                 "test_store_name",
		Url:                  "test_store_url",
		Address:              "test_store_address",
		Phone:                "test_store_phone",
		OperationalTimeStart: 8,
		OperationalTimeEnd:   16,
	}

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

	createStoreQueryMock := "INSERT INTO stores \\(name, url, address, phone, operational_time_start, operational_time_end, created_at\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5, \\$6, CURRENT_TIMESTAMP\\) RETURNING id"
	mock.ExpectQuery(createStoreQueryMock).
		WithArgs(expectedStoreData.Name, expectedStoreData.Url, expectedStoreData.Address, expectedStoreData.Phone, expectedStoreData.OperationalTimeStart, expectedStoreData.OperationalTimeEnd).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedStoreData.ID))

	mockStoreService.Mock.On("MockCreateStore", request).Return(StoreRes(expectedStoreData), nil)
	store, err := svc.storeService.MockCreateStore(ctx, request)

	assert.NoError(t, err)
	expectedResponse := StoreRes(expectedStoreData)
	assert.Equal(t, expectedResponse.ID, store.ID)
	assert.Equal(t, expectedResponse.Name, store.Name)
	assert.Equal(t, expectedResponse.Url, store.Url)
	assert.Equal(t, expectedResponse.Address, store.Address)
	assert.Equal(t, expectedResponse.Phone, store.Phone)
	assert.Equal(t, expectedResponse.OperationalTimeStart, store.OperationalTimeStart)
	assert.Equal(t, expectedResponse.OperationalTimeEnd, store.OperationalTimeEnd)

}

func TestGetStoreById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := repository.NewStoreRepo(dbx)
	mockStoreService := &StoreServiceMock{Mock: mocktest.Mock{}, repo: repo}
	svc := NewStoreMockService(mockStoreService, repo)

	ctx := context.Background()

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

	mockStoreService.Mock.On("MockGetStoreById", expectedStoreData.ID).Return(StoreRes(expectedStoreData), nil)
	store, err := svc.storeService.MockGetStoreById(ctx, expectedStoreData.ID)

	assert.NoError(t, err)
	expectedResponse := StoreRes(expectedStoreData)
	assert.Equal(t, expectedResponse.ID, store.ID)
	assert.Equal(t, expectedResponse.Name, store.Name)
	assert.Equal(t, expectedResponse.Url, store.Url)
	assert.Equal(t, expectedResponse.Address, store.Address)
	assert.Equal(t, expectedResponse.Phone, store.Phone)
	assert.Equal(t, expectedResponse.OperationalTimeStart, store.OperationalTimeStart)
	assert.Equal(t, expectedResponse.OperationalTimeEnd, store.OperationalTimeEnd)
}

func TestUpdateStore(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := repository.NewStoreRepo(dbx)
	mockStoreService := &StoreServiceMock{Mock: mocktest.Mock{}, repo: repo}
	svc := NewStoreMockService(mockStoreService, repo)

	ctx := context.Background()
	request := &domain.StoreRequest{
		Name:                 "test_store_name",
		Url:                  "test_store_url",
		Address:              "test_store_address",
		Phone:                "test_store_phone",
		OperationalTimeStart: 8,
		OperationalTimeEnd:   16,
	}

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

	updateStoreQueryMock := "UPDATE stores SET name = \\$2, url = \\$3, address = \\$4, phone = \\$5, operational_time_start = \\$6, operational_time_end = \\$7, updated_at = CURRENT_TIMESTAMP WHERE id = \\$1"
	mock.ExpectExec(updateStoreQueryMock).
		WithArgs(expectedStoreData.ID, expectedStoreData.Name, expectedStoreData.Url, expectedStoreData.Address, expectedStoreData.Phone, expectedStoreData.OperationalTimeStart, expectedStoreData.OperationalTimeEnd).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mockStoreService.Mock.On("MockUpdateStore", request, expectedStoreData.ID).Return(nil)
	err = svc.storeService.MockUpdateStore(ctx, request, expectedStoreData.ID)

	assert.NoError(t, err)
}

func TestShowStoreProducts(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := repository.NewStoreRepo(dbx)
	mockStoreService := &StoreServiceMock{Mock: mocktest.Mock{}, repo: repo}
	svc := NewStoreMockService(mockStoreService, repo)

	ctx := context.Background()

	storeId := "test_store_id"
	expectedProductData := []*repository.Product{
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
		AddRow(expectedProductData[0].ID, expectedProductData[0].StoreID, expectedProductData[0].Name, expectedProductData[0].Url, expectedProductData[0].Price, expectedProductData[0].Description, expectedProductData[0].CreatedAt, expectedProductData[0].UpdatedAt))

	expectedCount := int64(8)
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

	mockStoreService.Mock.On("MockShowStoreProducts", &pagination, searchAndFilter, storeId).Return(expectedProductData, expectedCount, nil)
	products, count, err := svc.storeService.MockShowStoreProducts(ctx, &pagination, searchAndFilter, storeId)
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	assert.Equal(t, []*domain.Product{ProductRes(expectedProductData[0], expectedStoreData)}, products)
}
