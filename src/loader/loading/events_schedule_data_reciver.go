package loading

import (
	"time"

	"github.com/denis-gudim/economic-calendar/loader/data"
)

type EventScheduleDataReciver interface {
	GetFirst(done bool) (*data.EventSchedule, error)
	GetByDates(from, to time.Time) ([]data.EventSchedule, error)
	Save(es data.EventSchedule) error
}
