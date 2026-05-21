package caseTypeService

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	casetype "github.com/hengadev/cluo_api/internal/domain/case_type"
)

func (s *Service) UpdateCaseType(ctx context.Context, id uuid.UUID, req *casetype.UpdateCaseTypeRequest) (*casetype.CaseTypeResponse, error) {
	ct, err := s.repo.GetCaseTypeByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get case type for update: %w", err)
	}
	ct.Name = req.Name
	ct.UpdatedAt = time.Now()
	if err := s.repo.UpdateCaseType(ctx, ct); err != nil {
		return nil, fmt.Errorf("update case type: %w", err)
	}
	return ct.ToResponse(), nil
}
