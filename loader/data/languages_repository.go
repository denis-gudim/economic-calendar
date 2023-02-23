package data

import (
	"context"
	"fmt"
)

type LanguagesRepository struct {
	baseRepository
}

func (r *LanguagesRepository) GetAll(ctx context.Context) (languages []Language, err error) {
	rows, err := r.initQueryBuilder().
		Select("*").
		From("languages").
		RunWith(r.db).
		QueryContext(ctx)

	if err != nil {
		return nil, fmt.Errorf("execute select query error: %w", err)
	}

	defer rows.Close()

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
			return nil, fmt.Errorf("scan row error: %w", err)
		}
		languages = append(languages, lang)
	}

	if err != nil {
		return nil, fmt.Errorf("execute select query error: %w", err)
	}

	return
}
