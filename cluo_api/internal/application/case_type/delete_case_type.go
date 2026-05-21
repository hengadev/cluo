package caseTypeService

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
)

func (s *Service) DeleteCaseType(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteCaseType(ctx, id); err != nil {
		if errors.Is(err, errs.ErrForeignKeyViolation) {
			return errs.NewConflictErr(fmt.Errorf("case type is referenced by existing cases"))
		}
		return fmt.Errorf("delete case type: %w", err)
	}
	return nil
}
