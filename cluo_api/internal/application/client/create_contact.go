package clientService

import (
	"context"
	"errors"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/client"

	"github.com/hengadev/encx"
)

func (s *Service) CreateContact(ctx context.Context, request *client.CreateContactRequest) error {
	if err := request.Valid(ctx); err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}

	// check if client ID exists and is valid
	clientIDBytes, err := encx.SerializeValue(request.ClientID)
	if err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}
	clientIDHashBasic := s.crypto.HashBasic(ctx, clientIDBytes)

	// Check if client exists in database
	exists, err := s.repo.ExistsByClientID(ctx, clientIDHashBasic)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return errs.NewExternalServiceErr(err, "database unavailable")
		case errors.Is(err, errs.ErrQueryCancelled):
			return fmt.Errorf("check client existence cancelled: %w", err)
		case errors.Is(err, errs.ErrTransactionFailure), errors.Is(err, errs.ErrDeadlock):
			return errs.NewExternalServiceErr(err, "database transaction failed")
		case errors.Is(err, errs.ErrResourceExhausted):
			return errs.NewExternalServiceErr(err, "database resources exhausted")
		case errors.Is(err, errs.ErrPermissionDenied):
			return errs.NewInternalErr(fmt.Errorf("database permission denied: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return errs.NewInternalErr(fmt.Errorf("database error: %w", err))
		case errors.Is(err, errs.ErrInvalidInput):
			return errs.NewInvalidValueErr("invalid client ID format")
		case errors.Is(err, context.Canceled):
			return fmt.Errorf("check client existence cancelled: %w", err)
		case errors.Is(err, context.DeadlineExceeded):
			return fmt.Errorf("check client existence timeout: %w", err)
		default:
			return errs.NewInternalErr(fmt.Errorf("failed to check client existence: %w", err))
		}
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
		switch {
		case errors.Is(err, errs.ErrUniqueViolation):
			return errs.NewAlreadyExistsError(err, "contact with this email or phone")
		case errors.Is(err, errs.ErrForeignKeyViolation):
			return errs.NewInvalidValueErr("invalid client reference")
		case errors.Is(err, errs.ErrNotNullViolation):
			return errs.NewInvalidValueErr("required contact field is missing")
		case errors.Is(err, errs.ErrCheckViolation):
			return errs.NewInvalidValueErr("contact data violates database constraints")
		case errors.Is(err, errs.ErrConnectionFailure), errors.Is(err, errs.ErrTooManyConnections):
			return errs.NewExternalServiceErr(err, "database unavailable")
		case errors.Is(err, errs.ErrQueryCancelled):
			return fmt.Errorf("create contact cancelled: %w", err)
		case errors.Is(err, errs.ErrTransactionFailure), errors.Is(err, errs.ErrDeadlock):
			return errs.NewExternalServiceErr(err, "database transaction failed")
		case errors.Is(err, errs.ErrResourceExhausted):
			return errs.NewExternalServiceErr(err, "database resources exhausted")
		case errors.Is(err, errs.ErrPermissionDenied):
			return errs.NewInternalErr(fmt.Errorf("database permission denied: %w", err))
		case errors.Is(err, errs.ErrDatabase):
			return errs.NewInternalErr(fmt.Errorf("database error: %w", err))
		case errors.Is(err, context.Canceled):
			return fmt.Errorf("create contact cancelled: %w", err)
		case errors.Is(err, context.DeadlineExceeded):
			return fmt.Errorf("create contact timeout: %w", err)
		default:
			return errs.NewInternalErr(fmt.Errorf("failed to create contact: %w", err))
		}
	}
	return nil
}
