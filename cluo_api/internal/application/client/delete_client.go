package clientService

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/client"
)

func (s *Service) DeleteClient(ctx context.Context, r *client.DeleteClientRequest) error {
	// Call repository with parsed UUID
	if err := s.repo.DeleteClient(ctx, r.ID); err != nil {
		return fmt.Errorf("failed to delete client: %w", err)
	}
	return nil
}
