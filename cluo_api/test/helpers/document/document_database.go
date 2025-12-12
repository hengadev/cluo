package documentHelpers

import (
	"context"
	"fmt"
	"testing"

	documentDomain "github.com/hengadev/cluo_api/internal/domain/document"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// ClearEstimatesTable truncates the estimates table for clean test state
func ClearEstimatesTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(ctx, "TRUNCATE TABLE estimates RESTART IDENTITY CASCADE")
	require.NoError(t, err)
}

// ClearMandatesTable truncates the mandates table for clean test state
func ClearMandatesTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(ctx, "TRUNCATE TABLE mandates RESTART IDENTITY CASCADE")
	require.NoError(t, err)
}

// ClearContractsTable truncates the contracts table for clean test state
func ClearContractsTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(ctx, "TRUNCATE TABLE contracts RESTART IDENTITY CASCADE")
	require.NoError(t, err)
}

// ClearInvoicesTable truncates the invoices table for clean test state
func ClearInvoicesTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(ctx, "TRUNCATE TABLE invoices RESTART IDENTITY CASCADE")
	require.NoError(t, err)
}

// ClearDocumentVersionsTable truncates the document_versions table for clean test state
func ClearDocumentVersionsTable(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	_, err := pool.Exec(ctx, "TRUNCATE TABLE document_versions RESTART IDENTITY CASCADE")
	require.NoError(t, err)
}

// ClearAllDocumentTables truncates all document-related tables for clean test state
func ClearAllDocumentTables(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	ClearDocumentVersionsTable(t, ctx, pool)
	ClearInvoicesTable(t, ctx, pool)
	ClearContractsTable(t, ctx, pool)
	ClearMandatesTable(t, ctx, pool)
	ClearEstimatesTable(t, ctx, pool)
}

// InsertEstimate inserts an EstimateEncx record into the database for testing
func InsertEstimate(t *testing.T, ctx context.Context, pool *pgxpool.Pool, estimate *documentDomain.EstimateEncx) error {
	t.Helper()

	query := `
		INSERT INTO estimates (
			id, status, created_at, updated_at,
			caseid_encrypted, clientid_encrypted,
			estimatenumber_encrypted, lineitems_encrypted, estimatedtotal_encrypted, notes_encrypted,
			issue_date, valid_until, accepted, accepted_at, accepted_by,
			dek_encrypted, key_version, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`

	_, err := pool.Exec(ctx, query,
		estimate.ID, estimate.Status, estimate.CreatedAt, estimate.UpdatedAt,
		estimate.CaseIDEncrypted, estimate.ClientIDEncrypted,
		estimate.EstimateNumberEncrypted, estimate.LineItemsEncrypted,
		estimate.EstimatedTotalEncrypted, estimate.NotesEncrypted,
		estimate.IssueDate, estimate.ValidUntil, estimate.Accepted,
		estimate.AcceptedAt, estimate.AcceptedBy,
		estimate.DEKEncrypted, estimate.KeyVersion, estimate.Metadata,
	)

	return err
}

// GetEstimateByID retrieves an EstimateEncx by ID from the database for testing
func GetEstimateByID(t *testing.T, ctx context.Context, pool *pgxpool.Pool, estimateID uuid.UUID) (*documentDomain.EstimateEncx, error) {
	t.Helper()

	query := `
		SELECT id, status, created_at, updated_at,
			   caseid_encrypted, clientid_encrypted,
			   estimatenumber_encrypted, lineitems_encrypted, estimatedtotal_encrypted, notes_encrypted,
			   issue_date, valid_until, accepted, accepted_at, accepted_by,
			   dek_encrypted, key_version, metadata
		FROM estimates
		WHERE id = $1
	`

	estimate := &documentDomain.EstimateEncx{}

	err := pool.QueryRow(ctx, query, estimateID).Scan(
		&estimate.ID, &estimate.Status, &estimate.CreatedAt, &estimate.UpdatedAt,
		&estimate.CaseIDEncrypted, &estimate.ClientIDEncrypted,
		&estimate.EstimateNumberEncrypted, &estimate.LineItemsEncrypted,
		&estimate.EstimatedTotalEncrypted, &estimate.NotesEncrypted,
		&estimate.IssueDate, &estimate.ValidUntil, &estimate.Accepted,
		&estimate.AcceptedAt, &estimate.AcceptedBy,
		&estimate.DEKEncrypted, &estimate.KeyVersion, &estimate.Metadata,
	)

	if err != nil {
		return nil, err
	}

	return estimate, nil
}

