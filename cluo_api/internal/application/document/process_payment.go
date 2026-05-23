package document

import (
	"context"
	"fmt"
	"time"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// ProcessPayment processes a payment for an invoice.
func (s *Service) ProcessPayment(ctx context.Context, invoiceID string, req *document.PaymentRequest) (*document.Invoice, error) {
	// Validate request
	if err := req.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

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

	// Check if payment can be processed
	if invoice.PaymentStatus == document.PaymentStatusPaid {
		return nil, errs.NewConflictErr(fmt.Errorf("invoice already fully paid"))
	}

	if invoice.PaymentStatus == document.PaymentStatusVoid {
		return nil, errs.NewInvalidValueErr("cannot process payment for voided invoice")
	}

	// Process payment
	invoice.AddPayment(req.Amount, req.PaymentMethod)
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
	if err := s.createDocumentVersion(ctx, invoice, &authorID, stringPtr(fmt.Sprintf("Payment processed: %.2f %s", req.Amount, req.PaymentMethod))); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return invoice, nil
}