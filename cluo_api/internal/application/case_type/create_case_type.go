package caseTypeService

import (
	"context"
	"errors"
	"fmt"

	casetype "github.com/hengadev/cluo_api/internal/domain/case_type"
	"github.com/hengadev/cluo_api/internal/common/errs"
)

func (s *Service) CreateCaseType(ctx context.Context, req *casetype.CreateCaseTypeRequest) (*casetype.CaseTypeResponse, error) {
	if req.Name == "" {
		return nil, errs.NewInvalidValueErr("name is required")
	}
	ct := casetype.New(req.Name)
	if err := s.repo.CreateCaseType(ctx, ct); err != nil {
		if errors.Is(err, errs.ErrUniqueViolation) {
			return nil, errs.NewConflictErr(fmt.Errorf("case type with name %q already exists", req.Name))
		}
		return nil, fmt.Errorf("create case type: %w", err)
	}
	return ct.ToResponse(), nil
}
