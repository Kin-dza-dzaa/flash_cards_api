// Package postgres implements potgres connection
package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ConnPool struct {
	Pool    *pgxpool.Pool
	Builder sq.StatementBuilderType
}

func (p *ConnPool) Close() {
	p.Pool.Close()
}

func New(pgurl string, maxPoolSize int) (*ConnPool, error) {
	poolConfig, err := pgxpool.ParseConfig(pgurl)
	if err != nil {
		return nil, fmt.Errorf("ConnPool - NewPool - pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MaxConns = int32(maxPoolSize)

	pool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("ConnPool - NewPool - ConnectConfig: %w", err)
	}

	postgres := new(ConnPool)
	postgres.Builder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	postgres.Pool = pool

	return postgres, nil
}
