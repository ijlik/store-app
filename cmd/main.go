package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	configenv "github.com/ijlik/store-app/pkg/config"
	configdata "github.com/ijlik/store-app/pkg/config/data"
	httpmiddlewaresdk "github.com/ijlik/store-app/pkg/http/middleware"
	"github.com/jmoiron/sqlx"

	// internal package
	"github.com/ijlik/store-app/internal/adapter/repository"
	"github.com/ijlik/store-app/internal/business/port"
	"github.com/ijlik/store-app/internal/business/service"
	httpdelivery "github.com/ijlik/store-app/internal/handler/http"
	_ "github.com/lib/pq"
)

var config configdata.Config

func getService(
	db *sqlx.DB,
) port.StoreDomainService {
	repo := repository.NewStoreRepo(db)
	services := service.NewStoreService(
		repo,
		config,
	)

	return services
}

func getConfig() configdata.Config {
	c := configenv.NewConfig("", 5)

	if c == nil {
		panic(errors.New("missing config"))
	}

	return c
}

func getDatabase() (*sqlx.DB, error) {
	var (
		host     = config.GetString("DB_HOST")
		dbPort   = config.GetInt("DB_PORT")
		user     = config.GetString("DB_USER")
		password = config.GetString("DB_PASSWORD")
		dbname   = config.GetString("DB_NAME")
		timeZone = "UTC"
	)

	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable TimeZone=%s",
		host, dbPort, user, password, dbname, timeZone)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return sqlx.NewDb(db, "postgres"), nil
}

func main() {
	// get config
	config = getConfig()

	db, err := getDatabase()
	if err != nil {
		panic(err)
	}

	defer db.Close()

	router := gin.Default()
	router.Use(
		httpmiddlewaresdk.WithAllowedCORS(),
	)

	services := getService(db)

	httpdelivery.HandlerHttp(
		router,
		config,
		services,
	)
}
