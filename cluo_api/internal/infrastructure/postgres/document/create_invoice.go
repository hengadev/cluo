package documentRepository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// CreateInvoice creates a new invoice in the database.
func (r *Repository) CreateInvoice(ctx context.Context, invoice *document.InvoiceEncx) error {
	query := `
		INSERT INTO invoices (
			id, status, created_at, updated_at,
			caseid_encrypted, clientid_encrypted,
			invoicenumber_encrypted, lineitems_encrypted, totalamount_encrypted,
			taxamount_encrypted, notes_encrypted, paidamount_encrypted,
			paymentmethod_encrypted, paymentterms_encrypted, latefee_encrypted,
			issue_date, due_date, tax_rate, payment_status, paid_at,
			linked_contract_id, currency, late_fee_rate,
			dek_encrypted, key_version, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26)
	`

	_, err := r.pool.Exec(ctx, query,
		invoice.ID, invoice.Status, invoice.CreatedAt, invoice.UpdatedAt,
		invoice.CaseIDEncrypted, invoice.ClientIDEncrypted,
		invoice.InvoiceNumberEncrypted, invoice.LineItemsEncrypted, invoice.TotalAmountEncrypted,
		invoice.TaxAmountEncrypted, invoice.NotesEncrypted, invoice.PaidAmountEncrypted,
		invoice.PaymentMethodEncrypted, invoice.PaymentTermsEncrypted, invoice.LateFeeEncrypted,
		invoice.IssueDate, invoice.DueDate, invoice.TaxRate, invoice.PaymentStatus, invoice.PaidAt,
		invoice.LinkedContractID, invoice.Currency, invoice.LateFeeRate,
		invoice.DEKEncrypted, invoice.KeyVersion, invoice.Metadata,
	)

	if err != nil {
		return fmt.Errorf("failed to create invoice: %w", err)
	}

	return nil
}
