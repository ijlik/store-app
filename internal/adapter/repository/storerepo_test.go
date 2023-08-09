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

func TestCreateStore(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewStoreRepo(dbx)

	expectedData := &Store{
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
		WithArgs(expectedData.Name, expectedData.Url, expectedData.Address, expectedData.Phone, expectedData.OperationalTimeStart, expectedData.OperationalTimeEnd).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedData.ID))

	ctx := context.Background()
	result, err := repo.CreateStore(ctx, expectedData)
	assert.NoError(t, err)
	assert.Equal(t, expectedData.ID, result.ID)
}

func TestGetStoreById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewStoreRepo(dbx)

	expectedData := &Store{
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
		AddRow(expectedData.ID, expectedData.Name, expectedData.Url, expectedData.Address, expectedData.Phone, expectedData.OperationalTimeStart, expectedData.OperationalTimeEnd, expectedData.CreatedAt, nil)
	mock.ExpectQuery(getStoreByIdQueryMock).WithArgs(expectedData.ID).WillReturnRows(rows)

	ctx := context.Background()
	result, err := repo.GetStoreById(ctx, expectedData.ID)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, result)
}

func TestUpdateStore(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewStoreRepo(dbx)

	expectedData := &Store{
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
		WithArgs(expectedData.ID, expectedData.Name, expectedData.Url, expectedData.Address, expectedData.Phone, expectedData.OperationalTimeStart, expectedData.OperationalTimeEnd).
		WillReturnResult(sqlmock.NewResult(1, 1))

	ctx := context.Background()
	err = repo.UpdateStore(ctx, expectedData)
	assert.NoError(t, err)
}
