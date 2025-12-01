package clientService

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/client"
)

func (s *Service) GetContactByID(ctx context.Context, request *client.GetContactByIDRequest) (*client.ContactResponse, error) {
	// Retrieve contact from repository
	contactEncx, err := s.repo.GetContactByID(ctx, request.ContactID)
	if err != nil {
		return nil, fmt.Errorf("failed to get contact by ID: %w", err)
	}

	// Decrypt contact data
	contact, err := client.DecryptContactEncx(ctx, s.crypto, contactEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("contact", err)
	}

	// Build response
	response := &client.ContactResponse{
		ID:        contact.ID.String(),
		ClientID:  contact.ClientID.String(),
		Lastname:  contact.Lastname,
		Firstname: contact.Firstname,
		Email:     contact.Email,
		Phone:     contact.Phone,
		Position:  contact.Position,
		CreatedAt: contact.CreatedAt,
	}

	return response, nil
}
