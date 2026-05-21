package caseTypeRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	casetype "github.com/hengadev/cluo_api/internal/domain/case_type"
)

func (r *Repository) UpdateCaseType(ctx context.Context, ct *casetype.CaseType) error {
	query := fmt.Sprintf(`
		UPDATE %s.case_types
		SET name = $1, updated_at = $2
		WHERE id = $3
	`, r.schema)
	tag, err := r.pool.Exec(ctx, query, ct.Name, ct.UpdatedAt, ct.ID)
	if err != nil {
		return errs.ClassifyPgError("update case type", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("update case type: %w", errs.ErrRepositoryNotFound)
	}
	return nil
}
