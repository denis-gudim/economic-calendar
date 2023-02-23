package data

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type CountriesRepository struct {
	Db *sqlx.DB
}

func NewCountriesRepository(db *sqlx.DB) *CountriesRepository {
	return &CountriesRepository{db}
}

func (r *CountriesRepository) GetCountriesByLanguage(ctx context.Context, langCode string) ([]Country, error) {
	countries := make([]Country, 0, 128)
	err := r.Db.SelectContext(ctx, &countries,
		`SELECT c.id, c.code, c.continent_code, c.currency, ct.title AS name
		 FROM countries AS c JOIN country_translations AS ct
		 ON c.id = ct.country_id JOIN languages AS l
		 ON l.id = ct.language_id AND l.code = $1`, langCode)
	if err != nil {
		return nil, fmt.Errorf("get countries by language '%s' failed: %w", langCode, err)
	}
	return countries, nil
}
