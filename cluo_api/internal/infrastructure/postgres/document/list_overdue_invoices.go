package documentRepository

import (
	"context"
	"fmt"
	"time"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// ListOverdueInvoices retrieves all invoices where payment_status = 'unpaid' and
// due_date < now, ordered by due_date ascending (most overdue first).
func (r *Repository) ListOverdueInvoices(ctx context.Context, pagination document.Pagination) ([]*document.InvoiceEncx, int, error) {
	now := time.Now()

	// Get total count
	var total int
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM invoices
		WHERE payment_status = 'unpaid' AND due_date < $1
	`, now).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count overdue invoices: %w", err)
	}

	if total == 0 {
		return []*document.InvoiceEncx{}, 0, nil
	}

	// Get paginated results
	offset := pagination.GetOffset()
	limit := pagination.PageSize
	query := `
		SELECT id, status, created_at, updated_at,
			   caseid_encrypted, clientid_encrypted,
			   invoicenumber_encrypted, lineitems_encrypted, totalamount_encrypted,
			   taxamount_encrypted, notes_encrypted, paidamount_encrypted,
			   paymentmethod_encrypted, paymentterms_encrypted, latefee_encrypted,
			   issue_date, due_date, tax_rate, payment_status, paid_at,
			   linked_contract_id, currency, late_fee_rate,
			   dek_encrypted, key_version, metadata
		FROM invoices
		WHERE payment_status = 'unpaid' AND due_date < $1
		ORDER BY due_date ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, now, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query overdue invoices: %w", err)
	}
	defer rows.Close()

	var invoices []*document.InvoiceEncx
	for rows.Next() {
		var invoice document.InvoiceEncx

		err := rows.Scan(
			&invoice.ID, &invoice.Status, &invoice.CreatedAt, &invoice.UpdatedAt,
			&invoice.CaseIDEncrypted, &invoice.ClientIDEncrypted,
			&invoice.InvoiceNumberEncrypted, &invoice.LineItemsEncrypted, &invoice.TotalAmountEncrypted,
			&invoice.TaxAmountEncrypted, &invoice.NotesEncrypted, &invoice.PaidAmountEncrypted,
			&invoice.PaymentMethodEncrypted, &invoice.PaymentTermsEncrypted, &invoice.LateFeeEncrypted,
			&invoice.IssueDate, &invoice.DueDate, &invoice.TaxRate, &invoice.PaymentStatus, &invoice.PaidAt,
			&invoice.LinkedContractID, &invoice.Currency, &invoice.LateFeeRate,
			&invoice.DEKEncrypted, &invoice.KeyVersion, &invoice.Metadata,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan overdue invoice: %w", err)
		}

		invoices = append(invoices, &invoice)
	}

	return invoices, total, nil
}
