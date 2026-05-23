package documentRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// UpdateInvoice updates an existing invoice in the database.
func (r *Repository) UpdateInvoice(ctx context.Context, invoice *document.InvoiceEncx) error {
	query := `
		UPDATE invoices SET
			status = $2, updated_at = $3,
			caseid_encrypted = $4, clientid_encrypted = $5,
			invoicenumber_encrypted = $6, lineitems_encrypted = $7, totalamount_encrypted = $8,
			taxamount_encrypted = $9, notes_encrypted = $10, paidamount_encrypted = $11,
			paymentmethod_encrypted = $12, paymentterms_encrypted = $13, latefee_encrypted = $14,
			issue_date = $15, due_date = $16, tax_rate = $17, payment_status = $18, paid_at = $19,
			linked_contract_id = $20, currency = $21, late_fee_rate = $22,
			dek_encrypted = $23, key_version = $24, metadata = $25
		WHERE id = $1
	`

	_, err := r.pool.Exec(ctx, query,
		invoice.ID, invoice.Status, invoice.UpdatedAt,
		invoice.CaseIDEncrypted, invoice.ClientIDEncrypted,
		invoice.InvoiceNumberEncrypted, invoice.LineItemsEncrypted, invoice.TotalAmountEncrypted,
		invoice.TaxAmountEncrypted, invoice.NotesEncrypted, invoice.PaidAmountEncrypted,
		invoice.PaymentMethodEncrypted, invoice.PaymentTermsEncrypted, invoice.LateFeeEncrypted,
		invoice.IssueDate, invoice.DueDate, invoice.TaxRate, invoice.PaymentStatus, invoice.PaidAt,
		invoice.LinkedContractID, invoice.Currency, invoice.LateFeeRate,
		invoice.DEKEncrypted, invoice.KeyVersion, invoice.Metadata,
	)

	if err != nil {
		return fmt.Errorf("failed to update invoice: %w", err)
	}

	return nil
}
