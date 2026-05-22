package caseTypeService

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	casetype "github.com/hengadev/cluo_api/internal/domain/case_type"
	"github.com/hengadev/cluo_api/internal/common/errs"
)

func (s *Service) UpdateCaseType(ctx context.Context, id uuid.UUID, req *casetype.UpdateCaseTypeRequest) (*casetype.CaseTypeResponse, error) {
	if req.Name == "" {
		return nil, errs.NewInvalidValueErr("name is required")
	}
	ct, err := s.repo.GetCaseTypeByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get case type for update: %w", err)
	}
	ct.Name = req.Name
	ct.UpdatedAt = time.Now()
	if err := s.repo.UpdateCaseType(ctx, ct); err != nil {
		if errors.Is(err, errs.ErrUniqueViolation) {
			return nil, errs.NewConflictErr(fmt.Errorf("case type with name %q already exists", req.Name))
		}
		return nil, fmt.Errorf("update case type: %w", err)
	}
	return ct.ToResponse(), nil
}
