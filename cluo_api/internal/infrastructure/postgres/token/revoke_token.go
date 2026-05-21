package tokenRepository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
)

func (r *Repository) RevokeToken(ctx context.Context, tokenID uuid.UUID) error {
	query := fmt.Sprintf(`
		UPDATE %s.case_access_tokens
		SET revoked_at = $1
		WHERE id = $2 AND revoked_at IS NULL
	`, r.schema)

	tag, err := r.pool.Exec(ctx, query, time.Now(), tokenID)
	if err != nil {
		return errs.ClassifyPgError("revoke token", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("revoke token: %w", errs.ErrRepositoryNotFound)
	}

	return nil
}
