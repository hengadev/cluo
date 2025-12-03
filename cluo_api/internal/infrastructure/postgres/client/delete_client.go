package clientRepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
)

func (r *Repository) DeleteClient(ctx context.Context, clientID uuid.UUID) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.clients
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query, clientID)
	if err != nil {
		return errs.ClassifyPgError("delete client", err)
	}

	// Check if any row was actually deleted
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.NewRepositoryNotFoundErr(nil, "client for deletion")
	}
	return nil
}