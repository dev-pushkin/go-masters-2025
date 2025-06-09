package pg

import (
	"log/slog"

	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresDB struct {
	p *pgxpool.Pool
}

func NewPostgresDB(log *slog.Logger, pool *pgxpool.Pool) *PostgresDB {
	return &PostgresDB{p: pool}
}
