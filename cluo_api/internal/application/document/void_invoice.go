package document

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// VoidInvoice voids an invoice.
func (s *Service) VoidInvoice(ctx context.Context, invoiceID string) (*document.Invoice, error) {
	// Get invoice
	invoiceEncx, err := s.repo.GetInvoiceByID(ctx, invoiceID)
	if err != nil {
		return nil, errs.NewNotFoundErr(err, "invoice")
	}

	// Decrypt invoice
	invoice, err := document.DecryptInvoiceEncx(ctx, s.crypto, invoiceEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("invoice", err)
	}

	// State machine enforcement: must be able to transition to cancelled
	if err := s.validateDocumentTransition(invoice, document.DocumentStatusCancelled); err != nil {
		return nil, err
	}

	if invoice.PaymentStatus == document.PaymentStatusPaid {
		return nil, errs.NewConflictErr(fmt.Errorf("cannot void paid invoice"))
	}

	if invoice.PaymentStatus == document.PaymentStatusRefunded {
		return nil, errs.NewConflictErr(fmt.Errorf("cannot void refunded invoice"))
	}

	if invoice.PaymentStatus == document.PaymentStatusVoid {
		return nil, errs.NewConflictErr(fmt.Errorf("invoice already voided"))
	}

	// Void invoice
	if err := invoice.Void(); err != nil {
		return nil, errs.NewConflictErr(err)
	}

	// Encrypt updated invoice
	updatedInvoiceEncx, err := document.ProcessInvoiceEncx(ctx, s.crypto, invoice)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("invoice", err)
	}

	// Save updates
	if err := s.repo.UpdateInvoice(ctx, updatedInvoiceEncx); err != nil {
		return nil, errs.NewNotUpdatedErr(err, "invoice")
	}

	// Create version record
	authorID := s.getUserIDFromContext(ctx)
	if err := s.createDocumentVersion(ctx, invoice, &authorID, stringPtr("Invoice voided")); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return invoice, nil
}