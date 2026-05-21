package caseTypeRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	casetype "github.com/hengadev/cluo_api/internal/domain/case_type"
)

func (r *Repository) CreateCaseType(ctx context.Context, ct *casetype.CaseType) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.case_types (id, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
	`, r.schema)
	_, err := r.pool.Exec(ctx, query, ct.ID, ct.Name, ct.CreatedAt, ct.UpdatedAt)
	if err != nil {
		return errs.ClassifyPgError("create case type", err)
	}
	return nil
}
