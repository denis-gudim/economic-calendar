package loading

import (
	"economic-calendar/loader/data"
	"time"
)

type EventScheduleDataReciver interface {
	GetFirst() (*data.EventSchedule, error)
	GetByDates(from, to time.Time) ([]data.EventSchedule, error)
	Save(es data.EventSchedule) error
}
