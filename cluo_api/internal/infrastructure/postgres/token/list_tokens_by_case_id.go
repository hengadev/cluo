package tokenRepository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/token"
)

func (r *Repository) ListTokensByCaseID(ctx context.Context, caseID uuid.UUID) ([]*token.Token, error) {
	query := fmt.Sprintf(`
		SELECT id, case_id, token_hash, expires_at, revoked_at, created_at
		FROM %s.case_access_tokens
		WHERE case_id = $1
		ORDER BY created_at DESC
	`, r.schema)

	rows, err := r.pool.Query(ctx, query, caseID)
	if err != nil {
		return nil, errs.ClassifyPgError("list tokens by case id", err)
	}
	defer rows.Close()

	var tokens []*token.Token
	for rows.Next() {
		t := &token.Token{}
		if err := rows.Scan(
			&t.ID,
			&t.CaseID,
			&t.TokenHash,
			&t.ExpiresAt,
			&t.RevokedAt,
			&t.CreatedAt,
		); err != nil {
			return nil, errs.ClassifyPgError("scan token row", err)
		}
		tokens = append(tokens, t)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.ClassifyPgError("iterate token rows", err)
	}

	return tokens, nil
}
