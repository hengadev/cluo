package tokenService

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/investigation"
)

func (s *Service) GetCaseSummaryByToken(ctx context.Context, rawToken string) (*investigation.CaseResponse, error) {
	caseID, err := s.ValidateToken(ctx, rawToken)
	if err != nil {
		return nil, err
	}

	caseEncx, err := s.caseRepo.GetCaseByID(ctx, caseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get case: %w", err)
	}

	c, err := investigation.DecryptInvestigationEncx(ctx, s.crypto, caseEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("case", err)
	}

	return c.ToResponse(), nil
}
