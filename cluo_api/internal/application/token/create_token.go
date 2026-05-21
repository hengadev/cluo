package tokenService

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/token"
)

func (s *Service) CreateToken(ctx context.Context, caseID uuid.UUID) (*token.CreateTokenResponse, error) {
	exists, err := s.caseRepo.ExistsCase(ctx, caseID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify case existence: %w", err)
	}
	if !exists {
		return nil, errs.NewNotFoundErr(fmt.Errorf("case %s not found", caseID), "case")
	}

	rawToken, tokenHash, err := token.GenerateRawToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	now := time.Now()
	t := &token.Token{
		ID:        uuid.New(),
		CaseID:    caseID,
		TokenHash: tokenHash,
		ExpiresAt: now.Add(token.TokenExpiryDays * 24 * time.Hour),
		RevokedAt: nil,
		CreatedAt: now,
	}

	if err := s.repo.CreateToken(ctx, t); err != nil {
		return nil, errs.NewNotCreatedErr(err, "token")
	}

	return &token.CreateTokenResponse{
		ID:        t.ID.String(),
		CaseID:    t.CaseID.String(),
		RawToken:  rawToken,
		ExpiresAt: t.ExpiresAt,
		CreatedAt: t.CreatedAt,
	}, nil
}
