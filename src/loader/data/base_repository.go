package data

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type baseRepository struct {
	ConnectionString string
}

func (r *baseRepository) createConnection() (*sqlx.DB, error) {
	return sqlx.Connect("postgres", r.ConnectionString)
}
