package tokenService

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/token"
)

// ValidateToken hashes the raw token, looks it up, checks IsValid(), and returns the caseID.
func (s *Service) ValidateToken(ctx context.Context, rawToken string) (uuid.UUID, error) {
	tokenHash := token.HashToken(rawToken)

	t, err := s.repo.GetTokenByHash(ctx, tokenHash)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to look up token: %w", err)
	}

	if !t.IsValid() {
		return uuid.Nil, errs.NewExpiredTokenErr("access", fmt.Errorf("token is expired or revoked"))
	}

	return t.CaseID, nil
}
