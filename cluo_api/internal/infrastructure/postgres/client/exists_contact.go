package clientRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/google/uuid"
)

func (r *Repository) ExistsContact(ctx context.Context, contactID uuid.UUID) (bool, error) {
	query := fmt.Sprintf(`
		SELECT EXISTS(
			SELECT 1 FROM %s.contacts WHERE id = $1
		)
	`, r.schema)

	var exists bool
	err := r.pool.QueryRow(ctx, query, contactID).Scan(&exists)
	if err != nil {
		return false, errs.ClassifyPgError("exists contact", err)
	}

	return exists, nil
}