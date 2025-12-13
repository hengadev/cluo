package document

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// CreateContract creates a new contract.
func (s *Service) CreateContract(ctx context.Context, contract *document.Contract) (*document.Contract, error) {
	// Validate contract
	if err := contract.Validate(); err != nil {
		return nil, errs.NewInvalidValueErr(fmt.Sprintf("contract validation failed: %s", err.Error()))
	}

	// Verify case exists
	caseExists, err := s.caseRepo.ExistsCase(ctx, contract.CaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to check case existence: %w", err)
	}
	if !caseExists {
		return nil, errs.NewNotFoundErr(fmt.Errorf("case not found"), "case")
	}

	// Verify client exists
	clientExists, err := s.clientRepo.ExistsClient(ctx, contract.ClientID)
	if err != nil {
		return nil, fmt.Errorf("failed to check client existence: %w", err)
	}
	if !clientExists {
		return nil, errs.NewNotFoundErr(fmt.Errorf("client not found"), "client")
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
	if err := s.createDocumentVersion(ctx, contract, nil, stringPtr("Initial creation")); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return contract, nil
}

