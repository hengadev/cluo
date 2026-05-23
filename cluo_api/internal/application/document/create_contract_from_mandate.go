package document

import (
	"context"

	"github.com/google/uuid"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// CreateContractFromMandate creates a contract from an activated mandate.
func (s *Service) CreateContractFromMandate(ctx context.Context, mandateID string, contract *document.Contract) (*document.Contract, error) {
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

	// Check if mandate is active
	if !mandate.IsActive() {
		return nil, errs.NewInvalidValueErr("mandate must be active to create contract")
	}

	// Link contract to mandate
	parsedMandateID, err := uuid.Parse(mandateID)
	if err != nil {
		return nil, errs.NewInvalidValueErr("invalid mandate ID format")
	}
	contract.LinkedMandateID = &parsedMandateID
	contract.CaseID = mandate.CaseID
	contract.ClientID = mandate.ClientID

	// Validate contract
	if err := contract.Validate(); err != nil {
		return nil, errs.NewInvalidValueErr("contract validation failed: " + err.Error())
	}

	// Set initial status
	if contract.Status == "" {
		contract.Status = document.DocumentStatusDraft
	}

	// Encrypt contract for storage
	contractEncx, err := document.ProcessContractEncx(ctx, s.crypto, contract)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("contract", err)
	}

	// Save contract
	if err := s.repo.CreateContract(ctx, contractEncx); err != nil {
		return nil, errs.NewNotCreatedErr(err, "contract")
	}

	// Create initial version
	if err := s.createDocumentVersion(ctx, contract, nil, stringPtr("Created from mandate")); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return contract, nil
}