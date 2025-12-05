package clientService

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/client"
)

func (s *Service) UpdateClient(ctx context.Context, request *client.UpdateClientRequest) (*client.ClientResponse, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Get existing client from repository
	clientEncx, err := s.repo.GetClientByID(ctx, request.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get client by ID: %w", err)
	}

	// Decrypt clientDecrypted data to allow field updates using the new generated function
	clientDecrypted, err := client.DecryptClientEncx(ctx, s.crypto, clientEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("client for update", err)
	}

	// Update only non-nil fields from request
	if request.Name != nil {
		clientDecrypted.Name = *request.Name
	}

	if request.Type != nil {
		clientDecrypted.Type = client.ClientType(*request.Type)
	}

	// process client

	// Encrypt the client data using the new generated function
	updatedClientEncx, err := client.ProcessClientEncx(ctx, s.crypto, clientDecrypted)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("client for update", err)
	}

	// update client

	// Save updated client to repository
	if err := s.repo.UpdateClient(ctx, updatedClientEncx); err != nil {
		return nil, fmt.Errorf("failed to update client: %w", err)
	}

	return clientDecrypted.ToResponse(), nil
}
