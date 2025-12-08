package ports

import (
	"context"

	"github.com/hengadev/cluo_api/internal/domain/case"
)

type CaseService interface {
	CreateCase(ctx context.Context, r *caseDomain.CreateCaseRequest) (*caseDomain.CaseResponse, error)
	GetCaseByID(ctx context.Context, r *caseDomain.GetCaseByIDRequest) (*caseDomain.CaseResponse, error)
	UpdateCase(ctx context.Context, r *caseDomain.UpdateCaseRequest) (*caseDomain.CaseResponse, error)
	DeleteCase(ctx context.Context, r *caseDomain.DeleteCaseByIDRequest) error
	List(ctx context.Context, r *caseDomain.ListCasesRequest) (*caseDomain.ListCasesResponse, error)
	ListByClient(ctx context.Context, r *caseDomain.ListByClientRequest) (*caseDomain.ListCasesResponse, error)
}
