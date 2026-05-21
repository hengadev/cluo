package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/rapport"
)

type RapportRepository interface {
	CreateRapport(ctx context.Context, r *rapport.RapportEncx) error
	GetRapportByCaseID(ctx context.Context, caseID uuid.UUID) (*rapport.RapportEncx, error)
	UpdateRapport(ctx context.Context, r *rapport.RapportEncx) error
	DeleteRapport(ctx context.Context, caseID uuid.UUID) error
}
