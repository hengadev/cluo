package caseSubjectService

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/subject"
)

func (s *Service) ListCaseSubjects(ctx context.Context, page, pageSize int) ([]*subject.CaseSubjectResponse, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	encxList, total, err := s.repo.ListCaseSubjects(ctx, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("list case subjects: %w", err)
	}

	responses := make([]*subject.CaseSubjectResponse, 0, len(encxList))
	for _, enc := range encxList {
		plain, err := subject.DecryptCaseSubjectEncx(ctx, s.crypto, enc)
		if err != nil {
			return nil, 0, errs.NewNotDecryptedErr("case subject", err)
		}
		responses = append(responses, plain.ToResponse())
	}

	return responses, total, nil
}
