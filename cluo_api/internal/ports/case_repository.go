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
	// UpdateCase(ctx context.Context, c *caseDomain.CaseEncx) error
	// DeleteCase(ctx context.Context, id uuid.UUID) error
	// ListByClient(ctx context.Context, clientID string, p Pagination, f caseDomainFilter) ([]*caseDomain.CaseEncx, int, error) // returns list + total
	// ListByInvestigator(ctx context.Context, investigatorID string, p Pagination) ([]*caseDomain.CaseEncx, int, error)
	// Search(ctx context.Context, query string, p Pagination) ([]*caseDomain.CaseEncx, int, error)
}
