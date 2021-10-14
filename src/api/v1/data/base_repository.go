package data

import (
	"context"

	"github.com/denis-gudim/economic-calendar/api/app"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type baseRepository struct {
	config app.Config
}

func (r *baseRepository) connectDB(ctx context.Context) (*sqlx.DB, error) {
	return sqlx.ConnectContext(ctx, "postgres", r.config.DB.ConnectionString)
}
