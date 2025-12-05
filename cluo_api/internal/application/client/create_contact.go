package clientService

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/client"
)

func (s *Service) CreateContact(ctx context.Context, request *client.CreateContactRequest) error {
	if err := request.Valid(ctx); err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}

	// Check if client exists in database
	exists, err := s.repo.ExistsClient(ctx, request.ClientID)
	if err != nil {
		return fmt.Errorf("failed to check client existence: %w", err)
	}

	if !exists {
		return errs.NewNotFoundErr(nil, "client")
	}

	contact := client.NewContact(request)
	contactEncx, err := client.ProcessContactEncx(ctx, s.crypto, contact)
	if err != nil {
		return errs.NewNotEncryptedErr("contact", err)
	}

	if err := s.repo.CreateContact(ctx, contactEncx); err != nil {
		return fmt.Errorf("failed to create contact: %w", err)
	}
	return nil
}
