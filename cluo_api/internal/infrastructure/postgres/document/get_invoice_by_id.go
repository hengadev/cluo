package documentRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

	query := `
		SELECT id, case_id, client_id, status, invoice_number, issue_date, due_date,
			   line_items, total_amount, tax_rate, tax_amount, notes, payment_status,
			   paid_at, paid_amount, payment_method, linked_contract_id, currency,
			   payment_terms, late_fee, late_fee_rate, created_at, updated_at
		FROM invoices
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	var invoice document.Invoice
	var lineItemsJSON []byte

	err := row.Scan(
		&invoice.ID, &invoice.CaseID, &invoice.ClientID, &invoice.Status,
		&invoice.InvoiceNumber, &invoice.IssueDate, &invoice.DueDate,
		&lineItemsJSON, &invoice.TotalAmount, &invoice.TaxRate, &invoice.TaxAmount,
		&invoice.Notes, &invoice.PaymentStatus, &invoice.PaidAt, &invoice.PaidAmount,
		&invoice.PaymentMethod, &invoice.LinkedContractID, &invoice.Currency,
		&invoice.PaymentTerms, &invoice.LateFee, &invoice.LateFeeRate,
		&invoice.CreatedAt, &invoice.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invoice not found")
		}
		return nil, fmt.Errorf("failed to get invoice: %w", err)
	}

	if err := json.Unmarshal(lineItemsJSON, &invoice.LineItems); err != nil {
		return nil, fmt.Errorf("failed to unmarshal line items: %w", err)
	}

	return &invoice, nil
}
