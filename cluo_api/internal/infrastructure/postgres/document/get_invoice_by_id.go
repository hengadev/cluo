package documentRepository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// GetInvoiceByID retrieves an invoice by its ID.
func (r *Repository) GetInvoiceByID(ctx context.Context, id string) (*document.InvoiceEncx, error) {
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
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	var invoice document.InvoiceEncx

	err := row.Scan(
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
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("invoice not found")
		}
		return nil, fmt.Errorf("failed to get invoice: %w", err)
	}

	return &invoice, nil
}
