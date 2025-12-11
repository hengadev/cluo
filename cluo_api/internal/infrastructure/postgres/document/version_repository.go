package documentRepository

import (
	"github.com/hengadev/cluo_api/internal/ports"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VersionRepository struct {
	pool *pgxpool.Pool
}

func NewVersionRepository(pool *pgxpool.Pool) ports.DocumentVersionRepository {
	return &VersionRepository{pool: pool}
}
