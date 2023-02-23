package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/denis-gudim/economic-calendar/api/httputil"
	"github.com/denis-gudim/economic-calendar/api/v1/data"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
)

type EventsDataReciver interface {
	GetScheduleByDates(ctx context.Context, from, to time.Time, langCode string) ([]data.Event, error)
	GetEventById(ctx context.Context, eventId int, langCode string) (*data.EventDetails, error)
	GetHistoryById(ctx context.Context, eventId int) ([]data.EventRow, error)
}

type EventsController struct {
	repository EventsDataReciver
	logger     *zap.Logger
}

func NewEventsController(r EventsDataReciver, l *zap.Logger) *EventsController {
	return &EventsController{
		repository: r,
		logger:     l,
	}
}

// GetEventsSchedule godoc
// @Summary Event schedule between dates
// @Schemes http|https
// @Description Returns event schedule list in dates diapasone
// @Tags Events
// @Accept json
// @Produce json
// @Param from query string true "from date string in ISO 8601 format e.g. 2021-10-10"
// @Param to query string true "to date string in ISO 8601 format e.g. 2021-10-10"
// @Param lang query string false "language code value" default(en)
// @Success 200 {array} data.Event
// @Failure 400 {object} httputil.BadRequestError
// @Failure 500 {object} httputil.InternalServerError
// @Router /events [get]
func (h *EventsController) GetEventsSchedule(ctx *gin.Context) {

	lang := ctx.DefaultQuery("lang", "en")
	from := ctx.Query("from")
	to := ctx.Query("to")

	fromDate, err := time.Parse("2006-01-02", from)

	if err != nil {
		err = xerrors.Errorf("invalid from date value '%s': %w", from, err)
		httputil.NewBadRequestError(ctx, err)
		return
	}

	toDate, err := time.ParseInLocation("2006-01-02", to, time.UTC)

	if err != nil {
		err = xerrors.Errorf("invalid to date value '%s': %w", to, err)
		httputil.NewBadRequestError(ctx, err)
		return
	}

	rows, err := h.repository.GetScheduleByDates(ctx, fromDate, toDate, lang)

	if err != nil {
		h.logger.Error(err.Error(),
			zap.Time("from", fromDate),
			zap.Time("to", toDate),
			zap.String("lang", lang),
		)
		httputil.NewInternalServerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, rows)
}

// GetEventDetails godoc
// @Summary Event details by id
// @Schemes http|https
// @Description Returns event details with last schedule information by specified identifier
// @Tags Events
// @Accept json
// @Produce json
// @Param eventId path int true "event identifier" example(368)
// @Param lang query string false "language code value" default(en)
// @Success 200 {object} data.EventDetails
// @Failure 400 {object} httputil.BadRequestError
// @Failure 404 {object} httputil.NotFoundError
// @Failure 500 {object} httputil.InternalServerError
// @Router /events/{eventId} [get]
func (h *EventsController) GetEventDetails(ctx *gin.Context) {

	lang := ctx.DefaultQuery("lang", "en")
	id := ctx.Param("eventId")

	eventId, err := strconv.Atoi(id)

	if err != nil {
		err = xerrors.Errorf("invalid event id value '%s': %w", id, err)
		httputil.NewBadRequestError(ctx, err)
		return
	}

	event, err := h.repository.GetEventById(ctx, eventId, lang)

	if err != nil {
		h.logger.Error(err.Error(),
			zap.Int("eventId", eventId),
			zap.String("lang", lang),
		)

		httputil.NewInternalServerError(ctx, err)

		return
	}

	if event == nil {
		httputil.NewNotFoundError(ctx, fmt.Errorf("event with id %d not found", eventId))
	}

	ctx.JSON(http.StatusOK, event)
}

// GetEventHistory godoc
// @Summary Event history by id
// @Schemes http|https
// @Description Returns event history list by event id
// @Tags Events
// @Accept json
// @Produce json
// @Param eventId path int true "event identifier" example(368)
// @Success 200 {array} data.EventRow
// @Failure 400 {object} httputil.BadRequestError
// @Failure 500 {object} httputil.InternalServerError
// @Router /events/{eventId}/history [get]
func (h *EventsController) GetEventHistory(ctx *gin.Context) {

	id := ctx.Param("eventId")

	eventId, err := strconv.Atoi(id)

	if err != nil {
		err = xerrors.Errorf("invalid event id value '%s': %w", id, err)
		httputil.NewBadRequestError(ctx, err)
		return
	}

	rows, err := h.repository.GetHistoryById(ctx, eventId)

	if err != nil {
		h.logger.Error(err.Error(),
			zap.Int("eventId", eventId),
		)
		httputil.NewInternalServerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, rows)
}
