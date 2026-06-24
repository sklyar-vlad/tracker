package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/sklyar-vlad/selfDev/internal/config"
)

func NewPostgres(ctx context.Context, config config.ConfigDatabase) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, config.DatabaseURL)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
