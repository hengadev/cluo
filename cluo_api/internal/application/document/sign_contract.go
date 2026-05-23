package document

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// SignContract signs a contract.
func (s *Service) SignContract(ctx context.Context, contractID string, req *document.SignDocumentRequest) (*document.Contract, error) {
	// Validate request
	if err := req.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

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

	// Create signature
	signature := document.Signature{
		ID:        uuid.New(),
		Name:      req.SignerName,
		Role:      req.SignerRole,
		Method:    req.Method,
		IPAddress: req.IPAddress,
		UserAgent: req.UserAgent,
		SignedAt:  time.Now(),
	}

	// Add signature
	if err := contract.AddSignature(signature); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

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
	if err := s.createDocumentVersion(ctx, contract, nil, stringPtr("Contract signed")); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return contract, nil
}