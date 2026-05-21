package subject

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
)

func (r *Repository) DeleteCaseSubject(ctx context.Context, id uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM %s.case_subjects WHERE id = $1`, r.schema)
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return errs.ClassifyPgError("delete case subject", err)
	}
	if result.RowsAffected() == 0 {
		return errs.NewRepositoryNotFoundErr(fmt.Errorf("case subject %s not found", id), "case_subject")
	}
	return nil
}
