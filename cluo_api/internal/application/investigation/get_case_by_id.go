package investigationService

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/investigation"
)

func (s *CaseService) GetCaseByID(ctx context.Context, request *investigation.GetCaseByIDRequest) (*investigation.CaseResponse, error) {
	caseEncx, err := s.repo.GetCaseByID(ctx, request.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get case by ID: %w", err)
	}

	// Decrypt c data
	c, err := investigation.DecryptInvestigationEncx(ctx, s.crypto, caseEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("case", err)
	}
	return c.ToResponse(), nil
}
