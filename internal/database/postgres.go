package database

import (
	"context"
	"fmt"
	"time"

	"inforce_task/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(ctx context.Context, cfg config.Database) (*pgxpool.Pool, error) {
	const fn = "database.NewPostgresPool"

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode,
	)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: parse config: %w", fn, err)
	}

	poolConfig.MaxConns = cfg.PoolConf.MaxConns
	poolConfig.MinConns = cfg.PoolConf.MinConns
	poolConfig.MaxConnIdleTime = cfg.PoolConf.MaxConnIdleTime

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("%s: create pool: %w", fn, err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("%s: ping db: %w", fn, err)
	}

	return pool, nil
}
