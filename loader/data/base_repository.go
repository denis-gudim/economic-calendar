package data

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

type baseRepository struct {
	db *sql.DB
}

func (r *baseRepository) initQueryBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}
