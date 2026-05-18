package ports

import (
	"context"

	"github.com/hengadev/cluo_api/internal/domain/investigation"
)

type CaseService interface {
	CreateCase(ctx context.Context, r *investigation.CreateCaseRequest) (*investigation.CaseResponse, error)
	GetCaseByID(ctx context.Context, r *investigation.GetCaseByIDRequest) (*investigation.CaseResponse, error)
	UpdateCase(ctx context.Context, r *investigation.UpdateCaseRequest) (*investigation.CaseResponse, error)
	DeleteCase(ctx context.Context, r *investigation.DeleteCaseByIDRequest) error
	List(ctx context.Context, r *investigation.ListCasesRequest) (*investigation.ListCasesResponse, error)
	ListByClient(ctx context.Context, r *investigation.ListByClientRequest) (*investigation.ListCasesResponse, error)
}
