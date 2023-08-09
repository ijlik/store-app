package http

import (
	"github.com/gin-gonic/gin"
	errpkg "github.com/ijlik/store-app/pkg/error"
)

type DefaultResponse struct {
	Code     string      `json:"code"`
	Message  string      `json:"message"`
	Data     interface{} `json:"data,omitempty"`
	HttpCode int         `json:"-"`
}

func DefaultSuccessResponse(data interface{}) DefaultResponse {
	return DefaultResponse{
		Code:     errpkg.GetCode(0),
		Message:  errpkg.GetMessage(0),
		Data:     data,
		HttpCode: errpkg.GetHttpStatus(0),
	}
}

func DefaultResponseErrorWithMessage(code errpkg.ErrCode, msg string) DefaultResponse {
	if msg == "" {
		msg = errpkg.GetMessage(code)
	}

	return DefaultResponse{
		Code:     errpkg.GetCode(code),
		Message:  msg,
		HttpCode: errpkg.GetHttpStatus(code),
	}
}

func BuildErrorResponse(c *gin.Context, code errpkg.ErrCode, msg string) {
	errResponse := DefaultResponseErrorWithMessage(code, msg)
	c.Header("Content-Type", "application/json")
	c.JSON(errResponse.HttpCode, errResponse)
	c.Abort()
}

func BuildSuccessResponse(data interface{}, c *gin.Context) {
	resp := DefaultSuccessResponse(data)
	c.Header("Content-Type", "application/json")
	c.JSON(resp.HttpCode, resp)
}
