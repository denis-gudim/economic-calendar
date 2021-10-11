package data

import (
	"golang.org/x/xerrors"
)

type LanguagesRepository struct {
	baseRepository
}

func (r *LanguagesRepository) GetAll() (languages []Language, err error) {
	fmtError := func(msg string, err error) error {
		return xerrors.Errorf("get all languages failed: %s: %w", msg, err)
	}

	db, err := r.createConnection()

	if err != nil {
		return nil, fmtError("create connection", err)
	}

	defer db.Close()

	rows, err := r.initQueryBuilder().
		Select("*").
		From("languages").
		RunWith(db).
		Query()

	if err != nil {
		return nil, fmtError("execute select query", err)
	}

	languages = make([]Language, 0, 24)

	for rows.Next() {
		lang := Language{}

		err = rows.Scan(
			&lang.Id,
			&lang.Code,
			&lang.Name,
			&lang.NativeName,
			&lang.Domain,
		)

		if err != nil {
			return nil, fmtError("scan row", err)
		}

		languages = append(languages, lang)
	}

	if err != nil {
		return nil, fmtError("execute select query", err)
	}

	return
}
