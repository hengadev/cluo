package clientService

import (
	"context"
	"errors"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/client"

	"github.com/google/uuid"
	"github.com/hengadev/encx"
)

func (s *Service) GetAllContactsByClientID(ctx context.Context, clientIDStr string) ([]*client.ContactResponse, error) {
	clientID, err := uuid.Parse(clientIDStr)
	if err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Convert client ID to hash for repository query
	clientIDBytes, err := encx.SerializeValue(clientID)
	if err != nil {
		return nil, errs.NewInvalidValueErr("failed to serialize client ID")
	}
	clientIDHash := s.crypto.HashBasic(ctx, clientIDBytes)

	contactsEncx, err := s.repo.GetAllContactsByClientID(ctx, clientIDHash)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return nil, errs.NewExternalServiceErr(err, "database unavailable")
		case errors.Is(err, errs.ErrQueryCancelled), errors.Is(err, context.Canceled):
			return fmt.Errorf("get all contacts by client ID cancelled: %w", err)
		case errors.Is(err, errs.ErrTransactionFailure), errors.Is(err, errs.ErrDeadlock):
			return nil, errs.NewExternalServiceErr(err, "database transaction failed")
		case errors.Is(err, errs.ErrResourceExhausted):
			return nil, errs.NewExternalServiceErr(err, "database resources exhausted")
		case errors.Is(err, errs.ErrPermissionDenied):
			return nil, errs.NewInternalErr(fmt.Errorf("database permission denied: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return nil, errs.NewInternalErr(fmt.Errorf("database error: %w", err))
		case errors.Is(err, context.DeadlineExceeded):
			return fmt.Errorf("get all contacts by client ID timeout: %w", err)
		default:
			return nil, errs.NewInternalErr(fmt.Errorf("failed to get all contacts by client ID: %w", err))
		}
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
