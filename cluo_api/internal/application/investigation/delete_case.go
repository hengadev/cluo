package investigationService

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/investigation"
)

func (s *CaseService) DeleteCase(ctx context.Context, request *investigation.DeleteCaseByIDRequest) error {
	err := s.repo.DeleteCase(ctx, request.ID)
	if err != nil {
		return fmt.Errorf("failed to delete case: %w", err)
	}

	return nil
}