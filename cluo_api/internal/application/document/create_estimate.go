package document

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// CreateEstimate creates a new estimate.
func (s *Service) CreateEstimate(ctx context.Context, estimate *document.Estimate) (*document.Estimate, error) {
	// Validate estimate
	if err := estimate.Validate(); err != nil {
		return nil, errs.NewInvalidValueErr(fmt.Sprintf("estimate validation failed: %s", err.Error()))
	}

	// Verify case exists
	caseExists, err := s.caseRepo.ExistsCase(ctx, estimate.CaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to check case existence: %w", err)
	}
	if !caseExists {
		return nil, errs.NewNotFoundErr(fmt.Errorf("case not found"), "case")
	}

	// Verify client exists
	clientExists, err := s.clientRepo.ExistsClient(ctx, estimate.ClientID)
	if err != nil {
		return nil, fmt.Errorf("failed to check client existence: %w", err)
	}
	if !clientExists {
		return nil, errs.NewNotFoundErr(fmt.Errorf("client not found"), "client")
	}

	// Set initial status
	if estimate.Status == "" {
		estimate.Status = document.DocumentStatusDraft
	}

	// Encrypt estimate for storage
	estimateEncx, err := document.ProcessEstimateEncx(ctx, s.crypto, estimate)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("estimate", err)
	}

	// Save estimate
	if err := s.repo.CreateEstimate(ctx, estimateEncx); err != nil {
		return nil, errs.NewNotCreatedErr(err, "estimate")
	}

	// Create initial version
	if err := s.createDocumentVersion(ctx, estimate, nil, stringPtr("Initial creation")); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return estimate, nil
}