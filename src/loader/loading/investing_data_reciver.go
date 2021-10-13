package loading

import (
	"economic-calendar/loader/investing"
	"time"
)

type InvestingDataReciver interface {
	GetEventsSchedule(dateFrom, dateTo time.Time) (map[int][]*investing.InvestingScheduleRow, error)
	GetEventsScheduleByLanguage(languageId int, dateFrom, dateTo time.Time) ([]*investing.InvestingScheduleRow, error)
	GetEventDetails(eventId int) ([]*investing.InvestingCalendarEvent, error)
	GetCountries() (map[int][]*investing.InvestingCountry, error)
}
