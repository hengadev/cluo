package ports

import (
	"context"

	"github.com/google/uuid"
	casetype "github.com/hengadev/cluo_api/internal/domain/case_type"
)

type CaseTypeService interface {
	CreateCaseType(ctx context.Context, req *casetype.CreateCaseTypeRequest) (*casetype.CaseTypeResponse, error)
	GetCaseTypeByID(ctx context.Context, id uuid.UUID) (*casetype.CaseTypeResponse, error)
	ListCaseTypes(ctx context.Context) ([]*casetype.CaseTypeResponse, error)
	UpdateCaseType(ctx context.Context, id uuid.UUID, req *casetype.UpdateCaseTypeRequest) (*casetype.CaseTypeResponse, error)
	DeleteCaseType(ctx context.Context, id uuid.UUID) error
}
