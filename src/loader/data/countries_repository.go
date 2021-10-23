package data

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"golang.org/x/xerrors"
)

type CountriesRepository struct {
	baseRepository
}

func NewCountriesRepository(db *sql.DB) *CountriesRepository {
	r := CountriesRepository{}
	r.db = db
	return &r
}

func (r *CountriesRepository) GetAll() (countries []Country, err error) {

	fmtError := func(msg string, err error) error {
		return xerrors.Errorf("get all countries failed: %s: %w", msg, err)
	}
	countries = make([]Country, 0, 100)

	rows, err := r.initQueryBuilder().
		Select("c.*, ct.language_id, ct.title").
		From("countries c").
		LeftJoin("country_translations ct ON c.id = ct.country_id").
		OrderBy("c.id").
		RunWith(r.db).
		Query()

	if err != nil {
		return nil, fmtError("execute select query", err)
	}

	defer rows.Close()

	var (
		langId    *int
		langTitle *string
		curr      Country
		prevId    int
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
			return nil, fmtError("scan row", err)
		}

		if curr.Id != prevId {
			trans = Translations{}
			curr.NameTranslations = trans
			countries = append(countries, curr)
			prevId = curr.Id
		}

		if langId != nil {
			trans[*langId] = *langTitle
		}
	}

	return
}

func (r *CountriesRepository) Save(c Country) error {

	fmtError := func(msg string, err error) error {
		return xerrors.Errorf("save country failed: %s: %w", msg, err)
	}

	tx, err := r.db.Begin()

	if err != nil {
		return fmtError("create db transaction", err)
	}

	defer func() {
		if tx != nil && err != nil {
			tx.Rollback()
		}
	}()

	upsertQuery := r.initQueryBuilder().
		Insert("countries").
		Columns("id", "code", "continent_code", "currency", "name").
		Values(c.Id, c.Code, c.ContinentCode, c.Currency, c.Name).
		Suffix("ON CONFLICT (id) DO").
		SuffixExpr(
			sq.Update(" ").
				Set("code", c.Code).
				Set("continent_code", c.ContinentCode).
				Set("currency", c.Currency).
				Set("name", c.Name))

	_, err = upsertQuery.RunWith(tx).Exec()

	if err != nil {
		return fmtError("execute upsert country query", err)
	}

	deleteQuery := r.initQueryBuilder().
		Delete("country_translations").
		Where(sq.Eq{"country_id": c.Id})

	_, err = deleteQuery.RunWith(tx).Exec()

	if err != nil {
		return fmtError("execute delete country translations query", err)
	}

	for langId, title := range c.NameTranslations {

		insertQuery := r.initQueryBuilder().
			Insert("country_translations").
			Columns("country_id", "language_id", "title").
			Values(c.Id, langId, title)

		_, err = insertQuery.RunWith(tx).Exec()

		if err != nil {
			return fmtError("execute insert country translation query", err)
		}
	}

	err = tx.Commit()

	if err != nil {
		return fmtError("commit transaction", err)
	}

	return nil
}
