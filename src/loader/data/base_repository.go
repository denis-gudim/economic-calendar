package data

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
)

type baseRepository struct {
	ConnectionString string
}

func (r *baseRepository) createConnection() (*sql.DB, error) {
	return sql.Open("postgres", r.ConnectionString)
}

func (r *baseRepository) initQueryBuilder() sq.StatementBuilderType {
	return sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
}