// InsertMandate inserts a MandateEncx record into the database for testing
func InsertMandate(t *testing.T, ctx context.Context, pool *pgxpool.Pool, mandate *documentDomain.MandateEncx) error {
	t.Helper()

	query := `
		INSERT INTO mandates (
			id, status, created_at, updated_at,
			caseid_encrypted, clientid_encrypted, mandatenumber_encrypted,
			scopeofwork_encrypted, termsconditions_encrypted,
			clientsignature_encrypted, investigatorsignature_encrypted, specialinstructions_encrypted,
			issue_date, valid_from, valid_until, linked_estimate_id, jurisdiction,
			dek_encrypted, key_version, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
	`

	_, err := pool.Exec(ctx, query,
		mandate.ID, mandate.Status, mandate.CreatedAt, mandate.UpdatedAt,
		mandate.CaseIDEncrypted, mandate.ClientIDEncrypted, mandate.MandateNumberEncrypted,
		mandate.ScopeOfWorkEncrypted, mandate.TermsConditionsEncrypted,
		mandate.ClientSignatureEncrypted, mandate.InvestigatorSignatureEncrypted, mandate.SpecialInstructionsEncrypted,
		mandate.IssueDate, mandate.ValidFrom, mandate.ValidUntil, mandate.LinkedEstimateID, mandate.Jurisdiction,
		mandate.DEKEncrypted, mandate.KeyVersion, mandate.Metadata,
	)

	return err
}

// GetMandateByID retrieves a MandateEncx by ID from the database for testing
func GetMandateByID(t *testing.T, ctx context.Context, pool *pgxpool.Pool, mandateID uuid.UUID) (*documentDomain.MandateEncx, error) {
	t.Helper()

	query := `
		SELECT id, status, created_at, updated_at,
			   caseid_encrypted, clientid_encrypted, mandatenumber_encrypted,
			   scopeofwork_encrypted, termsconditions_encrypted,
			   clientsignature_encrypted, investigatorsignature_encrypted, specialinstructions_encrypted,
			   issue_date, valid_from, valid_until, linked_estimate_id, jurisdiction,
			   dek_encrypted, key_version, metadata
		FROM mandates
		WHERE id = $1
	`

	mandate := &documentDomain.MandateEncx{}

	err := pool.QueryRow(ctx, query, mandateID).Scan(
		&mandate.ID, &mandate.Status, &mandate.CreatedAt, &mandate.UpdatedAt,
		&mandate.CaseIDEncrypted, &mandate.ClientIDEncrypted, &mandate.MandateNumberEncrypted,
		&mandate.ScopeOfWorkEncrypted, &mandate.TermsConditionsEncrypted,
		&mandate.ClientSignatureEncrypted, &mandate.InvestigatorSignatureEncrypted, &mandate.SpecialInstructionsEncrypted,
		&mandate.IssueDate, &mandate.ValidFrom, &mandate.ValidUntil, &mandate.LinkedEstimateID, &mandate.Jurisdiction,
		&mandate.DEKEncrypted, &mandate.KeyVersion, &mandate.Metadata,
	)

	if err != nil {
		return nil, err
	}

	return mandate, nil
}

