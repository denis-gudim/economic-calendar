package loading

import "economic-calendar/loader/data"

type EventsDataReciver interface {
	GetById(id int) (*data.Event, error)
	Save(e data.Event) error
}
