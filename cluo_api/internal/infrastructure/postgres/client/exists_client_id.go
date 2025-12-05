package clientRepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
)

func (r *Repository) ExistsClient(ctx context.Context, clientID uuid.UUID) (bool, error) {
	query := fmt.Sprintf(`
		SELECT EXISTS(
			SELECT 1 FROM %s.clients
			WHERE id = $1
		)
	`, r.schema)

	var exists bool
	err := r.pool.QueryRow(ctx, query, clientID).Scan(&exists)
	if err != nil {
		return false, errs.ClassifyPgError("check if client exists by ID", err)
	}

	return exists, nil
}
