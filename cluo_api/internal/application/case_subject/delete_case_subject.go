package caseSubjectService

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *Service) DeleteCaseSubject(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteCaseSubject(ctx, id); err != nil {
		return fmt.Errorf("delete case subject: %w", err)
	}
	return nil
}
