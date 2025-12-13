package document

import (
	"context"
	"time"

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

	// Check if invoice can be voided
	if invoice.PaymentStatus == document.PaymentStatusPaid {
		return nil, errs.NewInvalidValueErr("cannot void paid invoice")
	}

	if invoice.PaymentStatus == document.PaymentStatusVoid {
		return nil, errs.NewInvalidValueErr("invoice already voided")
	}

	// Void invoice
	invoice.Void()
	invoice.UpdatedAt = time.Now()

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