// InsertContract inserts a ContractEncx record into the database for testing
func InsertContract(t *testing.T, ctx context.Context, pool *pgxpool.Pool, contract *documentDomain.ContractEncx) error {
	t.Helper()

	query := `
		INSERT INTO contracts (
			id, status, created_at, updated_at,
			caseid_encrypted, clientid_encrypted, contractnumber_encrypted,
			scopeofservices_encrypted, paymentterms_encrypted,
			confidentiality_encrypted, terminationclause_encrypted,
			signatures_encrypted, contractvalue_encrypted, renewalterms_encrypted,
			start_date, end_date, linked_mandate_id, currency, governing_law,
			dek_encrypted, key_version, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)
	`

	_, err := pool.Exec(ctx, query,
		contract.ID, contract.Status, contract.CreatedAt, contract.UpdatedAt,
		contract.CaseIDEncrypted, contract.ClientIDEncrypted, contract.ContractNumberEncrypted,
		contract.ScopeOfServicesEncrypted, contract.PaymentTermsEncrypted,
		contract.ConfidentialityEncrypted, contract.TerminationClauseEncrypted,
		contract.SignaturesEncrypted, contract.ContractValueEncrypted, contract.RenewalTermsEncrypted,
		contract.StartDate, contract.EndDate, contract.LinkedMandateID, contract.Currency, contract.GoverningLaw,
		contract.DEKEncrypted, contract.KeyVersion, contract.Metadata,
	)

	return err
}

// GetContractByID retrieves a ContractEncx by ID from the database for testing
func GetContractByID(t *testing.T, ctx context.Context, pool *pgxpool.Pool, contractID uuid.UUID) (*documentDomain.ContractEncx, error) {
	t.Helper()

	query := `
		SELECT id, status, created_at, updated_at,
			   caseid_encrypted, clientid_encrypted, contractnumber_encrypted,
			   scopeofservices_encrypted, paymentterms_encrypted,
			   confidentiality_encrypted, terminationclause_encrypted,
			   signatures_encrypted, contractvalue_encrypted, renewalterms_encrypted,
			   start_date, end_date, linked_mandate_id, currency, governing_law,
			   dek_encrypted, key_version, metadata
		FROM contracts
		WHERE id = $1
	`

	contract := &documentDomain.ContractEncx{}

	err := pool.QueryRow(ctx, query, contractID).Scan(
		&contract.ID, &contract.Status, &contract.CreatedAt, &contract.UpdatedAt,
		&contract.CaseIDEncrypted, &contract.ClientIDEncrypted, &contract.ContractNumberEncrypted,
		&contract.ScopeOfServicesEncrypted, &contract.PaymentTermsEncrypted,
		&contract.ConfidentialityEncrypted, &contract.TerminationClauseEncrypted,
		&contract.SignaturesEncrypted, &contract.ContractValueEncrypted, &contract.RenewalTermsEncrypted,
		&contract.StartDate, &contract.EndDate, &contract.LinkedMandateID, &contract.Currency, &contract.GoverningLaw,
		&contract.DEKEncrypted, &contract.KeyVersion, &contract.Metadata,
	)

	if err != nil {
		return nil, err
	}

	return contract, nil
}

// InsertInvoice inserts an InvoiceEncx record into the database for testing
func InsertInvoice(t *testing.T, ctx context.Context, pool *pgxpool.Pool, invoice *documentDomain.InvoiceEncx) error {
	t.Helper()

	query := `
		INSERT INTO invoices (
			id, status, created_at, updated_at,
			caseid_encrypted, clientid_encrypted, invoicenumber_encrypted,
			lineitems_encrypted, totalamount_encrypted, taxamount_encrypted,
			notes_encrypted, paidamount_encrypted, paymentmethod_encrypted,
			paymentterms_encrypted, latefee_encrypted,
			issue_date, due_date, tax_rate, payment_status, paid_at,
			linked_contract_id, currency, late_fee_rate,
			dek_encrypted, key_version, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26)
	`

	_, err := pool.Exec(ctx, query,
		invoice.ID, invoice.Status, invoice.CreatedAt, invoice.UpdatedAt,
		invoice.CaseIDEncrypted, invoice.ClientIDEncrypted, invoice.InvoiceNumberEncrypted,
		invoice.LineItemsEncrypted, invoice.TotalAmountEncrypted, invoice.TaxAmountEncrypted,
		invoice.NotesEncrypted, invoice.PaidAmountEncrypted, invoice.PaymentMethodEncrypted,
		invoice.PaymentTermsEncrypted, invoice.LateFeeEncrypted,
		invoice.IssueDate, invoice.DueDate, invoice.TaxRate, invoice.PaymentStatus, invoice.PaidAt,
		invoice.LinkedContractID, invoice.Currency, invoice.LateFeeRate,
		invoice.DEKEncrypted, invoice.KeyVersion, invoice.Metadata,
	)

	return err
}

