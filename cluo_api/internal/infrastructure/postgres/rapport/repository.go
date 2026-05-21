package rapportRepository

import (
	"context"

	"github.com/hengadev/cluo_api/internal/ports"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool   *pgxpool.Pool
	schema string
}

func New(ctx context.Context, pool *pgxpool.Pool) ports.RapportRepository {
	return &Repository{pool: pool, schema: "cases"}
}
