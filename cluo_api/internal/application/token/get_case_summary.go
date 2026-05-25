package tokenService

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/investigation"
	"github.com/hengadev/cluo_api/internal/domain/token"
)

func (s *Service) GetCaseSummaryByToken(ctx context.Context, rawToken string) (*investigation.PortalCaseResponse, error) {
	tokenHash := token.HashToken(rawToken)

	t, err := s.repo.GetTokenByHash(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("failed to look up token: %w", err)
	}

	if !t.IsValid() {
		return nil, errs.NewExpiredTokenErr("access", fmt.Errorf("token is expired or revoked"))
	}

	caseEncx, err := s.caseRepo.GetCaseByID(ctx, t.CaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get case: %w", err)
	}

	c, err := investigation.DecryptInvestigationEncx(ctx, s.crypto, caseEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("case", err)
	}

	return &investigation.PortalCaseResponse{
		CaseResponse:   c.ToResponse(),
		TokenExpiresAt: t.ExpiresAt,
	}, nil
}
