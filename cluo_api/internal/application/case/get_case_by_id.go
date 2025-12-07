package caseService

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	caseDomain "github.com/hengadev/cluo_api/internal/domain/case"
)

func (s *CaseService) GetCaseByID(ctx context.Context, request *caseDomain.GetCaseByIDRequest) (*caseDomain.CaseResponse, error) {
	caseEncx, err := s.repo.GetCaseByID(ctx, request.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get case by ID: %w", err)
	}

	// Decrypt c data
	c, err := caseDomain.DecryptCaseEncx(ctx, s.crypto, caseEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("case", err)
	}
	return c.ToResponse(), nil
}
