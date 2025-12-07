package ports

import (
	"context"

	"github.com/hengadev/cluo_api/internal/domain/case"
)

type CaseService interface {
	CreateCase(ctx context.Context, r *caseDomain.CreateCaseRequest) (*caseDomain.CaseResponse, error)
	GetCaseByID(ctx context.Context, r *caseDomain.GetCaseByIDRequest) (*caseDomain.CaseResponse, error)
	UpdateCase(ctx context.Context, c *caseDomain.Case) error
	DeleteCase(ctx context.Context, r *caseDomain.DeleteCaseByIDRequest) error
}
