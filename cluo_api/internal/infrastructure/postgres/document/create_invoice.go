package documentRepository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// CreateInvoice creates a new invoice in the database.
func (r *Repository) CreateInvoice(ctx context.Context, invoice *document.Invoice) error {
	lineItemsJSON, err := json.Marshal(invoice.LineItems)
	if err != nil {
		return fmt.Errorf("failed to marshal line items: %w", err)
	}

	query := `
		INSERT INTO invoices (
			id, case_id, client_id, status, invoice_number, issue_date, due_date,
			line_items, total_amount, tax_rate, tax_amount, notes, payment_status,
			paid_at, paid_amount, payment_method, linked_contract_id, currency,
			payment_terms, late_fee, late_fee_rate
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
	`

	_, err = r.pool.Exec(ctx, query,
		invoice.ID, invoice.CaseID, invoice.ClientID, invoice.Status,
		invoice.InvoiceNumber, invoice.IssueDate, invoice.DueDate,
		lineItemsJSON, invoice.TotalAmount, invoice.TaxRate, invoice.TaxAmount,
		invoice.Notes, invoice.PaymentStatus, invoice.PaidAt, invoice.PaidAmount,
		invoice.PaymentMethod, invoice.LinkedContractID, invoice.Currency,
		invoice.PaymentTerms, invoice.LateFee, invoice.LateFeeRate,
	)

	if err != nil {
		return fmt.Errorf("failed to create invoice: %w", err)
	}

	return nil
}
