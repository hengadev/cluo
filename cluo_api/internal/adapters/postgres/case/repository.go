package caseRepository

import (
	"context"

	"github.com/hengadev/cluo_api/internal/ports"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Repository struct {
	pool   *pgxpool.Pool
	schema string
}

func New(ctx context.Context, pool *pgxpool.Pool) ports.CaseRepository {
	return &Repository{pool: pool, schema: "cases"}
}
