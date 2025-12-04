package clientService

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/client"
)

func (s *Service) CreateClient(ctx context.Context, request *client.CreateClientRequest) (*client.ClientResponse, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	newClient := client.NewClient(request)
	clientEncx, err := client.ProcessClientEncx(ctx, s.crypto, newClient)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("client", err)
	}

	if err := s.repo.CreateClient(ctx, clientEncx); err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return newClient.ToResponse(), nil
}
