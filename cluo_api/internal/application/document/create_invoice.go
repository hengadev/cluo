package document

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// CreateInvoice creates a new invoice.
func (s *Service) CreateInvoice(ctx context.Context, invoice *document.Invoice) (*document.Invoice, error) {
	// Validate invoice
	if err := invoice.Validate(); err != nil {
		return nil, errs.NewInvalidValueErr(fmt.Sprintf("invoice validation failed: %s", err.Error()))
	}

	// Verify case exists
	caseExists, err := s.caseRepo.ExistsCase(ctx, invoice.CaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to check case existence: %w", err)
	}
	if !caseExists {
		return nil, errs.NewNotFoundErr(fmt.Errorf("case not found"), "case")
	}

	// Verify client exists
	clientExists, err := s.clientRepo.ExistsClient(ctx, invoice.ClientID)
	if err != nil {
		return nil, fmt.Errorf("failed to check client existence: %w", err)
	}
	if !clientExists {
		return nil, errs.NewNotFoundErr(fmt.Errorf("client not found"), "client")
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
	if err := s.createDocumentVersion(ctx, invoice, nil, stringPtr("Initial creation")); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return invoice, nil
}