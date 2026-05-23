package document

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// GetOverdueInvoices retrieves overdue invoices.
func (s *Service) GetOverdueInvoices(ctx context.Context, pagination document.Pagination) ([]*document.Invoice, int, error) {
	if err := pagination.Validate(); err != nil {
		return nil, 0, fmt.Errorf("invalid pagination: %w", err)
	}

	// TODO: Implement overdue invoice filtering in repository
	// For now, we'll get all invoices and filter in service
	invoicesEncx, _, err := s.repo.ListInvoicesByCase(ctx, "", pagination)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list invoices: %w", err)
	}

	var overdueInvoices []*document.Invoice
	for _, invoiceEncx := range invoicesEncx {
		invoice, err := document.DecryptInvoiceEncx(ctx, s.crypto, invoiceEncx)
		if err != nil {
			continue // Skip if decryption fails
		}

		if invoice.IsOverdue() {
			overdueInvoices = append(overdueInvoices, invoice)
		}
	}

	return overdueInvoices, len(overdueInvoices), nil
}