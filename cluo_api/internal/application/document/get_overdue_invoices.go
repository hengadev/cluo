package document

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// GetOverdueInvoices retrieves a paginated list of invoices where
// payment_status = unpaid and due_date < now.
func (s *Service) GetOverdueInvoices(ctx context.Context, pagination document.Pagination) ([]*document.Invoice, int, error) {
	if err := pagination.Validate(); err != nil {
		return nil, 0, fmt.Errorf("invalid pagination: %w", err)
	}

	invoicesEncx, total, err := s.repo.ListOverdueInvoices(ctx, pagination)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list overdue invoices: %w", err)
	}

	overdueInvoices := make([]*document.Invoice, 0, len(invoicesEncx))
	for _, invoiceEncx := range invoicesEncx {
		invoice, err := document.DecryptInvoiceEncx(ctx, s.crypto, invoiceEncx)
		if err != nil {
			continue
		}
		overdueInvoices = append(overdueInvoices, invoice)
	}

	return overdueInvoices, total, nil
}
