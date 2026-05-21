package pieceRepository

import (
	"context"

	"github.com/hengadev/cluo_api/internal/ports"
	"github.com/jackc/pgx/v5/pgxpool"
)

const schema = "cases"

// Repository implements ports.PieceRepository backed by PostgreSQL.
type Repository struct {
	pool   *pgxpool.Pool
	schema string
}

// New creates a new PieceRepository.
func New(ctx context.Context, pool *pgxpool.Pool) ports.PieceRepository {
	return &Repository{pool: pool, schema: schema}
}
