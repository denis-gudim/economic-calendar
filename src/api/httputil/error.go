package httputil

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BadRequestError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"400 Bad Request: error details text"`
}

func NewBadRequestError(ctx *gin.Context, err error) {
	er := BadRequestError{
		Code:    http.StatusBadRequest,
		Message: fmt.Errorf("400 Bad Request: %w", err).Error(),
	}
	ctx.JSON(er.Code, er)
}

type NotFoundError struct {
	Code    int    `json:"code" example:"404"`
	Message string `json:"message" example:"404 Not Found: error details text"`
}

func NewNotFoundError(ctx *gin.Context, err error) {
	er := InternalServerError{
		Code:    http.StatusNotFound,
		Message: fmt.Errorf("404 Not Found: %w", err).Error(),
	}
	ctx.JSON(er.Code, er)
}

type InternalServerError struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"500 Internal Server Error"`
}

func NewInternalServerError(ctx *gin.Context, err error) {
	er := InternalServerError{
		Code:    http.StatusInternalServerError,
		Message: "500 Internal Server Error",
	}
	ctx.JSON(er.Code, er)
}
