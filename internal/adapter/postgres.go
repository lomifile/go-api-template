package adapter

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/lomifile/api/pkg/postgres"
)

type PostgresAdapter struct {
	*sqlx.DB
}

func NewPostgresAdapter(pg *postgres.Postgres) (*PostgresAdapter, error) {
	sqlDB := stdlib.OpenDBFromPool(pg.Pool)
	dbx := sqlx.NewDb(sqlDB, "pgx")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := dbx.PingContext(ctx); err != nil {
		return nil, err
	}

	return &PostgresAdapter{DB: dbx}, nil
}
