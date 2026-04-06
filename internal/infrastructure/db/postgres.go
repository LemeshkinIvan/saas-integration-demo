package db

import (
	"context"
	config "daos_core/internal/constants"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Cfg  *config.PostgresConfig
	Pool *pgxpool.Pool
}

var ErrGetFromConfig = errors.New("invalid config")

func Connect(cfg *config.PostgresConfig) (*Postgres, error) {
	url := cfg.GetConnectionString()

	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, err
	}

	if pool == nil {
		return nil, errors.New("Connect: pool is nil")
	}

	db := &Postgres{
		Pool: pool,
	}

	return db, nil
}

func (db *Postgres) Disconnect() {
	db.Pool.Close()
}