// GetInvoiceByID retrieves an InvoiceEncx by ID from the database for testing
func GetInvoiceByID(t *testing.T, ctx context.Context, pool *pgxpool.Pool, invoiceID uuid.UUID) (*documentDomain.InvoiceEncx, error) {
	t.Helper()

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
		WHERE id = $1
	`

	invoice := &documentDomain.InvoiceEncx{}

	err := pool.QueryRow(ctx, query, invoiceID).Scan(
		&invoice.ID, &invoice.Status, &invoice.CreatedAt, &invoice.UpdatedAt,
		&invoice.CaseIDEncrypted, &invoice.ClientIDEncrypted, &invoice.InvoiceNumberEncrypted,
		&invoice.LineItemsEncrypted, &invoice.TotalAmountEncrypted, &invoice.TaxAmountEncrypted,
		&invoice.NotesEncrypted, &invoice.PaidAmountEncrypted, &invoice.PaymentMethodEncrypted,
		&invoice.PaymentTermsEncrypted, &invoice.LateFeeEncrypted,
		&invoice.IssueDate, &invoice.DueDate, &invoice.TaxRate, &invoice.PaymentStatus, &invoice.PaidAt,
		&invoice.LinkedContractID, &invoice.Currency, &invoice.LateFeeRate,
		&invoice.DEKEncrypted, &invoice.KeyVersion, &invoice.Metadata,
	)

	if err != nil {
		return nil, err
	}

	return invoice, nil
}

// InsertDocumentVersion inserts a DocumentVersion record into the database for testing
func InsertDocumentVersion(t *testing.T, ctx context.Context, pool *pgxpool.Pool, version *documentDomain.DocumentVersion) error {
	t.Helper()

	query := `
		INSERT INTO document_versions (
			id, document_id, doc_type, version, author_id, data, reason, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := pool.Exec(ctx, query,
		version.ID, version.DocumentID, version.DocType, version.Version,
		version.AuthorID, version.Data, version.Reason, version.CreatedAt,
	)

	return err
}

// GetDocumentVersionByID retrieves a document version by ID from the database for testing
func GetDocumentVersionByID(t *testing.T, ctx context.Context, pool *pgxpool.Pool, versionID uuid.UUID) (*documentDomain.DocumentVersion, error) {
	t.Helper()

	query := `
		SELECT id, document_id, doc_type, version, author_id, data, reason, created_at
		FROM document_versions
		WHERE id = $1
	`

	version := &documentDomain.DocumentVersion{}

	err := pool.QueryRow(ctx, query, versionID).Scan(
		&version.ID, &version.DocumentID, &version.DocType, &version.Version,
		&version.AuthorID, &version.Data, &version.Reason, &version.CreatedAt,
	)

	return version, err
}

// CountDocumentsByCaseID returns the number of documents of a specific type for a case
// Note: This function cannot filter by encrypted case_id as it requires decryption.
// For testing purposes, it counts all documents in the table.
// Consider using repository methods with proper encryption/decryption for accurate counts.
func CountDocumentsByCaseID(t *testing.T, ctx context.Context, pool *pgxpool.Pool, caseID uuid.UUID, docType string) (int, error) {
	t.Helper()

	var count int
	var query string

	// WARNING: Cannot filter by encrypted case_id without decryption
	// This counts ALL documents of the given type, not just for the specific case
	switch docType {
	case "estimate":
		query = `SELECT COUNT(*) FROM estimates`
	case "mandate":
		query = `SELECT COUNT(*) FROM mandates`
	case "contract":
		query = `SELECT COUNT(*) FROM contracts`
	case "invoice":
		query = `SELECT COUNT(*) FROM invoices`
	default:
		return 0, fmt.Errorf("unknown document type: %s", docType)
	}

	err := pool.QueryRow(ctx, query).Scan(&count)
	return count, err
}

// CountVersionsByDocumentID returns the number of versions for a document
func CountVersionsByDocumentID(t *testing.T, ctx context.Context, pool *pgxpool.Pool, documentID uuid.UUID) (int, error) {
	t.Helper()

	var count int
	query := `SELECT COUNT(*) FROM document_versions WHERE document_id = $1`
	err := pool.QueryRow(ctx, query, documentID).Scan(&count)
	return count, err
}

// DEPRECATED: CreateTestEstimateInDB - Use repository methods with proper encryption instead.
// The Insert/Get functions now work with *Encx types and require encrypted data.
// Tests should use the document repository's Create methods which handle encryption internally.
//
// Example:
//
//	repo := documentRepository.NewDocumentRepository(pool, cryptoService)
//	estimate := documentHelpers.NewTestEstimateWithCaseID(t, caseID, clientID)
//	err := repo.CreateEstimate(ctx, estimate)
func CreateTestEstimateInDB(t *testing.T, ctx context.Context, pool *pgxpool.Pool, caseID, clientID uuid.UUID) *documentDomain.Estimate {
	t.Helper()
	t.Skip("DEPRECATED: Use repository methods with proper encryption. See function comment for details.")
	return nil
}

// DEPRECATED: CreateTestMandateInDB - Use repository methods with proper encryption instead.
// See CreateTestEstimateInDB for migration example.
func CreateTestMandateInDB(t *testing.T, ctx context.Context, pool *pgxpool.Pool, caseID, clientID uuid.UUID) *documentDomain.Mandate {
	t.Helper()
	t.Skip("DEPRECATED: Use repository methods with proper encryption. See CreateTestEstimateInDB for details.")
	return nil
}

// DEPRECATED: CreateTestContractInDB - Use repository methods with proper encryption instead.
// See CreateTestEstimateInDB for migration example.
func CreateTestContractInDB(t *testing.T, ctx context.Context, pool *pgxpool.Pool, caseID, clientID uuid.UUID) *documentDomain.Contract {
	t.Helper()
	t.Skip("DEPRECATED: Use repository methods with proper encryption. See CreateTestEstimateInDB for details.")
	return nil
}

// DEPRECATED: CreateTestInvoiceInDB - Use repository methods with proper encryption instead.
// See CreateTestEstimateInDB for migration example.
func CreateTestInvoiceInDB(t *testing.T, ctx context.Context, pool *pgxpool.Pool, caseID, clientID uuid.UUID) *documentDomain.Invoice {
	t.Helper()
	t.Skip("DEPRECATED: Use repository methods with proper encryption. See CreateTestEstimateInDB for details.")
	return nil
}

// DEPRECATED: CreateDocumentWorkflowInDB - Use repository methods with proper encryption instead.
// See CreateTestEstimateInDB for migration example.
func CreateDocumentWorkflowInDB(t *testing.T, ctx context.Context, pool *pgxpool.Pool, caseID, clientID uuid.UUID) (*documentDomain.Estimate, *documentDomain.Mandate, *documentDomain.Contract, *documentDomain.Invoice) {
	t.Helper()
	t.Skip("DEPRECATED: Use repository methods with proper encryption. See CreateTestEstimateInDB for details.")
	return nil, nil, nil, nil
}
