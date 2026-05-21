package pieceRepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
)

func (r *Repository) DeletePiece(ctx context.Context, id uuid.UUID) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.pieces WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return errs.ClassifyPgError("delete piece", err)
	}

	if result.RowsAffected() == 0 {
		return errs.NewRepositoryNotFoundErr(fmt.Errorf("piece with ID %s not found", id), "piece")
	}

	return nil
}
