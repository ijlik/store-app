package http

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	configdata "github.com/ijlik/store-app/pkg/config/data"
	httppkg "github.com/ijlik/store-app/pkg/http"

	"github.com/ijlik/store-app/internal/business/port"
)

type requestHandler struct {
	config  configdata.Config
	service port.StoreDomainService
}

func HandlerHttp(
	router *gin.Engine,
	config configdata.Config,
	service port.StoreDomainService,
) {
	rh := requestHandler{
		config:  config,
		service: service,
	}

	routeHandler(router, rh)

	addr := config.GetString("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	fmt.Println("HTTP Running ON ", addr)

	httppkg.Serve(router, addr)
}

func routeHandler(router *gin.Engine, rh requestHandler) {
	storeRoute := router.Group("/store")
	storeRoute.POST("", rh.CreateStore)
	storeRoute.GET("/:id", rh.ShowStore)
	storeRoute.PUT("/:id", rh.UpdateStore)
	storeRoute.GET("/:id/products", rh.ShowStoreProducts)

	productRoute := router.Group("/product")
	productRoute.GET("", rh.ListProducts)
	productRoute.POST("", rh.CreateProduct)
	productRoute.GET("/:url", rh.ShowProduct)
	productRoute.PUT("/:id", rh.UpdateProduct)
	productRoute.DELETE("/:id", rh.DeleteProduct)

}

func decodeRequest(c *gin.Context, i interface{}) error {
	err := json.NewDecoder(c.Request.Body).Decode(i)
	if err != nil {
		return err
	}

	return nil
}
