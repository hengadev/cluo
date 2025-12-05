package clientService

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/client"
)

func (s *Service) GetAllClients(ctx context.Context) ([]*client.ClientResponse, error) {
	clientsEncx, err := s.repo.GetAllClients(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all clients: %w", err)
	}

	// Decrypt clients and build response
	clientResponses := make([]*client.ClientResponse, 0, len(clientsEncx))
	for _, clientEncx := range clientsEncx {
		// Decrypt clientDecrypted
		clientDecrypted, err := client.DecryptClientEncx(ctx, s.crypto, clientEncx)
		if err != nil {
			return nil, errs.NewNotDecryptedErr("client", err)
		}

		// Use domain model's ToResponse method
		clientResponse := clientDecrypted.ToResponse()
		clientResponses = append(clientResponses, clientResponse)
	}
	return clientResponses, nil
}
