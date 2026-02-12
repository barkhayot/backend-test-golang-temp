package pgdb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	connectionTimeout = 5 * time.Second
)

func New(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, connectionTimeout)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

func Close(db *pgxpool.Pool) {
	if db != nil {
		log.Println("closing database connection...")
		db.Close()
	}
}
