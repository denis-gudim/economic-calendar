package data

import (
	"economic-calendar/loader/app"
	"fmt"
)

type CountriesRepository struct {
	baseRepository
}

func NewCountriesRepository(cnf app.Config) *CountriesRepository {
	r := CountriesRepository{}
	r.ConnectionString = cnf.DB.ConnectionString
	return &r
}

func (r *CountriesRepository) GetAll() (countries []Country, err error) {
	db, err := r.createConnection()

	if err != nil {
		return nil, fmt.Errorf("get all countries create db connection error: %w", err)
	}

	defer db.Close()

	countries = make([]Country, 0, 248)

	rows, err := db.Queryx(
		`SELECT c.*, ct.language_id, ct.title
		 FROM countries AS c LEFT JOIN country_translations AS ct
		 ON c.id = ct.country_id
		 ORDER BY c.id`)

	if err != nil {
		return nil, fmt.Errorf("get all countries query error: %w", err)
	}

	var (
		langId    *int
		langTitle *string
		curr      Country
		prevCode  string
		trans     Translations
	)

	for rows.Next() {
		err = rows.Scan(
			&curr.Id,
			&curr.Code,
			&curr.ContinentCode,
			&curr.Name,
			&curr.Currency,
			&langId,
			&langTitle,
		)

		if err != nil {
			return nil, fmt.Errorf("get all countries scan row error: %w", err)
		}

		if curr.Code != prevCode {
			trans = Translations{}
			curr.Translations = trans
			countries = append(countries, curr)
			prevCode = curr.Code
		}

		if langId != nil {
			trans[*langId] = *langTitle
		}
	}

	return
}

func (r *CountriesRepository) Save(country Country) error {
	db, err := r.createConnection()

	if err != nil {
		return fmt.Errorf("save country create db connection error: %w", err)
	}

	defer db.Close()

	tx, err := db.Beginx()

	if err != nil {
		return fmt.Errorf("save country create db transaction error: %w", err)
	}

	_, err = tx.NamedExec(
		`INSERT INTO countries (id, code, continent_code, currency, name)
		 VALUES (:id, :code, :continent_code, :currency, :name)
		 ON CONFLICT (id) DO UPDATE 
		 SET code = :code, continent_code = :continent_code, currency = :currency, name = :name`,
		country)

	if err != nil {
		tx.Rollback() // TODO: Add rollback error handling
		return fmt.Errorf("save country execute update error: %w", err)
	}

	_, err = tx.Exec(`DELETE FROM country_translations WHERE country_id=$1`, country.Id)

	if err != nil {
		tx.Rollback() // TODO: Add rollback error handling
		return fmt.Errorf("save country delete translations error: %w", err)
	}

	for langId, title := range country.Translations {
		_, err = tx.Exec(
			`INSERT INTO country_translations (country_id, language_id, title)
			 VALUES($1, $2, $3)`,
			country.Id, langId, title)

		if err != nil {
			tx.Rollback() //TODO: Add rollaback error handling
			return fmt.Errorf("save country execute update translation error: %w", err)
		}
	}

	err = tx.Commit()

	if err != nil {
		return fmt.Errorf("save country commit transaction error: %w", err)
	}

	return nil
}
