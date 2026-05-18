package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/investigation"
)

// CaseRepository handles persistent storage for cases.
type CaseRepository interface {
	CreateCase(ctx context.Context, c *investigation.InvestigationEncx) error
	GetCaseByID(ctx context.Context, id uuid.UUID) (*investigation.InvestigationEncx, error)
	UpdateCase(ctx context.Context, c *investigation.InvestigationEncx) error
	DeleteCase(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, f investigation.Filter, p investigation.Pagination) ([]*investigation.InvestigationEncx, int, error)
	ListByClient(ctx context.Context, clientID uuid.UUID, p investigation.Pagination) ([]*investigation.InvestigationEncx, int, error)
	ExistsCase(ctx context.Context, caseID uuid.UUID) (bool, error)
}
