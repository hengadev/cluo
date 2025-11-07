package clientService

import (
	"context"
	"errors"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/client"

	"github.com/google/uuid"
)

func (s *Service) DeleteContact(ctx context.Context, r *client.DeleteContactRequest) error {
	if err := r.Valid(ctx); err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}

	// Parse contact ID from string to UUID
	contactID, err := parseUUID(r.ContactID)
	if err != nil {
		return errs.NewInvalidValueErr("invalid contact ID format")
	}

	// Call repository with parsed UUID
	if err := s.repo.DeleteContact(ctx, contactID); err != nil {
		switch {
		case errors.Is(err, errs.ErrRepositoryNotFound):
			return errs.NewNotFoundErr(err, "contact")
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return errs.NewExternalServiceErr(err, "database unavailable")
		case errors.Is(err, errs.ErrQueryCancelled), errors.Is(err, context.Canceled):
			return fmt.Errorf("delete contact cancelled: %w", err)
		case errors.Is(err, errs.ErrTransactionFailure), errors.Is(err, errs.ErrDeadlock):
			return errs.NewExternalServiceErr(err, "database transaction failed")
		case errors.Is(err, errs.ErrResourceExhausted):
			return errs.NewExternalServiceErr(err, "database resources exhausted")
		case errors.Is(err, errs.ErrPermissionDenied):
			return errs.NewInternalErr(fmt.Errorf("database permission denied: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return errs.NewInternalErr(fmt.Errorf("database error: %w", err))
		case errors.Is(err, context.DeadlineExceeded):
			return fmt.Errorf("delete contact timeout: %w", err)
		default:
			return errs.NewInternalErr(fmt.Errorf("failed to delete contact: %w", err))
		}
	}
	return nil
}

// parseUUID parses a string into a UUID
func parseUUID(uuidStr string) (uuid.UUID, error) {
	return uuid.Parse(uuidStr)
}
