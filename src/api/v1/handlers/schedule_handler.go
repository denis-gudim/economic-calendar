package v1

import (
	"net/http"
	"strconv"
	"time"

	"github.com/denis-gudim/economic-calendar/api/app"
	"github.com/denis-gudim/economic-calendar/api/v1/data"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ScheduleHandler struct {
	repository *data.ScheduleRepository
	logger     *zap.Logger
	baseHandler
}

func InitScheduleHandler(rg *gin.RouterGroup, cnf app.Config, logger *zap.Logger) {

	handler := ScheduleHandler{
		repository: data.NewScheduleRepository(cnf),
		logger:     logger,
	}

	rg.GET("schedule", handler.Get)
	rg.GET("schedule/:id", handler.GetEventById)
}

func (h *ScheduleHandler) Get(c *gin.Context) {

	lang := c.DefaultQuery("lang", "en")
	from := c.Query("from")
	to := c.Query("to")

	fromDate, err := time.Parse("2006-01-02", from)

	if err != nil {
		h.writeBadRequest(c, "invalid from date value '%s'", from)
		return
	}

	toDate, err := time.ParseInLocation("2006-01-02", to, time.UTC)

	if err != nil {
		h.writeBadRequest(c, "invalid to date value '%s'", to)
		return
	}

	rows, err := h.repository.GetScheduleByDates(c, fromDate, toDate, lang)

	if err != nil {
		h.logger.Error(err.Error(),
			zap.Time("from", fromDate),
			zap.Time("to", toDate),
			zap.String("lang", lang),
		)
		h.writeServerError(c)
		return
	}

	c.JSON(http.StatusOK, rows)
}

func (h *ScheduleHandler) GetEventById(c *gin.Context) {

	lang := c.DefaultQuery("lang", "en")
	id := c.Param("id")

	eventId, err := strconv.Atoi(id)

	if err != nil {
		h.writeBadRequest(c, "invalid event id value '%s'", id)
		return
	}

	rows, err := h.repository.GetEventById(c, eventId, lang)

	if err != nil {
		h.logger.Error(err.Error(),
			zap.Int("eventId", eventId),
			zap.String("lang", lang),
		)
		h.writeServerError(c)
		return
	}

	c.JSON(http.StatusOK, rows)
}
