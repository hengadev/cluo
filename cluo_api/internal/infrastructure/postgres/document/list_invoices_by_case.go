package documentRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

	// Get total count
	var total int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM invoices WHERE case_id = $1", caseID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count invoices: %w", err)
	}

	// Get paginated results
	offset := pagination.GetOffset()
	limit := pagination.PageSize
	query := `
		SELECT id, case_id, client_id, status, invoice_number, issue_date, due_date,
			   line_items, total_amount, tax_rate, tax_amount, notes, payment_status,
			   paid_at, paid_amount, payment_method, linked_contract_id, currency,
			   payment_terms, late_fee, late_fee_rate, created_at, updated_at
		FROM invoices
		WHERE case_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, caseID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query invoices: %w", err)
	}
	defer rows.Close()

	var invoices []*document.Invoice
	for rows.Next() {
		var invoice document.Invoice
		var lineItemsJSON []byte

		err := rows.Scan(
			&invoice.ID, &invoice.CaseID, &invoice.ClientID, &invoice.Status,
			&invoice.InvoiceNumber, &invoice.IssueDate, &invoice.DueDate,
			&lineItemsJSON, &invoice.TotalAmount, &invoice.TaxRate, &invoice.TaxAmount,
			&invoice.Notes, &invoice.PaymentStatus, &invoice.PaidAt, &invoice.PaidAmount,
			&invoice.PaymentMethod, &invoice.LinkedContractID, &invoice.Currency,
			&invoice.PaymentTerms, &invoice.LateFee, &invoice.LateFeeRate,
			&invoice.CreatedAt, &invoice.UpdatedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan invoice: %w", err)
		}

		if err := json.Unmarshal(lineItemsJSON, &invoice.LineItems); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal line items: %w", err)
		}

		invoices = append(invoices, &invoice)
	}

	return invoices, total, nil
}

