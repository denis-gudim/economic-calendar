package v1

import (
	"net/http"

	"github.com/denis-gudim/economic-calendar/api/app"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	data "github.com/denis-gudim/economic-calendar/api/v1/data"
)

type CountriesHandler struct {
	repository *data.CountriesRepository
	logger     *zap.Logger
	baseHandler
}

func InitCountriesHandler(rg *gin.RouterGroup, cnf app.Config, logger *zap.Logger) {

	handler := CountriesHandler{
		repository: data.NewCountriesRepository(cnf),
		logger:     logger,
	}

	rg.GET("countries", handler.Get)
}

func (h *CountriesHandler) Get(c *gin.Context) {

	lang := c.DefaultQuery("lang", "en")

	countries, err := h.repository.GetCountriesByLanguage(c, lang)

	if err != nil {
		h.logger.Error(err.Error(), zap.String("lang", lang))
		h.writeServerError(c)
		return
	}

	c.JSON(http.StatusOK, countries)
}
