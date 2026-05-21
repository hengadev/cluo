package caseTypeService

import (
	"context"
	"fmt"

	casetype "github.com/hengadev/cluo_api/internal/domain/case_type"
)

func (s *Service) ListCaseTypes(ctx context.Context) ([]*casetype.CaseTypeResponse, error) {
	caseTypes, err := s.repo.ListCaseTypes(ctx)
	if err != nil {
		return nil, fmt.Errorf("list case types: %w", err)
	}
	responses := make([]*casetype.CaseTypeResponse, 0, len(caseTypes))
	for _, ct := range caseTypes {
		responses = append(responses, ct.ToResponse())
	}
	return responses, nil
}
