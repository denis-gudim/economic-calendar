package data

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

type CountriesRepository struct {
	baseRepository
}

func NewCountriesRepository(db *sql.DB) *CountriesRepository {
	r := CountriesRepository{}
	r.db = db
	return &r
}

func (r *CountriesRepository) GetAll(ctx context.Context) (countries []Country, err error) {
	countries = make([]Country, 0, 100)
	rows, err := r.initQueryBuilder().
		Select("c.*, ct.language_id, ct.title").
		From("countries c").
		LeftJoin("country_translations ct ON c.id = ct.country_id").
		OrderBy("c.id").
		RunWith(r.db).
		QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("execute select query error: %w", err)
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
			return nil, fmt.Errorf("scan row error: %w", err)
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

func (r *CountriesRepository) Save(ctx context.Context, c Country) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return fmt.Errorf("create db transaction error: %w", err)
	}
	defer func() {
		if tx != nil && err != nil {
			if rerr := tx.Rollback(); rerr != nil {
				err = rerr
			}
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

	_, err = upsertQuery.RunWith(tx).ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("execute upsert query error: %w", err)
	}

	deleteQuery := r.initQueryBuilder().
		Delete("country_translations").
		Where(sq.Eq{"country_id": c.Id})

	_, err = deleteQuery.RunWith(tx).ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("delete translations query error: %w", err)
	}

	for langId, title := range c.NameTranslations {

		insertQuery := r.initQueryBuilder().
			Insert("country_translations").
			Columns("country_id", "language_id", "title").
			Values(c.Id, langId, title)

		_, err = insertQuery.RunWith(tx).ExecContext(ctx)

		if err != nil {
			return fmt.Errorf("execute insert translation query error: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
