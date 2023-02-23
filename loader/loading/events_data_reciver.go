package loading

import (
	"context"

	"github.com/denis-gudim/economic-calendar/loader/data"
)

type EventsDataReciver interface {
	GetById(ctx context.Context, id int) (*data.Event, error)
	Save(ctx context.Context, e data.Event) error
}
