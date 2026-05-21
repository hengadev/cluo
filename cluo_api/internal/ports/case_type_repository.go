package ports

import (
	"context"

	"github.com/google/uuid"
	casetype "github.com/hengadev/cluo_api/internal/domain/case_type"
)

type CaseTypeRepository interface {
	CreateCaseType(ctx context.Context, ct *casetype.CaseType) error
	GetCaseTypeByID(ctx context.Context, id uuid.UUID) (*casetype.CaseType, error)
	ListCaseTypes(ctx context.Context) ([]*casetype.CaseType, error)
	UpdateCaseType(ctx context.Context, ct *casetype.CaseType) error
	DeleteCaseType(ctx context.Context, id uuid.UUID) error
}
