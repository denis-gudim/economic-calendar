package loading

import (
	"time"

	"github.com/denis-gudim/economic-calendar/loader/investing"
)

type InvestingDataReciver interface {
	GetEventsSchedule(dateFrom, dateTo time.Time) (map[int][]*investing.InvestingScheduleRow, error)
	GetEventsScheduleByLanguage(languageId int, dateFrom, dateTo time.Time) ([]*investing.InvestingScheduleRow, error)
	GetEventDetails(eventId int) ([]*investing.InvestingCalendarEvent, error)
	GetCountries() (map[int][]*investing.InvestingCountry, error)
}
