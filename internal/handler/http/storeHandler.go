package http

import (
	"github.com/gin-gonic/gin"
	"github.com/ijlik/store-app/internal/business/domain"
	errpkg "github.com/ijlik/store-app/pkg/error"
	httppkg "github.com/ijlik/store-app/pkg/http"
	httppagination "github.com/ijlik/store-app/pkg/http/pagination"
)

func (rh *requestHandler) CreateStore(c *gin.Context) {
	ctx := c.Request.Context()
	var request domain.StoreRequest

	err := decodeRequest(c, &request)
	if err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}
	if err := request.Validate(); err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}

	store, err := rh.service.CreateStore(ctx, &request)
	if err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}

	response := httppkg.DefaultSuccessResponse(store)
	c.JSON(response.HttpCode, response)
}

func (rh *requestHandler) ShowStore(c *gin.Context) {
	ctx := c.Request.Context()
	var params = domain.HttpStoreIdParams{}

	if errJson := c.ShouldBindUri(&params); errJson != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, "")
		return
	}

	store, err := rh.service.GetStoreById(ctx, params.ID)
	if err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}

	response := httppkg.DefaultSuccessResponse(store)
	c.JSON(response.HttpCode, response)
}

func (rh *requestHandler) UpdateStore(c *gin.Context) {
	ctx := c.Request.Context()
	var params = domain.HttpStoreIdParams{}

	if errJson := c.ShouldBindUri(&params); errJson != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, "")
		return
	}

	var request domain.StoreRequest

	err := decodeRequest(c, &request)
	if err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}
	if err := request.Validate(); err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}

	err = rh.service.UpdateStore(ctx, &request, params.ID)
	if err != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, err.Error())
		return
	}

	response := httppkg.DefaultSuccessResponse(nil)
	c.JSON(response.HttpCode, response)
}

func (rh *requestHandler) ShowStoreProducts(c *gin.Context) {
	var (
		sf         = domain.SearchAndFilterProduct{}
		pagination *httppagination.Pagination
	)

	var params = domain.HttpStoreIdParams{}

	if errJson := c.ShouldBindUri(&params); errJson != nil {
		httppkg.BuildErrorResponse(c, errpkg.ErrBadRequest, "")
		return
	}

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

	err = rh.service.ShowStoreProducts(c.Request.Context(), pagination, &sf, params.ID)
	if err != nil {
		httppkg.BuildErrorResponse(c, err.GetCode(), err.Error())
		return
	}

	pagination.BuildPaginationResponse(c)
}
