package rapportService

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (s *Service) DeleteRapport(ctx context.Context, caseID uuid.UUID) error {
	if err := s.repo.DeleteRapport(ctx, caseID); err != nil {
		return fmt.Errorf("failed to delete rapport: %w", err)
	}
	return nil
}
