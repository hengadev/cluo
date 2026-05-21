package tokenRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/token"
)

func (r *Repository) GetTokenByHash(ctx context.Context, tokenHash string) (*token.Token, error) {
	query := fmt.Sprintf(`
		SELECT id, case_id, token_hash, expires_at, revoked_at, created_at
		FROM %s.case_access_tokens
		WHERE token_hash = $1
	`, r.schema)

	t := &token.Token{}
	err := r.pool.QueryRow(ctx, query, tokenHash).Scan(
		&t.ID,
		&t.CaseID,
		&t.TokenHash,
		&t.ExpiresAt,
		&t.RevokedAt,
		&t.CreatedAt,
	)
	if err != nil {
		return nil, errs.ClassifyPgError("get token by hash", err)
	}

	return t, nil
}
