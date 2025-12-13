package document

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// AcceptEstimate accepts an estimate and creates a mandate.
func (s *Service) AcceptEstimate(ctx context.Context, estimateID string, acceptedByStr string) (*document.Mandate, error) {
	// Get estimate
	estimateEncx, err := s.repo.GetEstimateByID(ctx, estimateID)
	if err != nil {
		return nil, errs.NewNotFoundErr(err, "estimate")
	}

	// Decrypt estimate
	estimate, err := document.DecryptEstimateEncx(ctx, s.crypto, estimateEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("estimate", err)
	}

	// Check if estimate can be accepted
	if estimate.Accepted {
		return nil, errs.NewConflictErr(fmt.Errorf("estimate already accepted"))
	}

	if estimate.IsExpired() {
		return nil, errs.NewInvalidValueErr("cannot accept expired estimate")
	}

	acceptedBy, err := uuid.Parse(acceptedByStr)
	if err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Accept estimate
	estimate.Accept(acceptedBy)

	// Encrypt updated estimate
	updatedEstimateEncx, err := document.ProcessEstimateEncx(ctx, s.crypto, estimate)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("estimate", err)
	}

	// Update estimate
	if err := s.repo.UpdateEstimate(ctx, updatedEstimateEncx); err != nil {
		return nil, errs.NewNotUpdatedErr(err, "estimate")
	}

	// Create mandate from estimate
	mandate := document.NewMandateFromEstimate(estimate)
	mandate.Status = document.DocumentStatusDraft

	// Encrypt mandate for storage
	mandateEncx, err := document.ProcessMandateEncx(ctx, s.crypto, mandate)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("mandate", err)
	}

	// Save mandate
	if err := s.repo.CreateMandate(ctx, mandateEncx); err != nil {
		return nil, errs.NewNotCreatedErr(err, "mandate")
	}

	// Create version records
	authorID := s.getUserIDFromContext(ctx)
	if err := s.createDocumentVersion(ctx, estimate, &authorID, stringPtr("Estimate accepted")); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	if err := s.createDocumentVersion(ctx, mandate, &authorID, stringPtr("Created from accepted estimate")); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return mandate, nil
}

