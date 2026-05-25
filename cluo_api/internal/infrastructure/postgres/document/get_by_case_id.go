package documentRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// GetFirstByCaseAndType fetches the first document of the given type for the
// specified case ID, querying by the plain case_id column.
func (r *Repository) GetFirstByCaseAndType(ctx context.Context, caseID string, docType document.DocumentType) (document.Documentable, error) {
	switch docType {
	case document.DocumentTypeEstimate:
		return r.getFirstEstimateByCaseID(ctx, caseID)
	case document.DocumentTypeMandate:
		return r.getFirstMandateByCaseID(ctx, caseID)
	case document.DocumentTypeContract:
		return r.getFirstContractByCaseID(ctx, caseID)
	case document.DocumentTypeInvoice:
		return r.getFirstInvoiceByCaseID(ctx, caseID)
	default:
		return nil, fmt.Errorf("unsupported document type: %s", docType)
	}
}

func (r *Repository) getFirstEstimateByCaseID(ctx context.Context, caseID string) (*document.EstimateEncx, error) {
	query := `
		SELECT id, status, created_at, updated_at,
			   caseid_encrypted, clientid_encrypted,
			   estimatenumber_encrypted, lineitems_encrypted, estimatedtotal_encrypted, notes_encrypted,
			   issue_date, valid_until, accepted, accepted_at, accepted_by,
			   dek_encrypted, key_version, metadata
		FROM estimates
		WHERE case_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	var e document.EstimateEncx
	err := r.pool.QueryRow(ctx, query, caseID).Scan(
		&e.ID, &e.Status, &e.CreatedAt, &e.UpdatedAt,
		&e.CaseIDEncrypted, &e.ClientIDEncrypted,
		&e.EstimateNumberEncrypted, &e.LineItemsEncrypted, &e.EstimatedTotalEncrypted, &e.NotesEncrypted,
		&e.IssueDate, &e.ValidUntil, &e.Accepted, &e.AcceptedAt, &e.AcceptedBy,
		&e.DEKEncrypted, &e.KeyVersion, &e.Metadata,
	)
	if err != nil {
		return nil, errs.ClassifyPgError("get estimate by case ID", err)
	}
	return &e, nil
}

func (r *Repository) getFirstMandateByCaseID(ctx context.Context, caseID string) (*document.MandateEncx, error) {
	query := `
		SELECT id, status, created_at, updated_at,
			   caseid_encrypted, clientid_encrypted, mandatenumber_encrypted,
			   scopeofwork_encrypted, termsconditions_encrypted,
			   clientsignature_encrypted, investigatorsignature_encrypted, specialinstructions_encrypted,
			   issue_date, valid_from, valid_until, linked_estimate_id, jurisdiction,
			   dek_encrypted, key_version, metadata
		FROM mandates
		WHERE case_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	var m document.MandateEncx
	err := r.pool.QueryRow(ctx, query, caseID).Scan(
		&m.ID, &m.Status, &m.CreatedAt, &m.UpdatedAt,
		&m.CaseIDEncrypted, &m.ClientIDEncrypted, &m.MandateNumberEncrypted,
		&m.ScopeOfWorkEncrypted, &m.TermsConditionsEncrypted,
		&m.ClientSignatureEncrypted, &m.InvestigatorSignatureEncrypted, &m.SpecialInstructionsEncrypted,
		&m.IssueDate, &m.ValidFrom, &m.ValidUntil, &m.LinkedEstimateID, &m.Jurisdiction,
		&m.DEKEncrypted, &m.KeyVersion, &m.Metadata,
	)
	if err != nil {
		return nil, errs.ClassifyPgError("get mandate by case ID", err)
	}
	return &m, nil
}

func (r *Repository) getFirstContractByCaseID(ctx context.Context, caseID string) (*document.ContractEncx, error) {
	query := `
		SELECT id, status, created_at, updated_at,
			   caseid_encrypted, clientid_encrypted, contractnumber_encrypted,
			   scopeofservices_encrypted, paymentterms_encrypted,
			   confidentiality_encrypted, terminationclause_encrypted,
			   signatures_encrypted, contractvalue_encrypted, renewalterms_encrypted,
			   start_date, end_date, linked_mandate_id, currency, governing_law,
			   dek_encrypted, key_version, metadata
		FROM contracts
		WHERE case_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	var c document.ContractEncx
	err := r.pool.QueryRow(ctx, query, caseID).Scan(
		&c.ID, &c.Status, &c.CreatedAt, &c.UpdatedAt,
		&c.CaseIDEncrypted, &c.ClientIDEncrypted, &c.ContractNumberEncrypted,
		&c.ScopeOfServicesEncrypted, &c.PaymentTermsEncrypted,
		&c.ConfidentialityEncrypted, &c.TerminationClauseEncrypted,
		&c.SignaturesEncrypted, &c.ContractValueEncrypted, &c.RenewalTermsEncrypted,
		&c.StartDate, &c.EndDate, &c.LinkedMandateID, &c.Currency, &c.GoverningLaw,
		&c.DEKEncrypted, &c.KeyVersion, &c.Metadata,
	)
	if err != nil {
		return nil, errs.ClassifyPgError("get contract by case ID", err)
	}
	return &c, nil
}

func (r *Repository) getFirstInvoiceByCaseID(ctx context.Context, caseID string) (*document.InvoiceEncx, error) {
	query := `
		SELECT id, status, created_at, updated_at,
			   caseid_encrypted, clientid_encrypted, invoicenumber_encrypted,
			   lineitems_encrypted, totalamount_encrypted, taxamount_encrypted,
			   notes_encrypted, paidamount_encrypted, paymentmethod_encrypted,
			   paymentterms_encrypted, latefee_encrypted,
			   issue_date, due_date, tax_rate, payment_status, paid_at,
			   linked_contract_id, currency, late_fee_rate,
			   dek_encrypted, key_version, metadata
		FROM invoices
		WHERE case_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	var inv document.InvoiceEncx
	err := r.pool.QueryRow(ctx, query, caseID).Scan(
		&inv.ID, &inv.Status, &inv.CreatedAt, &inv.UpdatedAt,
		&inv.CaseIDEncrypted, &inv.ClientIDEncrypted, &inv.InvoiceNumberEncrypted,
		&inv.LineItemsEncrypted, &inv.TotalAmountEncrypted, &inv.TaxAmountEncrypted,
		&inv.NotesEncrypted, &inv.PaidAmountEncrypted, &inv.PaymentMethodEncrypted,
		&inv.PaymentTermsEncrypted, &inv.LateFeeEncrypted,
		&inv.IssueDate, &inv.DueDate, &inv.TaxRate, &inv.PaymentStatus, &inv.PaidAt,
		&inv.LinkedContractID, &inv.Currency, &inv.LateFeeRate,
		&inv.DEKEncrypted, &inv.KeyVersion, &inv.Metadata,
	)
	if err != nil {
		return nil, errs.ClassifyPgError("get invoice by case ID", err)
	}
	return &inv, nil
}
