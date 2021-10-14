package loading

import "github.com/denis-gudim/economic-calendar/loader/data"

type CountriesDataReciver interface {
	GetAll() ([]data.Country, error)
	Save(c data.Country) error
}
