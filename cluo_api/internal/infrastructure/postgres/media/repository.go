package mediaRepository

import (
	"context"

	"github.com/hengadev/cluo_api/internal/ports"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const Schema = "media"

type Repository struct {
	pool   *pgxpool.Pool
	schema string
}

func New(ctx context.Context, pool *pgxpool.Pool) ports.MediaRepository {
	return &Repository{pool: pool, schema: Schema}
}
