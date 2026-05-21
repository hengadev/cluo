package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/rapport"
)

type RapportService interface {
	CreateRapport(ctx context.Context, req *rapport.CreateRapportRequest) (*rapport.RapportResponse, error)
	GetRapportByCaseID(ctx context.Context, caseID uuid.UUID) (*rapport.RapportResponse, error)
	UpdateRapport(ctx context.Context, req *rapport.UpdateRapportRequest) (*rapport.RapportResponse, error)
	DeleteRapport(ctx context.Context, caseID uuid.UUID) error
}
