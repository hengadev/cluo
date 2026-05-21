package rapportRepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
)

func (r *Repository) DeleteRapport(ctx context.Context, caseID uuid.UUID) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.rapports WHERE case_id = $1
	`, r.schema)

	cmdTag, err := r.pool.Exec(ctx, query, caseID)
	if err != nil {
		return errs.ClassifyPgError("delete rapport", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return errs.NewRepositoryNotFoundErr(nil, "rapport for deletion")
	}

	return nil
}
