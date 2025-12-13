package document

import (
	"context"
	"fmt"
	"time"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// UpdateEstimate updates the line items of an estimate.
func (s *Service) UpdateEstimate(ctx context.Context, estimateID string, lineItems []document.EstimateItem) (*document.Estimate, error) {
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

	// Check if estimate can be modified
	if estimate.Accepted {
		return nil, errs.NewConflictErr(fmt.Errorf("cannot modify accepted estimate"))
	}

	// Update line items
	estimate.LineItems = lineItems
	estimate.CalculateTotal()
	estimate.UpdatedAt = time.Now()

	// Validate updated estimate
	if err := estimate.Validate(); err != nil {
		return nil, errs.NewInvalidValueErr(fmt.Sprintf("estimate validation failed: %s", err.Error()))
	}

	// Encrypt updated estimate
	updatedEstimateEncx, err := document.ProcessEstimateEncx(ctx, s.crypto, estimate)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("estimate", err)
	}

	// Save updates
	if err := s.repo.UpdateEstimate(ctx, updatedEstimateEncx); err != nil {
		return nil, errs.NewNotUpdatedErr(err, "estimate")
	}

	// Create version record
	authorID := s.getUserIDFromContext(ctx)
	if err := s.createDocumentVersion(ctx, estimate, &authorID, stringPtr("Updated line items")); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return estimate, nil
}