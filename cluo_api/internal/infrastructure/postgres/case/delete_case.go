package caseRepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
)

func (r *Repository) DeleteCase(ctx context.Context, caseID uuid.UUID) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.cases WHERE id = $1
	`, r.schema)

	cmdTag, err := r.pool.Exec(ctx, query, caseID)
	if err != nil {
		return errs.ClassifyPgError("delete case", err)
	}

	// Check if a row was actually deleted
	if cmdTag.RowsAffected() == 0 {
		return errs.NewRepositoryNotFoundErr(nil, "case for deletion")
	}

	return nil
}