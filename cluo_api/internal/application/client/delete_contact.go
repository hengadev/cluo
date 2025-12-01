package clientService

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/client"
)

func (s *Service) DeleteContact(ctx context.Context, r *client.DeleteContactRequest) error {
	// Call repository with parsed UUID
	if err := s.repo.DeleteContact(ctx, r.ContactID); err != nil {
		return fmt.Errorf("failed to delete contact: %w", err)
	}
	return nil
}
