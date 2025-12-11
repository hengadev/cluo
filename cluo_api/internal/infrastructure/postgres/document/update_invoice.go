package documentRepository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// UpdateInvoice updates an existing invoice in the database.
func (r *Repository) UpdateInvoice(ctx context.Context, invoice *document.Invoice) error {
	lineItemsJSON, err := json.Marshal(invoice.LineItems)
	if err != nil {
		return fmt.Errorf("failed to marshal line items: %w", err)
	}

	query := `
		UPDATE invoices SET
			case_id = $2, client_id = $3, status = $4, invoice_number = $5,
			issue_date = $6, due_date = $7, line_items = $8, total_amount = $9,
			tax_rate = $10, tax_amount = $11, notes = $12, payment_status = $13,
			paid_at = $14, paid_amount = $15, payment_method = $16, linked_contract_id = $17,
			currency = $18, payment_terms = $19, late_fee = $20, late_fee_rate = $21,
			updated_at = NOW()
		WHERE id = $1
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
		return fmt.Errorf("failed to update invoice: %w", err)
	}

	return nil
}
