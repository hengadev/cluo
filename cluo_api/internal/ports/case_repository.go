package ports

import (
	"context"

	"github.com/google/uuid"
	caseDomain "github.com/hengadev/cluo_api/internal/domain/case"
)

// CaseRepository handles persistent storage for cases.
type CaseRepository interface {
	CreateCase(ctx context.Context, c *caseDomain.CaseEncx) error
	GetCaseByID(ctx context.Context, id uuid.UUID) (*caseDomain.CaseEncx, error)
	UpdateCase(ctx context.Context, c *caseDomain.CaseEncx) error
	DeleteCase(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, f caseDomain.CaseFilter, p caseDomain.Pagination) ([]*caseDomain.CaseEncx, int, error)
	ListByClient(ctx context.Context, clientID uuid.UUID, p caseDomain.Pagination) ([]*caseDomain.CaseEncx, int, error)
}
