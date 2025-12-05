package clientRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
)

func (r *Repository) ExistsByClientID(ctx context.Context, clientID string) (bool, error) {
	query := fmt.Sprintf(`
		SELECT EXISTS(
			SELECT 1 FROM %s.contacts 
			WHERE client_id = $1
		)
	`, r.schema)

	var exists bool
	err := r.pool.QueryRow(ctx, query, clientID).Scan(&exists)
	if err != nil {
		return false, errs.ClassifyPgError("check if contact exists by client ID", err)
	}

	return exists, nil
}
