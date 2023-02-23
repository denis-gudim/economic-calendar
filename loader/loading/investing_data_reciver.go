package loading

import (
	"context"
	"time"

	"github.com/denis-gudim/economic-calendar/loader/investing"
)

type InvestingDataReciver interface {
	GetEventsSchedule(ctx context.Context, dateFrom, dateTo time.Time) (map[int][]*investing.InvestingScheduleRow, error)
	GetEventsScheduleByLanguage(ctx context.Context, languageId int, dateFrom, dateTo time.Time) ([]*investing.InvestingScheduleRow, error)
	GetEventDetails(ctx context.Context, eventId int) ([]*investing.InvestingCalendarEvent, error)
	GetCountries(ctx context.Context) (map[int][]*investing.InvestingCountry, error)
}
