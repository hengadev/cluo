package ports

import (
	"context"

	"github.com/hengadev/cluo_api/internal/domain/case"
)

// CaseRepository handles persistent storage for cases.
type CaseRepository interface {
	CreateCase(ctx context.Context, c *caseDomain.Case) error
	GetCaseByID(ctx context.Context, id string) (*caseDomain.Case, error)
	UpdateCase(ctx context.Context, c *caseDomain.Case) error
	DeleteCase(ctx context.Context, id string) error
	// ListByClient(ctx context.Context, clientID string, p Pagination, f caseDomainFilter) ([]*caseDomain.Case, int, error) // returns list + total
	// ListByInvestigator(ctx context.Context, investigatorID string, p Pagination) ([]*case.Case, int, error)
	// Search(ctx context.Context, query string, p Pagination) ([]*case.Case, int, error)
}
