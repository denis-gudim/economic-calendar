package loading

import (
	"context"
	"time"

	"github.com/denis-gudim/economic-calendar/loader/data"
)

type EventScheduleDataReciver interface {
	GetFirst(ctx context.Context, done bool) (*data.EventSchedule, error)
	GetByDates(ctx context.Context, from, to time.Time) ([]data.EventSchedule, error)
	Save(ctx context.Context, es data.EventSchedule) error
}
