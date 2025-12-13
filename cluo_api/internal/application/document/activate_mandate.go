package document

import (
	"context"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// ActivateMandate activates a mandate.
func (s *Service) ActivateMandate(ctx context.Context, mandateID string) (*document.Mandate, error) {
	// Get mandate
	mandateEncx, err := s.repo.GetMandateByID(ctx, mandateID)
	if err != nil {
		return nil, errs.NewNotFoundErr(err, "mandate")
	}

	// Decrypt mandate
	mandate, err := document.DecryptMandateEncx(ctx, s.crypto, mandateEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("mandate", err)
	}

	// Check if mandate can be activated
	if mandate.ClientSignature == nil {
		return nil, errs.NewInvalidValueErr("mandate must have client signature to be activated")
	}

	if mandate.IsExpired() {
		return nil, errs.NewInvalidValueErr("cannot activate expired mandate")
	}

	// Activate mandate
	mandate.Activate()

	// Encrypt updated mandate
	updatedMandateEncx, err := document.ProcessMandateEncx(ctx, s.crypto, mandate)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("mandate", err)
	}

	// Save updates
	if err := s.repo.UpdateMandate(ctx, updatedMandateEncx); err != nil {
		return nil, errs.NewNotUpdatedErr(err, "mandate")
	}

	// Create version record
	authorID := s.getUserIDFromContext(ctx)
	if err := s.createDocumentVersion(ctx, mandate, &authorID, stringPtr("Mandate activated")); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return mandate, nil
}

