package clientRepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
)

func (r *Repository) DeleteContact(ctx context.Context, contactID uuid.UUID) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.contacts
		WHERE id = $1
	`, r.schema)

	result, err := r.pool.Exec(ctx, query, contactID)
	if err != nil {
		return errs.ClassifyPgError("delete contact", err)
	}

	// Check if any row was actually deleted
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return errs.NewRepositoryNotFoundErr(nil, "contact for deletion")
	}
	return nil
}
