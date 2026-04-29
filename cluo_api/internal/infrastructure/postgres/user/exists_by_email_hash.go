package userRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
)

func (r *Repository) ExistsByEmailHash(ctx context.Context, emailHash string) (bool, error) {
	query := fmt.Sprintf(`
		SELECT EXISTS(
			SELECT 1 FROM %s.users
			WHERE email_hash = $1
		)
	`, r.schema)

	var exists bool
	err := r.pool.QueryRow(ctx, query, emailHash).Scan(&exists)
	if err != nil {
		return false, errs.ClassifyPgError("check if user exists by email hash", err)
	}

	return exists, nil
}
