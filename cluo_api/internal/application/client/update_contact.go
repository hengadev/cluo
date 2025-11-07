package clientService

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/client"
)

func (s *Service) UpdateContact(ctx context.Context, request *client.UpdateContactRequest) error {
	if err := request.Valid(ctx); err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}

	contactID, err := uuid.Parse(request.ID)
	if err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}

	// Get existing contact from repository
	contactEncx, err := s.repo.GetContactByID(ctx, contactID)
	if err != nil {
		switch {

		}
	}

	// Decrypt contact data to allow field updates using the new generated function
	contact, err := client.DecryptContactEncx(ctx, s.crypto, contactEncx)
	if err != nil {
		return errs.NewNotDecryptedErr("contact for update", err)
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
		return errs.NewNotEncryptedErr("contact for update", err)
	}

	// update contact

	// Save updated contact to repository
	if err := s.repo.UpdateContact(ctx, updatedcontactEncx); err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "contact")
		case errors.Is(err, errs.ErrRepositoryNotUpdated):
			return errs.NewNotUpdatedErr(err, "contact")
		case errors.Is(err, errs.ErrUniqueViolation):
			return errs.NewAlreadyExistsError(err, "contact with this email or phone")
		case errors.Is(err, errs.ErrForeignKeyViolation):
			return errs.NewInvalidValueErr("invalid reference in contact data")
		case errors.Is(err, errs.ErrNotNullViolation):
			return errs.NewInvalidValueErr("required field is missing")
		case errors.Is(err, errs.ErrCheckViolation):
			return errs.NewInvalidValueErr("contact data violates database constraints")
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return errs.NewExternalServiceErr(err, "database unavailable")
		case errors.Is(err, errs.ErrQueryCancelled):
			return fmt.Errorf("update contact cancelled: %w", err)
		case errors.Is(err, errs.ErrTransactionFailure), errors.Is(err, errs.ErrDeadlock):
			return errs.NewExternalServiceErr(err, "database transaction failed")
		case errors.Is(err, errs.ErrResourceExhausted):
			return errs.NewExternalServiceErr(err, "database resources exhausted")
		case errors.Is(err, errs.ErrPermissionDenied):
			return errs.NewInternalErr(fmt.Errorf("database permission denied: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return errs.NewInternalErr(fmt.Errorf("database error: %w", err))
		case errors.Is(err, context.Canceled):
			return fmt.Errorf("update contact cancelled: %w", err)
		case errors.Is(err, context.DeadlineExceeded):
			return fmt.Errorf("update contact timeout: %w", err)
		default:
			return errs.NewInternalErr(fmt.Errorf("failed to update contact: %w", err))
		}
	}

	return nil
}
