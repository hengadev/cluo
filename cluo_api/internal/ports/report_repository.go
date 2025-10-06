package ports

import (
	"context"

	"github.com/hengadev/cluo_api/internal/domain"
)

// ReportRepository stores tiptap JSON content and optional versioning.
type ReportRepository interface {
	CreateOrUpdate(ctx context.Context, r *domain.Report) error
	GetByID(ctx context.Context, id string) (*domain.Report, error)
	GetByCaseID(ctx context.Context, caseID string) (*domain.Report, error)
	GetVersions(ctx context.Context, caseID string) ([]*domain.Report, error)
	Delete(ctx context.Context, id string) error
}
