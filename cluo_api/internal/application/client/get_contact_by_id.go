package clientService

import (
	"context"
	"errors"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/client"

	"github.com/google/uuid"
)

func (s *Service) GetContactByID(ctx context.Context, request *client.GetContactByIDRequest) (*client.ContactResponse, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	contactID, err := uuid.Parse(request.ContactID)
	if err != nil {
		return nil, errs.NewInvalidValueErr("invalid contact ID format")
	}

	// Retrieve contact from repository
	contactEncx, err := s.repo.GetContactByID(ctx, contactID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrDomainNotFound):
			return nil, errs.NewNotFoundErr(err, "contact")
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return nil, errs.NewExternalServiceErr(err, "database unavailable")
		case errors.Is(err, errs.ErrQueryCancelled):
			return nil, fmt.Errorf("get contact cancelled: %w", err)
		case errors.Is(err, errs.ErrTransactionFailure), errors.Is(err, errs.ErrDeadlock):
			return nil, errs.NewExternalServiceErr(err, "database transaction failed")
		case errors.Is(err, errs.ErrResourceExhausted):
			return nil, errs.NewExternalServiceErr(err, "database resources exhausted")
		case errors.Is(err, errs.ErrPermissionDenied):
			return nil, errs.NewInternalErr(fmt.Errorf("database permission denied: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return nil, errs.NewInternalErr(fmt.Errorf("database error: %w", err))
		case errors.Is(err, context.Canceled):
			return nil, fmt.Errorf("get contact cancelled: %w", err)
		case errors.Is(err, context.DeadlineExceeded):
			return nil, fmt.Errorf("get contact timeout: %w", err)
		default:
			return nil, errs.NewInternalErr(fmt.Errorf("failed to get contact: %w", err))
		}
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

