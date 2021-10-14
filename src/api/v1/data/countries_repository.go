package data

import (
	"context"

	"github.com/denis-gudim/economic-calendar/api/app"
	"golang.org/x/xerrors"
)

type CountriesRepository struct {
	baseRepository
}

func NewCountriesRepository(cnf app.Config) *CountriesRepository {
	r := CountriesRepository{}
	r.config = cnf

	return &r
}

func (r *CountriesRepository) GetCountriesByLanguage(ctx context.Context, langCode string) ([]Country, error) {
	const errFormat = "get countries by language '%s' failed: %s: %w"

	db, err := r.connectDB(ctx)

	if err != nil {
		return nil, xerrors.Errorf(errFormat, langCode, "connect to db", err)
	}

	defer db.Close()

	countries := make([]Country, 0, 128)

	err = db.SelectContext(ctx, &countries,
		`SELECT c.id, c.code, c.continent_code, c.currency, ct.title AS name
		 FROM countries AS c JOIN country_translations AS ct
		 ON c.id = ct.country_id JOIN languages AS l
		 ON l.id = ct.language_id AND l.code = $1`, langCode)

	if err != nil {
		return nil, xerrors.Errorf(errFormat, langCode, "execute select query", err)
	}

	return countries, nil
}
