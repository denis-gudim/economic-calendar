package loading

import (
	"economic-calendar/loader/data"
	"time"
)

type EventScheduleDataReciver interface {
	GetFirst(done bool) (*data.EventSchedule, error)
	GetByDates(from, to time.Time) ([]data.EventSchedule, error)
	Save(es data.EventSchedule) error
}
