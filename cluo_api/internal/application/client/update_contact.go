package clientService

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/client"
)

func (s *Service) UpdateContact(ctx context.Context, request *client.UpdateContactRequest) (*client.ContactResponse, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Get existing contact from repository
	contactEncx, err := s.repo.GetContactByID(ctx, request.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get contact by ID: %w", err)
	}

	// Decrypt contact data to allow field updates using the new generated function
	contact, err := client.DecryptContactEncx(ctx, s.crypto, contactEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("contact for update", err)
	}

	// Update only non-nil fields from request
	if request.Lastname != nil {
		contact.Lastname = *request.Lastname
	}

	if request.Firstname != nil {
		contact.Firstname = *request.Firstname
	}

	if request.Email != nil {
		contact.Email = *request.Email
	}

	if request.Phone != nil {
		contact.Phone = *request.Phone
	}

	if request.Position != nil {
		contact.Position = *request.Position
	}

	// process contact

	// Encrypt the contact data using the new generated function
	updatedcontactEncx, err := client.ProcessContactEncx(ctx, s.crypto, contact)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("contact for update", err)
	}

	// update contact

	// Save updated contact to repository
	if err := s.repo.UpdateContact(ctx, updatedcontactEncx); err != nil {
		return nil, fmt.Errorf("failed to update contact: %w", err)
	}

	return contact.ToResponse(), nil
}
