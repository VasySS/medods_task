package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, connString string) (*Repository, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать подключение к postgres: %v", err)
	}

	return &Repository{
		pool: pool,
	}, nil
}

func (r *Repository) Close() {
	r.pool.Close()
}
