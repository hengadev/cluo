package caseTypeRepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
)

func (r *Repository) DeleteCaseType(ctx context.Context, id uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM %s.case_types WHERE id = $1`, r.schema)
	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return errs.ClassifyPgError("delete case type", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("delete case type: %w", errs.ErrRepositoryNotFound)
	}
	return nil
}
