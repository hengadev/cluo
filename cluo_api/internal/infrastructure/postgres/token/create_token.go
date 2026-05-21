package tokenRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/token"
)

func (r *Repository) CreateToken(ctx context.Context, t *token.Token) error {
	if t == nil {
		return fmt.Errorf("token cannot be nil")
	}

	query := fmt.Sprintf(`
		INSERT INTO %s.case_access_tokens (id, case_id, token_hash, expires_at, revoked_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, r.schema)

	_, err := r.pool.Exec(ctx, query,
		t.ID,
		t.CaseID,
		t.TokenHash,
		t.ExpiresAt,
		t.RevokedAt,
		t.CreatedAt,
	)
	if err != nil {
		return errs.ClassifyPgError("create token", err)
	}

	return nil
}
