package clientService

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/client"
)

func (s *Service) CreateContact(ctx context.Context, request *client.CreateContactRequest) (*client.ContactResponse, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Check if client exists in database
	exists, err := s.repo.ExistsClient(ctx, request.ClientID)
	if err != nil {
		return nil, fmt.Errorf("failed to check client existence: %w", err)
	}

	if !exists {
		return nil, errs.NewRepositoryNotFoundErr(fmt.Errorf("client with ID %s not found", request.ClientID), "client")
	}

	contact := client.NewContact(request)
	contactEncx, err := client.ProcessContactEncx(ctx, s.crypto, contact)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("contact", err)
	}

	if err := s.repo.CreateContact(ctx, contactEncx); err != nil {
		return nil, fmt.Errorf("failed to create contact: %w", err)
	}
	return contact.ToResponse(), nil
}
