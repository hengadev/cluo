package document

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// CreateMandate creates a new mandate.
func (s *Service) CreateMandate(ctx context.Context, mandate *document.Mandate) (*document.Mandate, error) {
	// Validate mandate
	if err := mandate.Validate(); err != nil {
		return nil, errs.NewInvalidValueErr(fmt.Sprintf("mandate validation failed: %s", err.Error()))
	}

	// Verify case exists
	caseExists, err := s.caseRepo.ExistsCase(ctx, mandate.CaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to check case existence: %w", err)
	}
	if !caseExists {
		return nil, errs.NewNotFoundErr(fmt.Errorf("case not found"), "case")
	}

	// Verify client exists
	clientExists, err := s.clientRepo.ExistsClient(ctx, mandate.ClientID)
	if err != nil {
		return nil, fmt.Errorf("failed to check client existence: %w", err)
	}
	if !clientExists {
		return nil, errs.NewNotFoundErr(fmt.Errorf("client not found"), "client")
	}

	// Set initial status
	if mandate.Status == "" {
		mandate.Status = document.DocumentStatusDraft
	}

	// Encrypt mandate for storage
	mandateEncx, err := document.ProcessMandateEncx(ctx, s.crypto, mandate)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("mandate", err)
	}

	// Save mandate
	if err := s.repo.CreateMandate(ctx, mandateEncx); err != nil {
		return nil, errs.NewNotCreatedErr(err, "mandate")
	}

	// Create initial version
	if err := s.createDocumentVersion(ctx, mandate, nil, stringPtr("Initial creation")); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return mandate, nil
}