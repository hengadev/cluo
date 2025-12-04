package clientService

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/client"
)

func (s *Service) GetClientByID(ctx context.Context, request *client.GetClientByIDRequest) (*client.ClientResponse, error) {
	// Retrieve client from repository
	clientEncx, err := s.repo.GetClientByID(ctx, request.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get client by ID: %w", err)
	}

	// Decrypt c data
	c, err := client.DecryptClientEncx(ctx, s.crypto, clientEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("client", err)
	}

	// Build response
	response := &client.ClientResponse{
		ID:         c.ID.String(),
		Name:       c.Name,
		Type:       string(c.Type),
		ContactIDs: c.ContactIDs,
	}

	return response, nil
}
