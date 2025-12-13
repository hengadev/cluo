package document

import (
	"context"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// CreateInvoiceFromContract creates an invoice from an active contract.
func (s *Service) CreateInvoiceFromContract(ctx context.Context, contractID string, invoice *document.Invoice) (*document.Invoice, error) {
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

	// Check if contract is active
	if !contract.IsActive() {
		return nil, errs.NewInvalidValueErr("contract must be active to create invoice")
	}

	// Link invoice to contract
	invoice.LinkedContractID = &contractID
	invoice.CaseID = contract.CaseID
	invoice.ClientID = contract.ClientID

	// Validate invoice
	if err := invoice.Validate(); err != nil {
		return nil, errs.NewInvalidValueErr("invoice validation failed: " + err.Error())
	}

	// Set initial status
	if invoice.Status == "" {
		invoice.Status = document.DocumentStatusDraft
	}

	// Encrypt invoice for storage
	invoiceEncx, err := document.ProcessInvoiceEncx(ctx, s.crypto, invoice)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("invoice", err)
	}

	// Save invoice
	if err := s.repo.CreateInvoice(ctx, invoiceEncx); err != nil {
		return nil, errs.NewNotCreatedErr(err, "invoice")
	}

	// Create initial version
	if err := s.createDocumentVersion(ctx, invoice, nil, stringPtr("Created from contract")); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return invoice, nil
}