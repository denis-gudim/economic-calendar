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

type EventsHandler struct {
	repository *data.EventsRepository
	logger     *zap.Logger
	baseHandler
}

func NewEventsHandler(cnf app.Config, logger *zap.Logger) EventsHandler {
	return EventsHandler{
		repository: data.NewEventsRepository(cnf),
		logger:     logger,
	}
}

func (h *EventsHandler) GetEventsSchdule(c *gin.Context) {

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

func (h *EventsHandler) GetEventDetails(c *gin.Context) {

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

func (h *EventsHandler) GetEventHistory(c *gin.Context) {

	id := c.Param("id")

	eventId, err := strconv.Atoi(id)

	if err != nil {
		h.writeBadRequest(c, "invalid event id value '%s'", id)
		return
	}

	rows, err := h.repository.GetHistoryById(c, eventId)

	if err != nil {
		h.logger.Error(err.Error(),
			zap.Int("eventId", eventId),
		)
		h.writeServerError(c)
		return
	}

	c.JSON(http.StatusOK, rows)
}
