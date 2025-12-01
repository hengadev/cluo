package clientService

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/client"

	"github.com/google/uuid"
	"github.com/hengadev/encx"
)

func (s *Service) GetAllContactsByClientID(ctx context.Context, clientID uuid.UUID) ([]*client.ContactResponse, error) {
	// Convert client ID to hash for repository query
	clientIDBytes, err := encx.SerializeValue(clientID)
	if err != nil {
		return nil, errs.NewInvalidValueErr("failed to serialize client ID")
	}
	clientIDHash := s.crypto.HashBasic(ctx, clientIDBytes)

	contactsEncx, err := s.repo.GetAllContactsByClientID(ctx, clientIDHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get all contacts by client ID: %w", err)
	}

	// Decrypt contacts and build response
	contactResponses := make([]*client.ContactResponse, 0, len(contactsEncx))
	for _, contactEncx := range contactsEncx {
		// Decrypt contact
		contact, err := client.DecryptContactEncx(ctx, s.crypto, contactEncx)
		if err != nil {
			return nil, errs.NewNotDecryptedErr("contact", err)
		}

		// Build contact response
		contactResponse := &client.ContactResponse{
			ID:        contact.ID.String(),
			ClientID:  contact.ClientID.String(),
			Lastname:  contact.Lastname,
			Firstname: contact.Firstname,
			Email:     contact.Email,
			Phone:     contact.Phone,
			Position:  contact.Position,
			CreatedAt: contact.CreatedAt,
		}

		contactResponses = append(contactResponses, contactResponse)
	}
	return contactResponses, nil
}
