package clientService

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
)

func (s *Service) GetContactIDsForClient(ctx context.Context, clientID uuid.UUID) ([]uuid.UUID, error) {
	// Check if client exists in database
	exists, err := s.repo.ExistsClient(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to check client existence: %w", err)
	}

	if !exists {
		return nil, errs.NewRepositoryNotFoundErr(fmt.Errorf("client with ID %s not found", clientID), "client")
	}

	// Get contact IDs from repository
	contactIDs, err := s.repo.GetContactIDsForClient(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get contact IDs for client: %w", err)
	}

	return contactIDs, nil
}

