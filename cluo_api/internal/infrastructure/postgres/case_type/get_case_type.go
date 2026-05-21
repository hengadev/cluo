package caseTypeRepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	casetype "github.com/hengadev/cluo_api/internal/domain/case_type"
)

func (r *Repository) GetCaseTypeByID(ctx context.Context, id uuid.UUID) (*casetype.CaseType, error) {
	query := fmt.Sprintf(`SELECT id, name, created_at, updated_at FROM %s.case_types WHERE id = $1`, r.schema)
	ct := &casetype.CaseType{}
	err := r.pool.QueryRow(ctx, query, id).Scan(&ct.ID, &ct.Name, &ct.CreatedAt, &ct.UpdatedAt)
	if err != nil {
		return nil, errs.ClassifyPgError("get case type by id", err)
	}
	return ct, nil
}
