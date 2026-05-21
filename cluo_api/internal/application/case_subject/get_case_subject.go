package caseSubjectService

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/subject"
)

func (s *Service) GetCaseSubjectByID(ctx context.Context, id uuid.UUID) (*subject.CaseSubjectResponse, error) {
	encxVal, err := s.repo.GetCaseSubjectByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get case subject: %w", err)
	}

	plain, err := subject.DecryptCaseSubjectEncx(ctx, s.crypto, encxVal)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("case subject", err)
	}

	return plain.ToResponse(), nil
}
