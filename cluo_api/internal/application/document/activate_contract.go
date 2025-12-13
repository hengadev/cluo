package document

import (
	"context"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// ActivateContract activates a contract.
func (s *Service) ActivateContract(ctx context.Context, contractID string) (*document.Contract, error) {
	// Get contract
	contractEncx, err := s.repo.GetContractByID(ctx, contractID)
	if err != nil {
		return nil, errs.NewNotFoundErr(err, "contract")
	}

	// Decrypt contract
	contract, err := document.DecryptContractEncx(ctx, s.crypto, contractEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("contract", err)
	}

	// Check if contract can be activated
	if len(contract.Signatures) == 0 {
		return nil, errs.NewInvalidValueErr("contract must have at least one signature to be activated")
	}

	if contract.IsExpired() {
		return nil, errs.NewInvalidValueErr("cannot activate expired contract")
	}

	// Activate contract
	contract.Activate()

	// Encrypt updated contract
	updatedContractEncx, err := document.ProcessContractEncx(ctx, s.crypto, contract)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("contract", err)
	}

	// Save updates
	if err := s.repo.UpdateContract(ctx, updatedContractEncx); err != nil {
		return nil, errs.NewNotUpdatedErr(err, "contract")
	}

	// Create version record
	authorID := s.getUserIDFromContext(ctx)
	if err := s.createDocumentVersion(ctx, contract, &authorID, stringPtr("Contract activated")); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return contract, nil
}