package http

import (
	"github.com/gin-gonic/gin"
	"github.com/ijlik/store-app/internal/business/domain"
	errpkg "github.com/ijlik/store-app/pkg/error"
	httppkg "github.com/ijlik/store-app/pkg/http"
	httppagination "github.com/ijlik/store-app/pkg/http/pagination"
)

func (rh *requestHandler) ListProducts(c *gin.Context) {
	var (
		sf         = domain.SearchAndFilterProduct{}
		pagination *httppagination.Pagination
	)

	if errQuery := c.ShouldBindQuery(&sf); errQuery != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, "")
		return
	}

	err := sf.Validate()
	if err != nil {
		httppkg.BuildErrorResponse(c, err.GetCode(), err.Error())
		return
	}
	pagination = httppagination.NewPaginate(sf.Limit, sf.Page)

	err = rh.service.ShowProducts(c.Request.Context(), pagination, &sf)
	if err != nil {
		httppkg.BuildErrorResponse(c, err.GetCode(), err.Error())
		return
	}

	pagination.BuildPaginationResponse(c)
}

func (rh *requestHandler) CreateProduct(c *gin.Context) {
	ctx := c.Request.Context()
	var request domain.ProductRequest

	err := decodeRequest(c, &request)
	if err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}
	if err := request.Validate(); err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}

	product, err := rh.service.CreateProduct(ctx, &request)
	if err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}

	response := httppkg.DefaultSuccessResponse(product)
	c.JSON(response.HttpCode, response)
}

func (rh *requestHandler) ShowProduct(c *gin.Context) {
	ctx := c.Request.Context()
	var params = domain.HttpProductUrlParams{}

	if errJson := c.ShouldBindUri(&params); errJson != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, "")
		return
	}

	product, err := rh.service.GetProductByUrl(ctx, params.Url)
	if err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}

	response := httppkg.DefaultSuccessResponse(product)
	c.JSON(response.HttpCode, response)
}

func (rh *requestHandler) UpdateProduct(c *gin.Context) {
	ctx := c.Request.Context()
	var params = domain.HttpProductIdParams{}

	if errJson := c.ShouldBindUri(&params); errJson != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, "")
		return
	}

	var request domain.ProductRequest

	err := decodeRequest(c, &request)
	if err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}
	if err := request.Validate(); err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}

	err = rh.service.UpdateProduct(ctx, &request, params.ID)
	if err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}

	response := httppkg.DefaultSuccessResponse(nil)
	c.JSON(response.HttpCode, response)
}

func (rh *requestHandler) DeleteProduct(c *gin.Context) {
	ctx := c.Request.Context()
	var params = domain.HttpProductIdParams{}

	if errJson := c.ShouldBindUri(&params); errJson != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, "")
		return
	}

	err := rh.service.DeleteProduct(ctx, params.ID)
	if err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}

	response := httppkg.DefaultSuccessResponse(nil)
	c.JSON(response.HttpCode, response)
}
