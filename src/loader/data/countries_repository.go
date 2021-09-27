package data

import (
	"economic-calendar/loader/investing/data"
	"fmt"
)

type CountriesRepository struct {
	baseRepository
}

func (r *CountriesRepository) GetAll() (countries []Country, err error) {
	db, err := r.createConnection()

	if err != nil {
		return nil, fmt.Errorf("get all countries create db connection error: %w", err)
	}

	defer db.Close()

	countries = make([]Country, 0, 248)

	rows, err := db.Queryx(
		`SELECT c.*, ct.language_code, ct.title
		 FROM countries AS c LEFT JOIN country_translations AS ct
		 ON c.code = ct.country_code
		 ORDER BY c.code`)

	if err != nil {
		return nil, fmt.Errorf("get all countries query error: %w", err)
	}

	var (
		languageCode  *string
		languageTitle *string
		curr          Country
		prevCode      string
		trans         map[string]string
	)

	for rows.Next() {
		err = rows.Scan(
			&curr.Code,
			&curr.ContinentCode,
			&curr.Name,
			&curr.Currency,
			&curr.InvestingId,
			&languageCode,
			&languageTitle,
		)

		if err != nil {
			return nil, fmt.Errorf("get all countries scan row error: %w", err)
		}

		if curr.Code != prevCode {
			trans = make(map[string]string, len(data.InvestingLanguagesMap))
			curr.Translations = trans
			countries = append(countries, curr)
			prevCode = curr.Code
		}

		if languageCode != nil {
			trans[*languageCode] = *languageTitle
		}
	}

	return
}

func (r *CountriesRepository) Save(country *Country) error {
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
		`UPDATE countries 
		 SET code=:code 
		   , continent_code=:continent_code
		   , currency=:currency
		   , investing_id=:investing_id
		   , name=:name
		 WHERE code=:code`, country)

	if err != nil {
		tx.Rollback() // TODO: Add rollback error handling
		return fmt.Errorf("save country execute update error: %w", err)
	}

	for lang, title := range country.Translations {
		_, err = tx.Exec(
			`INSERT INTO country_translations (country_code, language_code, title)
			 VALUES($1, $2, $3)
			 ON CONFLICT (country_code, language_code) DO UPDATE
			 SET country_code=$1, language_code=$2, title=$3`,
			country.Code, lang, title)

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
