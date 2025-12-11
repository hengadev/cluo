package documentHelpers

import (
	"context"
	"encoding/json"
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

// InsertEstimate inserts an Estimate record into the database for testing
func InsertEstimate(t *testing.T, ctx context.Context, pool *pgxpool.Pool, estimate *documentDomain.Estimate) error {
	t.Helper()

	lineItemsJSON, err := json.Marshal(estimate.LineItems)
	if err != nil {
		return fmt.Errorf("failed to marshal line items: %w", err)
	}

	query := `
		INSERT INTO estimates (
			id, case_id, client_id, status, estimate_number, issue_date, valid_until,
			line_items, estimated_total, notes, accepted, accepted_at, accepted_by,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	_, err = pool.Exec(ctx, query,
		estimate.ID, estimate.CaseID, estimate.ClientID, estimate.Status,
		estimate.EstimateNumber, estimate.IssueDate, estimate.ValidUntil,
		lineItemsJSON, estimate.EstimatedTotal, estimate.Notes,
		estimate.Accepted, estimate.AcceptedAt, estimate.AcceptedBy,
		estimate.CreatedAt, estimate.UpdatedAt,
	)

	return err
}

// GetEstimateByID retrieves an estimate by ID from the database for testing
func GetEstimateByID(t *testing.T, ctx context.Context, pool *pgxpool.Pool, estimateID uuid.UUID) (*documentDomain.Estimate, error) {
	t.Helper()

	query := `
		SELECT id, case_id, client_id, status, estimate_number, issue_date, valid_until,
			   line_items, estimated_total, notes, accepted, accepted_at, accepted_by,
			   created_at, updated_at
		FROM estimates
		WHERE id = $1
	`

	estimate := &documentDomain.Estimate{}
	var lineItemsJSON []byte

	err := pool.QueryRow(ctx, query, estimateID).Scan(
		&estimate.ID, &estimate.CaseID, &estimate.ClientID, &estimate.Status,
		&estimate.EstimateNumber, &estimate.IssueDate, &estimate.ValidUntil,
		&lineItemsJSON, &estimate.EstimatedTotal, &estimate.Notes,
		&estimate.Accepted, &estimate.AcceptedAt, &estimate.AcceptedBy,
		&estimate.CreatedAt, &estimate.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(lineItemsJSON, &estimate.LineItems); err != nil {
		return nil, fmt.Errorf("failed to unmarshal line items: %w", err)
	}

	return estimate, nil
}

// InsertMandate inserts a Mandate record into the database for testing
func InsertMandate(t *testing.T, ctx context.Context, pool *pgxpool.Pool, mandate *documentDomain.Mandate) error {
	t.Helper()

	var clientSignatureJSON, investigatorSignatureJSON []byte
	var err error

	if mandate.ClientSignature != nil {
		clientSignatureJSON, err = json.Marshal(mandate.ClientSignature)
		if err != nil {
			return fmt.Errorf("failed to marshal client signature: %w", err)
		}
	}

	if mandate.InvestigatorSignature != nil {
		investigatorSignatureJSON, err = json.Marshal(mandate.InvestigatorSignature)
		if err != nil {
			return fmt.Errorf("failed to marshal investigator signature: %w", err)
		}
	}

	query := `
		INSERT INTO mandates (
			id, case_id, client_id, status, mandate_number, issue_date, scope_of_work,
			valid_from, valid_until, terms_conditions, client_signature,
			investigator_signature, linked_estimate_id, special_instructions, jurisdiction,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`

	_, err = pool.Exec(ctx, query,
		mandate.ID, mandate.CaseID, mandate.ClientID, mandate.Status,
		mandate.MandateNumber, mandate.IssueDate, mandate.ScopeOfWork,
		mandate.ValidFrom, mandate.ValidUntil, mandate.TermsConditions,
		clientSignatureJSON, investigatorSignatureJSON, mandate.LinkedEstimateID,
		mandate.SpecialInstructions, mandate.Jurisdiction,
		mandate.CreatedAt, mandate.UpdatedAt,
	)

	return err
}

// GetMandateByID retrieves a mandate by ID from the database for testing
func GetMandateByID(t *testing.T, ctx context.Context, pool *pgxpool.Pool, mandateID uuid.UUID) (*documentDomain.Mandate, error) {
	t.Helper()

	query := `
		SELECT id, case_id, client_id, status, mandate_number, issue_date, scope_of_work,
			   valid_from, valid_until, terms_conditions, client_signature,
			   investigator_signature, linked_estimate_id, special_instructions, jurisdiction,
			   created_at, updated_at
		FROM mandates
		WHERE id = $1
	`

	mandate := &documentDomain.Mandate{}
	var clientSignatureJSON, investigatorSignatureJSON []byte

	err := pool.QueryRow(ctx, query, mandateID).Scan(
		&mandate.ID, &mandate.CaseID, &mandate.ClientID, &mandate.Status,
		&mandate.MandateNumber, &mandate.IssueDate, &mandate.ScopeOfWork,
		&mandate.ValidFrom, &mandate.ValidUntil, &mandate.TermsConditions,
		&clientSignatureJSON, &investigatorSignatureJSON, &mandate.LinkedEstimateID,
		&mandate.SpecialInstructions, &mandate.Jurisdiction,
		&mandate.CreatedAt, &mandate.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if clientSignatureJSON != nil {
		if err := json.Unmarshal(clientSignatureJSON, &mandate.ClientSignature); err != nil {
			return nil, fmt.Errorf("failed to unmarshal client signature: %w", err)
		}
	}

	if investigatorSignatureJSON != nil {
		if err := json.Unmarshal(investigatorSignatureJSON, &mandate.InvestigatorSignature); err != nil {
			return nil, fmt.Errorf("failed to unmarshal investigator signature: %w", err)
		}
	}

	return mandate, nil
}

// InsertContract inserts a Contract record into the database for testing
func InsertContract(t *testing.T, ctx context.Context, pool *pgxpool.Pool, contract *documentDomain.Contract) error {
	t.Helper()

	signaturesJSON, err := json.Marshal(contract.Signatures)
	if err != nil {
		return fmt.Errorf("failed to marshal signatures: %w", err)
	}

	query := `
		INSERT INTO contracts (
			id, case_id, client_id, status, contract_number, start_date, end_date,
			scope_of_services, payment_terms, confidentiality, termination_clause,
			signatures, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err = pool.Exec(ctx, query,
		contract.ID, contract.CaseID, contract.ClientID, contract.Status,
		contract.ContractNumber, contract.StartDate, contract.EndDate,
		contract.ScopeOfServices, contract.PaymentTerms, contract.Confidentiality,
		contract.TerminationClause, signaturesJSON,
		contract.CreatedAt, contract.UpdatedAt,
	)

	return err
}

// GetContractByID retrieves a contract by ID from the database for testing
func GetContractByID(t *testing.T, ctx context.Context, pool *pgxpool.Pool, contractID uuid.UUID) (*documentDomain.Contract, error) {
	t.Helper()

	query := `
		SELECT id, case_id, client_id, status, contract_number, start_date, end_date,
			   scope_of_services, payment_terms, confidentiality, termination_clause,
			   signatures, created_at, updated_at
		FROM contracts
		WHERE id = $1
	`

	contract := &documentDomain.Contract{}
	var signaturesJSON []byte

	err := pool.QueryRow(ctx, query, contractID).Scan(
		&contract.ID, &contract.CaseID, &contract.ClientID, &contract.Status,
		&contract.ContractNumber, &contract.StartDate, &contract.EndDate,
		&contract.ScopeOfServices, &contract.PaymentTerms, &contract.Confidentiality,
		&contract.TerminationClause, &signaturesJSON,
		&contract.CreatedAt, &contract.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(signaturesJSON, &contract.Signatures); err != nil {
		return nil, fmt.Errorf("failed to unmarshal signatures: %w", err)
	}

	return contract, nil
}

// InsertInvoice inserts an Invoice record into the database for testing
func InsertInvoice(t *testing.T, ctx context.Context, pool *pgxpool.Pool, invoice *documentDomain.Invoice) error {
	t.Helper()

	lineItemsJSON, err := json.Marshal(invoice.LineItems)
	if err != nil {
		return fmt.Errorf("failed to marshal line items: %w", err)
	}

	query := `
		INSERT INTO invoices (
			id, case_id, client_id, status, invoice_number, issue_date, due_date,
			line_items, total_amount, tax_rate, tax_amount, notes, payment_status,
			paid_at, paid_amount, payment_method, payment_reference, linked_contract_id,
			currency, payment_terms, late_fee, late_fee_rate, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)
	`

	_, err = pool.Exec(ctx, query,
		invoice.ID, invoice.CaseID, invoice.ClientID, invoice.Status,
		invoice.InvoiceNumber, invoice.IssueDate, invoice.DueDate,
		lineItemsJSON, invoice.TotalAmount, invoice.TaxRate, invoice.TaxAmount,
		invoice.Notes, invoice.PaymentStatus, invoice.PaidAt, invoice.PaidAmount,
		invoice.PaymentMethod, invoice.PaymentReference, invoice.LinkedContractID,
		invoice.Currency, invoice.PaymentTerms, invoice.LateFee, invoice.LateFeeRate,
		invoice.CreatedAt, invoice.UpdatedAt,
	)

	return err
}

// GetInvoiceByID retrieves an invoice by ID from the database for testing
func GetInvoiceByID(t *testing.T, ctx context.Context, pool *pgxpool.Pool, invoiceID uuid.UUID) (*documentDomain.Invoice, error) {
	t.Helper()

	query := `
		SELECT id, case_id, client_id, status, invoice_number, issue_date, due_date,
			   line_items, total_amount, tax_rate, tax_amount, notes, payment_status,
			   paid_at, paid_amount, payment_method, payment_reference, linked_contract_id,
			   currency, payment_terms, late_fee, late_fee_rate, created_at, updated_at
		FROM invoices
		WHERE id = $1
	`

	invoice := &documentDomain.Invoice{}
	var lineItemsJSON []byte

	err := pool.QueryRow(ctx, query, invoiceID).Scan(
		&invoice.ID, &invoice.CaseID, &invoice.ClientID, &invoice.Status,
		&invoice.InvoiceNumber, &invoice.IssueDate, &invoice.DueDate,
		&lineItemsJSON, &invoice.TotalAmount, &invoice.TaxRate, &invoice.TaxAmount,
		&invoice.Notes, &invoice.PaymentStatus, &invoice.PaidAt, &invoice.PaidAmount,
		&invoice.PaymentMethod, &invoice.PaymentReference, &invoice.LinkedContractID,
		&invoice.Currency, &invoice.PaymentTerms, &invoice.LateFee, &invoice.LateFeeRate,
		&invoice.CreatedAt, &invoice.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(lineItemsJSON, &invoice.LineItems); err != nil {
		return nil, fmt.Errorf("failed to unmarshal line items: %w", err)
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
func CountDocumentsByCaseID(t *testing.T, ctx context.Context, pool *pgxpool.Pool, caseID uuid.UUID, docType string) (int, error) {
	t.Helper()

	var count int
	var query string

	switch docType {
	case "estimate":
		query = `SELECT COUNT(*) FROM estimates WHERE case_id = $1`
	case "mandate":
		query = `SELECT COUNT(*) FROM mandates WHERE case_id = $1`
	case "contract":
		query = `SELECT COUNT(*) FROM contracts WHERE case_id = $1`
	case "invoice":
		query = `SELECT COUNT(*) FROM invoices WHERE case_id = $1`
	default:
		return 0, fmt.Errorf("unknown document type: %s", docType)
	}

	err := pool.QueryRow(ctx, query, caseID).Scan(&count)
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

// CreateTestEstimateInDB creates and inserts an estimate with the given case and client IDs
func CreateTestEstimateInDB(t *testing.T, ctx context.Context, pool *pgxpool.Pool, caseID, clientID uuid.UUID) *documentDomain.Estimate {
	t.Helper()

	estimate := NewTestEstimateWithCaseID(t, caseID, clientID)
	err := InsertEstimate(t, ctx, pool, estimate)
	require.NoError(t, err)
	return estimate
}

// CreateTestMandateInDB creates and inserts a mandate with the given case and client IDs
func CreateTestMandateInDB(t *testing.T, ctx context.Context, pool *pgxpool.Pool, caseID, clientID uuid.UUID) *documentDomain.Mandate {
	t.Helper()

	mandate := NewTestMandateWithCaseID(t, caseID, clientID)
	err := InsertMandate(t, ctx, pool, mandate)
	require.NoError(t, err)
	return mandate
}

// CreateTestContractInDB creates and inserts a contract with the given case and client IDs
func CreateTestContractInDB(t *testing.T, ctx context.Context, pool *pgxpool.Pool, caseID, clientID uuid.UUID) *documentDomain.Contract {
	t.Helper()

	contract := NewTestContractWithCaseID(t, caseID, clientID)
	err := InsertContract(t, ctx, pool, contract)
	require.NoError(t, err)
	return contract
}

// CreateTestInvoiceInDB creates and inserts an invoice with the given case and client IDs
func CreateTestInvoiceInDB(t *testing.T, ctx context.Context, pool *pgxpool.Pool, caseID, clientID uuid.UUID) *documentDomain.Invoice {
	t.Helper()

	invoice := NewTestInvoiceWithCaseID(t, caseID, clientID)
	err := InsertInvoice(t, ctx, pool, invoice)
	require.NoError(t, err)
	return invoice
}

// CreateDocumentWorkflowInDB creates a complete document workflow in the database:
// Estimate -> Mandate -> Contract -> Invoice
func CreateDocumentWorkflowInDB(t *testing.T, ctx context.Context, pool *pgxpool.Pool, caseID, clientID uuid.UUID) (*documentDomain.Estimate, *documentDomain.Mandate, *documentDomain.Contract, *documentDomain.Invoice) {
	t.Helper()

	estimate, mandate, contract, invoice := CreateDocumentWorkflow(t, caseID, clientID)

	// Insert all documents into the database
	err := InsertEstimate(t, ctx, pool, estimate)
	require.NoError(t, err)

	err = InsertMandate(t, ctx, pool, mandate)
	require.NoError(t, err)

	err = InsertContract(t, ctx, pool, contract)
	require.NoError(t, err)

	err = InsertInvoice(t, ctx, pool, invoice)
	require.NoError(t, err)

	return estimate, mandate, contract, invoice
}