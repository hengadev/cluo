package caseTypeService

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	casetype "github.com/hengadev/cluo_api/internal/domain/case_type"
)

func (s *Service) GetCaseTypeByID(ctx context.Context, id uuid.UUID) (*casetype.CaseTypeResponse, error) {
	ct, err := s.repo.GetCaseTypeByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get case type: %w", err)
	}
	return ct.ToResponse(), nil
}
