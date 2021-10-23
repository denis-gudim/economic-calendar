package controllers

import (
	"context"
	"net/http"

	"github.com/denis-gudim/economic-calendar/api/httputil"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	data "github.com/denis-gudim/economic-calendar/api/v1/data"
)

type CountriesDataReciver interface {
	GetCountriesByLanguage(ctx context.Context, langCode string) ([]data.Country, error)
}

type CountriesController struct {
	repository CountriesDataReciver
	logger     *zap.Logger
}

func NewCountriesController(r CountriesDataReciver, l *zap.Logger) *CountriesController {
	return &CountriesController{
		repository: r,
		logger:     l,
	}
}

// CountriesGet godoc
// @Summary Countries list by language code
// @Schemes http|https
// @Description Returns list of countries translated to specified language.
// @Tags Countries
// @Accept json
// @Produce json
// @Param lang query string false "language code value" default(en)
// @Success 200 {array} data.Country
// @Failure 500 {object} httputil.InternalServerError
// @Router /countries [get]
func (h *CountriesController) GetByLanguage(ctx *gin.Context) {

	lang := ctx.DefaultQuery("lang", "en")

	countries, err := h.repository.GetCountriesByLanguage(ctx, lang)

	if err != nil {
		h.logger.Error(err.Error(), zap.String("lang", lang))
		httputil.NewInternalServerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, countries)
}
