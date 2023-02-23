package loading

import (
	"context"

	"github.com/denis-gudim/economic-calendar/loader/data"
)

type CountriesDataReciver interface {
	GetAll(ctx context.Context) ([]data.Country, error)
	Save(ctx context.Context, c data.Country) error
}
