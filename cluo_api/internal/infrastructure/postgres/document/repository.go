package documentRepository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hengadev/cluo_api/internal/domain"
	"github.com/hengadev/cluo_api/internal/ports"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) ports.DocumentRepository {
	return &Repository{pool: pool}
}

// Generic document operations
func (r *Repository) Create(ctx context.Context, doc domain.Documentable) error {
	switch d := doc.(type) {
	case *domain.Estimate:
		return r.CreateEstimate(ctx, d)
	case *domain.Mandate:
		return r.CreateMandate(ctx, d)
	case *domain.Contract:
		return r.CreateContract(ctx, d)
	case *domain.Invoice:
		return r.CreateInvoice(ctx, d)
	default:
		return fmt.Errorf("unsupported document type: %T", doc)
	}
}

func (r *Repository) GetByID(ctx context.Context, id string, docType domain.DocumentType) (domain.Documentable, error) {
	switch docType {
	case domain.DocumentTypeEstimate:
		return r.GetEstimateByID(ctx, id)
	case domain.DocumentTypeMandate:
		return r.GetMandateByID(ctx, id)
	case domain.DocumentTypeContract:
		return r.GetContractByID(ctx, id)
	case domain.DocumentTypeInvoice:
		return r.GetInvoiceByID(ctx, id)
	default:
		return nil, fmt.Errorf("unsupported document type: %s", docType)
	}
}

func (r *Repository) Update(ctx context.Context, doc domain.Documentable) error {
	switch d := doc.(type) {
	case *domain.Estimate:
		return r.UpdateEstimate(ctx, d)
	case *domain.Mandate:
		return r.UpdateMandate(ctx, d)
	case *domain.Contract:
		return r.UpdateContract(ctx, d)
	case *domain.Invoice:
		return r.UpdateInvoice(ctx, d)
	default:
		return fmt.Errorf("unsupported document type: %T", doc)
	}
}

func (r *Repository) Delete(ctx context.Context, id string, docType domain.DocumentType) error {
	switch docType {
	case domain.DocumentTypeEstimate:
		return r.DeleteEstimate(ctx, id)
	case domain.DocumentTypeMandate:
		return r.DeleteMandate(ctx, id)
	case domain.DocumentTypeContract:
		return r.DeleteContract(ctx, id)
	case domain.DocumentTypeInvoice:
		return r.DeleteInvoice(ctx, id)
	default:
		return fmt.Errorf("unsupported document type: %s", docType)
	}
}

func (r *Repository) List(ctx context.Context, filter domain.DocumentFilter, pagination domain.Pagination) ([]domain.DocumentSummary, int, error) {
	// Build WHERE clause
	whereConditions := []string{}
	args := []interface{}{}
	argIndex := 1

	if filter.Type != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("doc_type = $%d", argIndex))
		args = append(args, *filter.Type)
		argIndex++
	}

	if filter.Status != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *filter.Status)
		argIndex++
	}

	if filter.CaseID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("case_id = $%d", argIndex))
		args = append(args, *filter.CaseID)
		argIndex++
	}

	if filter.ClientID != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("client_id = $%d", argIndex))
		args = append(args, *filter.ClientID)
		argIndex++
	}

	if filter.DateFrom != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("created_at >= $%d", argIndex))
		args = append(args, *filter.DateFrom)
		argIndex++
	}

	if filter.DateTo != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("created_at <= $%d", argIndex))
		args = append(args, *filter.DateTo)
		argIndex++
	}

	if filter.Search != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("(document_ref ILIKE $%d OR doc_type ILIKE $%d)", argIndex, argIndex+1))
		args = append(args, "%"+*filter.Search+"%", "%"+*filter.Search+"%")
		argIndex += 2
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Get total count
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM (
			SELECT id, case_id, client_id, 'estimate' as doc_type, status, estimate_number as document_ref, created_at, updated_at FROM estimates
			UNION ALL
			SELECT id, case_id, client_id, 'mandate' as doc_type, status, mandate_number as document_ref, created_at, updated_at FROM mandates
			UNION ALL
			SELECT id, case_id, client_id, 'contract' as doc_type, status, contract_number as document_ref, created_at, updated_at FROM contracts
			UNION ALL
			SELECT id, case_id, client_id, 'invoice' as doc_type, status, invoice_number as document_ref, created_at, updated_at FROM invoices
		) documents %s
	`, whereClause)

	var total int
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count documents: %w", err)
	}

	// Get paginated results
	offset := pagination.GetOffset()
	limit := pagination.PageSize
	args = append(args, limit, offset)

	query := fmt.Sprintf(`
		SELECT * FROM (
			SELECT id, case_id, client_id, 'estimate' as doc_type, status, estimate_number as document_ref, created_at, updated_at FROM estimates
			UNION ALL
			SELECT id, case_id, client_id, 'mandate' as doc_type, status, mandate_number as document_ref, created_at, updated_at FROM mandates
			UNION ALL
			SELECT id, case_id, client_id, 'contract' as doc_type, status, contract_number as document_ref, created_at, updated_at FROM contracts
			UNION ALL
			SELECT id, case_id, client_id, 'invoice' as doc_type, status, invoice_number as document_ref, created_at, updated_at FROM invoices
		) documents %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query documents: %w", err)
	}
	defer rows.Close()

	var documents []domain.DocumentSummary
	for rows.Next() {
		var doc domain.DocumentSummary
		err := rows.Scan(
			&doc.ID, &doc.CaseID, &doc.ClientID, &doc.Type, &doc.Status,
			&doc.DocumentRef, &doc.CreatedAt, &doc.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan document: %w", err)
		}
		documents = append(documents, doc)
	}

	return documents, total, nil
}

// Estimate operations
func (r *Repository) CreateEstimate(ctx context.Context, estimate *domain.Estimate) error {
	lineItemsJSON, err := json.Marshal(estimate.LineItems)
	if err != nil {
		return fmt.Errorf("failed to marshal line items: %w", err)
	}

	query := `
		INSERT INTO estimates (
			id, case_id, client_id, status, estimate_number, issue_date, valid_until,
			line_items, estimated_total, notes, accepted, accepted_at, accepted_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err = r.pool.Exec(ctx, query,
		estimate.ID, estimate.CaseID, estimate.ClientID, estimate.Status,
		estimate.EstimateNumber, estimate.IssueDate, estimate.ValidUntil,
		lineItemsJSON, estimate.EstimatedTotal, estimate.Notes,
		estimate.Accepted, estimate.AcceptedAt, estimate.AcceptedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to create estimate: %w", err)
	}

	return nil
}

func (r *Repository) GetEstimateByID(ctx context.Context, id string) (*domain.Estimate, error) {
	query := `
		SELECT id, case_id, client_id, status, estimate_number, issue_date, valid_until,
			   line_items, estimated_total, notes, accepted, accepted_at, accepted_by,
			   created_at, updated_at
		FROM estimates
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	var estimate domain.Estimate
	var lineItemsJSON []byte

	err := row.Scan(
		&estimate.ID, &estimate.CaseID, &estimate.ClientID, &estimate.Status,
		&estimate.EstimateNumber, &estimate.IssueDate, &estimate.ValidUntil,
		&lineItemsJSON, &estimate.EstimatedTotal, &estimate.Notes,
		&estimate.Accepted, &estimate.AcceptedAt, &estimate.AcceptedBy,
		&estimate.CreatedAt, &estimate.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("estimate not found")
		}
		return nil, fmt.Errorf("failed to get estimate: %w", err)
	}

	if err := json.Unmarshal(lineItemsJSON, &estimate.LineItems); err != nil {
		return nil, fmt.Errorf("failed to unmarshal line items: %w", err)
	}

	return &estimate, nil
}

func (r *Repository) UpdateEstimate(ctx context.Context, estimate *domain.Estimate) error {
	lineItemsJSON, err := json.Marshal(estimate.LineItems)
	if err != nil {
		return fmt.Errorf("failed to marshal line items: %w", err)
	}

	query := `
		UPDATE estimates SET
			case_id = $2, client_id = $3, status = $4, estimate_number = $5,
			issue_date = $6, valid_until = $7, line_items = $8, estimated_total = $9,
			notes = $10, accepted = $11, accepted_at = $12, accepted_by = $13,
			updated_at = NOW()
		WHERE id = $1
	`

	_, err = r.pool.Exec(ctx, query,
		estimate.ID, estimate.CaseID, estimate.ClientID, estimate.Status,
		estimate.EstimateNumber, estimate.IssueDate, estimate.ValidUntil,
		lineItemsJSON, estimate.EstimatedTotal, estimate.Notes,
		estimate.Accepted, estimate.AcceptedAt, estimate.AcceptedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to update estimate: %w", err)
	}

	return nil
}

func (r *Repository) DeleteEstimate(ctx context.Context, id string) error {
	query := `DELETE FROM estimates WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete estimate: %w", err)
	}
	return nil
}

func (r *Repository) ListEstimatesByCase(ctx context.Context, caseID string, pagination domain.Pagination) ([]*domain.Estimate, int, error) {
	// Get total count
	var total int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM estimates WHERE case_id = $1", caseID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count estimates: %w", err)
	}

	// Get paginated results
	offset := pagination.GetOffset()
	limit := pagination.PageSize
	query := `
		SELECT id, case_id, client_id, status, estimate_number, issue_date, valid_until,
			   line_items, estimated_total, notes, accepted, accepted_at, accepted_by,
			   created_at, updated_at
		FROM estimates
		WHERE case_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, caseID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query estimates: %w", err)
	}
	defer rows.Close()

	var estimates []*domain.Estimate
	for rows.Next() {
		var estimate domain.Estimate
		var lineItemsJSON []byte

		err := rows.Scan(
			&estimate.ID, &estimate.CaseID, &estimate.ClientID, &estimate.Status,
			&estimate.EstimateNumber, &estimate.IssueDate, &estimate.ValidUntil,
			&lineItemsJSON, &estimate.EstimatedTotal, &estimate.Notes,
			&estimate.Accepted, &estimate.AcceptedAt, &estimate.AcceptedBy,
			&estimate.CreatedAt, &estimate.UpdatedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan estimate: %w", err)
		}

		if err := json.Unmarshal(lineItemsJSON, &estimate.LineItems); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal line items: %w", err)
		}

		estimates = append(estimates, &estimate)
	}

	return estimates, total, nil
}

// Mandate operations
func (r *Repository) CreateMandate(ctx context.Context, mandate *domain.Mandate) error {
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
			investigator_signature, linked_estimate_id, special_instructions, jurisdiction
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`

	_, err = r.pool.Exec(ctx, query,
		mandate.ID, mandate.CaseID, mandate.ClientID, mandate.Status,
		mandate.MandateNumber, mandate.IssueDate, mandate.ScopeOfWork,
		mandate.ValidFrom, mandate.ValidUntil, mandate.TermsConditions,
		clientSignatureJSON, investigatorSignatureJSON, mandate.LinkedEstimateID,
		mandate.SpecialInstructions, mandate.Jurisdiction,
	)

	if err != nil {
		return fmt.Errorf("failed to create mandate: %w", err)
	}

	return nil
}

func (r *Repository) GetMandateByID(ctx context.Context, id string) (*domain.Mandate, error) {
	query := `
		SELECT id, case_id, client_id, status, mandate_number, issue_date, scope_of_work,
			   valid_from, valid_until, terms_conditions, client_signature,
			   investigator_signature, linked_estimate_id, special_instructions, jurisdiction,
			   created_at, updated_at
		FROM mandates
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	var mandate domain.Mandate
	var clientSignatureJSON, investigatorSignatureJSON []byte

	err := row.Scan(
		&mandate.ID, &mandate.CaseID, &mandate.ClientID, &mandate.Status,
		&mandate.MandateNumber, &mandate.IssueDate, &mandate.ScopeOfWork,
		&mandate.ValidFrom, &mandate.ValidUntil, &mandate.TermsConditions,
		&clientSignatureJSON, &investigatorSignatureJSON, &mandate.LinkedEstimateID,
		&mandate.SpecialInstructions, &mandate.Jurisdiction,
		&mandate.CreatedAt, &mandate.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("mandate not found")
		}
		return nil, fmt.Errorf("failed to get mandate: %w", err)
	}

	if len(clientSignatureJSON) > 0 {
		if err := json.Unmarshal(clientSignatureJSON, &mandate.ClientSignature); err != nil {
			return nil, fmt.Errorf("failed to unmarshal client signature: %w", err)
		}
	}

	if len(investigatorSignatureJSON) > 0 {
		if err := json.Unmarshal(investigatorSignatureJSON, &mandate.InvestigatorSignature); err != nil {
			return nil, fmt.Errorf("failed to unmarshal investigator signature: %w", err)
		}
	}

	return &mandate, nil
}

func (r *Repository) UpdateMandate(ctx context.Context, mandate *domain.Mandate) error {
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
		UPDATE mandates SET
			case_id = $2, client_id = $3, status = $4, mandate_number = $5,
			issue_date = $6, scope_of_work = $7, valid_from = $8, valid_until = $9,
			terms_conditions = $10, client_signature = $11, investigator_signature = $12,
			linked_estimate_id = $13, special_instructions = $14, jurisdiction = $15,
			updated_at = NOW()
		WHERE id = $1
	`

	_, err = r.pool.Exec(ctx, query,
		mandate.ID, mandate.CaseID, mandate.ClientID, mandate.Status,
		mandate.MandateNumber, mandate.IssueDate, mandate.ScopeOfWork,
		mandate.ValidFrom, mandate.ValidUntil, mandate.TermsConditions,
		clientSignatureJSON, investigatorSignatureJSON, mandate.LinkedEstimateID,
		mandate.SpecialInstructions, mandate.Jurisdiction,
	)

	if err != nil {
		return fmt.Errorf("failed to update mandate: %w", err)
	}

	return nil
}

func (r *Repository) DeleteMandate(ctx context.Context, id string) error {
	query := `DELETE FROM mandates WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete mandate: %w", err)
	}
	return nil
}

func (r *Repository) ListMandatesByCase(ctx context.Context, caseID string, pagination domain.Pagination) ([]*domain.Mandate, int, error) {
	// Get total count
	var total int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM mandates WHERE case_id = $1", caseID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count mandates: %w", err)
	}

	// Get paginated results
	offset := pagination.GetOffset()
	limit := pagination.PageSize
	query := `
		SELECT id, case_id, client_id, status, mandate_number, issue_date, scope_of_work,
			   valid_from, valid_until, terms_conditions, client_signature,
			   investigator_signature, linked_estimate_id, special_instructions, jurisdiction,
			   created_at, updated_at
		FROM mandates
		WHERE case_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, caseID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query mandates: %w", err)
	}
	defer rows.Close()

	var mandates []*domain.Mandate
	for rows.Next() {
		var mandate domain.Mandate
		var clientSignatureJSON, investigatorSignatureJSON []byte

		err := rows.Scan(
			&mandate.ID, &mandate.CaseID, &mandate.ClientID, &mandate.Status,
			&mandate.MandateNumber, &mandate.IssueDate, &mandate.ScopeOfWork,
			&mandate.ValidFrom, &mandate.ValidUntil, &mandate.TermsConditions,
			&clientSignatureJSON, &investigatorSignatureJSON, &mandate.LinkedEstimateID,
			&mandate.SpecialInstructions, &mandate.Jurisdiction,
			&mandate.CreatedAt, &mandate.UpdatedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan mandate: %w", err)
		}

		if len(clientSignatureJSON) > 0 {
			if err := json.Unmarshal(clientSignatureJSON, &mandate.ClientSignature); err != nil {
				return nil, 0, fmt.Errorf("failed to unmarshal client signature: %w", err)
			}
		}

		if len(investigatorSignatureJSON) > 0 {
			if err := json.Unmarshal(investigatorSignatureJSON, &mandate.InvestigatorSignature); err != nil {
				return nil, 0, fmt.Errorf("failed to unmarshal investigator signature: %w", err)
			}
		}

		mandates = append(mandates, &mandate)
	}

	return mandates, total, nil
}

// Contract operations
func (r *Repository) CreateContract(ctx context.Context, contract *domain.Contract) error {
	signaturesJSON, err := json.Marshal(contract.Signatures)
	if err != nil {
		return fmt.Errorf("failed to marshal signatures: %w", err)
	}

	query := `
		INSERT INTO contracts (
			id, case_id, client_id, status, contract_number, start_date, end_date,
			scope_of_services, payment_terms, confidentiality, termination_clause,
			signatures, linked_mandate_id, contract_value, currency, renewal_terms, governing_law
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`

	_, err = r.pool.Exec(ctx, query,
		contract.ID, contract.CaseID, contract.ClientID, contract.Status,
		contract.ContractNumber, contract.StartDate, contract.EndDate,
		contract.ScopeOfServices, contract.PaymentTerms, contract.Confidentiality,
		contract.TerminationClause, signaturesJSON, contract.LinkedMandateID,
		contract.ContractValue, contract.Currency, contract.RenewalTerms, contract.GoverningLaw,
	)

	if err != nil {
		return fmt.Errorf("failed to create contract: %w", err)
	}

	return nil
}

func (r *Repository) GetContractByID(ctx context.Context, id string) (*domain.Contract, error) {
	query := `
		SELECT id, case_id, client_id, status, contract_number, start_date, end_date,
			   scope_of_services, payment_terms, confidentiality, termination_clause,
			   signatures, linked_mandate_id, contract_value, currency, renewal_terms, governing_law,
			   created_at, updated_at
		FROM contracts
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	var contract domain.Contract
	var signaturesJSON []byte

	err := row.Scan(
		&contract.ID, &contract.CaseID, &contract.ClientID, &contract.Status,
		&contract.ContractNumber, &contract.StartDate, &contract.EndDate,
		&contract.ScopeOfServices, &contract.PaymentTerms, &contract.Confidentiality,
		&contract.TerminationClause, &signaturesJSON, &contract.LinkedMandateID,
		&contract.ContractValue, &contract.Currency, &contract.RenewalTerms, &contract.GoverningLaw,
		&contract.CreatedAt, &contract.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("contract not found")
		}
		return nil, fmt.Errorf("failed to get contract: %w", err)
	}

	if err := json.Unmarshal(signaturesJSON, &contract.Signatures); err != nil {
		return nil, fmt.Errorf("failed to unmarshal signatures: %w", err)
	}

	return &contract, nil
}

func (r *Repository) UpdateContract(ctx context.Context, contract *domain.Contract) error {
	signaturesJSON, err := json.Marshal(contract.Signatures)
	if err != nil {
		return fmt.Errorf("failed to marshal signatures: %w", err)
	}

	query := `
		UPDATE contracts SET
			case_id = $2, client_id = $3, status = $4, contract_number = $5,
			start_date = $6, end_date = $7, scope_of_services = $8, payment_terms = $9,
			confidentiality = $10, termination_clause = $11, signatures = $12,
			linked_mandate_id = $13, contract_value = $14, currency = $15,
			renewal_terms = $16, governing_law = $17, updated_at = NOW()
		WHERE id = $1
	`

	_, err = r.pool.Exec(ctx, query,
		contract.ID, contract.CaseID, contract.ClientID, contract.Status,
		contract.ContractNumber, contract.StartDate, contract.EndDate,
		contract.ScopeOfServices, contract.PaymentTerms, contract.Confidentiality,
		contract.TerminationClause, signaturesJSON, contract.LinkedMandateID,
		contract.ContractValue, contract.Currency, contract.RenewalTerms, contract.GoverningLaw,
	)

	if err != nil {
		return fmt.Errorf("failed to update contract: %w", err)
	}

	return nil
}

func (r *Repository) DeleteContract(ctx context.Context, id string) error {
	query := `DELETE FROM contracts WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete contract: %w", err)
	}
	return nil
}

func (r *Repository) ListContractsByCase(ctx context.Context, caseID string, pagination domain.Pagination) ([]*domain.Contract, int, error) {
	// Get total count
	var total int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM contracts WHERE case_id = $1", caseID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count contracts: %w", err)
	}

	// Get paginated results
	offset := pagination.GetOffset()
	limit := pagination.PageSize
	query := `
		SELECT id, case_id, client_id, status, contract_number, start_date, end_date,
			   scope_of_services, payment_terms, confidentiality, termination_clause,
			   signatures, linked_mandate_id, contract_value, currency, renewal_terms, governing_law,
			   created_at, updated_at
		FROM contracts
		WHERE case_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.pool.Query(ctx, query, caseID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query contracts: %w", err)
	}
	defer rows.Close()

	var contracts []*domain.Contract
	for rows.Next() {
		var contract domain.Contract
		var signaturesJSON []byte

		err := rows.Scan(
			&contract.ID, &contract.CaseID, &contract.ClientID, &contract.Status,
			&contract.ContractNumber, &contract.StartDate, &contract.EndDate,
			&contract.ScopeOfServices, &contract.PaymentTerms, &contract.Confidentiality,
			&contract.TerminationClause, &signaturesJSON, &contract.LinkedMandateID,
			&contract.ContractValue, &contract.Currency, &contract.RenewalTerms, &contract.GoverningLaw,
			&contract.CreatedAt, &contract.UpdatedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan contract: %w", err)
		}

		if err := json.Unmarshal(signaturesJSON, &contract.Signatures); err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal signatures: %w", err)
		}

		contracts = append(contracts, &contract)
	}

	return contracts, total, nil
}

// Invoice operations
func (r *Repository) CreateInvoice(ctx context.Context, invoice *domain.Invoice) error {
	lineItemsJSON, err := json.Marshal(invoice.LineItems)
	if err != nil {
		return fmt.Errorf("failed to marshal line items: %w", err)
	}

	query := `
		INSERT INTO invoices (
			id, case_id, client_id, status, invoice_number, issue_date, due_date,
			line_items, total_amount, tax_rate, tax_amount, notes, payment_status,
			paid_at, paid_amount, payment_method, linked_contract_id, currency,
			payment_terms, late_fee, late_fee_rate
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
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
		return fmt.Errorf("failed to create invoice: %w", err)
	}

	return nil
}

func (r *Repository) GetInvoiceByID(ctx context.Context, id string) (*domain.Invoice, error) {
	query := `
		SELECT id, case_id, client_id, status, invoice_number, issue_date, due_date,
			   line_items, total_amount, tax_rate, tax_amount, notes, payment_status,
			   paid_at, paid_amount, payment_method, linked_contract_id, currency,
			   payment_terms, late_fee, late_fee_rate, created_at, updated_at
		FROM invoices
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	var invoice domain.Invoice
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

func (r *Repository) UpdateInvoice(ctx context.Context, invoice *domain.Invoice) error {
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

func (r *Repository) DeleteInvoice(ctx context.Context, id string) error {
	query := `DELETE FROM invoices WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete invoice: %w", err)
	}
	return nil
}

func (r *Repository) ListInvoicesByCase(ctx context.Context, caseID string, pagination domain.Pagination) ([]*domain.Invoice, int, error) {
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

	var invoices []*domain.Invoice
	for rows.Next() {
		var invoice domain.Invoice
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

// Document linking operations
func (r *Repository) GetLinkedDocuments(ctx context.Context, documentID string, docType domain.DocumentType) ([]domain.Documentable, error) {
	var linkedDocs []domain.Documentable

	switch docType {
	case domain.DocumentTypeEstimate:
		// Get mandates linked to this estimate
		query := `
			SELECT id, case_id, client_id, status, mandate_number, issue_date, scope_of_work,
				   valid_from, valid_until, terms_conditions, client_signature,
				   investigator_signature, linked_estimate_id, special_instructions, jurisdiction,
				   created_at, updated_at
			FROM mandates
			WHERE linked_estimate_id = $1
		`

		rows, err := r.pool.Query(ctx, query, documentID)
		if err != nil {
			return nil, fmt.Errorf("failed to query linked mandates: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var mandate domain.Mandate
			var clientSignatureJSON, investigatorSignatureJSON []byte

			err := rows.Scan(
				&mandate.ID, &mandate.CaseID, &mandate.ClientID, &mandate.Status,
				&mandate.MandateNumber, &mandate.IssueDate, &mandate.ScopeOfWork,
				&mandate.ValidFrom, &mandate.ValidUntil, &mandate.TermsConditions,
				&clientSignatureJSON, &investigatorSignatureJSON, &mandate.LinkedEstimateID,
				&mandate.SpecialInstructions, &mandate.Jurisdiction,
				&mandate.CreatedAt, &mandate.UpdatedAt,
			)

			if err != nil {
				return nil, fmt.Errorf("failed to scan linked mandate: %w", err)
			}

			if len(clientSignatureJSON) > 0 {
				if err := json.Unmarshal(clientSignatureJSON, &mandate.ClientSignature); err != nil {
					return nil, fmt.Errorf("failed to unmarshal client signature: %w", err)
				}
			}

			if len(investigatorSignatureJSON) > 0 {
				if err := json.Unmarshal(investigatorSignatureJSON, &mandate.InvestigatorSignature); err != nil {
					return nil, fmt.Errorf("failed to unmarshal investigator signature: %w", err)
				}
			}

			linkedDocs = append(linkedDocs, &mandate)
		}

	case domain.DocumentTypeMandate:
		// Get contracts linked to this mandate
		query := `
			SELECT id, case_id, client_id, status, contract_number, start_date, end_date,
				   scope_of_services, payment_terms, confidentiality, termination_clause,
				   signatures, linked_mandate_id, contract_value, currency, renewal_terms, governing_law,
				   created_at, updated_at
			FROM contracts
			WHERE linked_mandate_id = $1
		`

		rows, err := r.pool.Query(ctx, query, documentID)
		if err != nil {
			return nil, fmt.Errorf("failed to query linked contracts: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var contract domain.Contract
			var signaturesJSON []byte

			err := rows.Scan(
				&contract.ID, &contract.CaseID, &contract.ClientID, &contract.Status,
				&contract.ContractNumber, &contract.StartDate, &contract.EndDate,
				&contract.ScopeOfServices, &contract.PaymentTerms, &contract.Confidentiality,
				&contract.TerminationClause, &signaturesJSON, &contract.LinkedMandateID,
				&contract.ContractValue, &contract.Currency, &contract.RenewalTerms, &contract.GoverningLaw,
				&contract.CreatedAt, &contract.UpdatedAt,
			)

			if err != nil {
				return nil, fmt.Errorf("failed to scan linked contract: %w", err)
			}

			if err := json.Unmarshal(signaturesJSON, &contract.Signatures); err != nil {
				return nil, fmt.Errorf("failed to unmarshal signatures: %w", err)
			}

			linkedDocs = append(linkedDocs, &contract)
		}

	case domain.DocumentTypeContract:
		// Get invoices linked to this contract
		query := `
			SELECT id, case_id, client_id, status, invoice_number, issue_date, due_date,
				   line_items, total_amount, tax_rate, tax_amount, notes, payment_status,
				   paid_at, paid_amount, payment_method, linked_contract_id, currency,
				   payment_terms, late_fee, late_fee_rate, created_at, updated_at
			FROM invoices
			WHERE linked_contract_id = $1
		`

		rows, err := r.pool.Query(ctx, query, documentID)
		if err != nil {
			return nil, fmt.Errorf("failed to query linked invoices: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var invoice domain.Invoice
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
				return nil, fmt.Errorf("failed to scan linked invoice: %w", err)
			}

			if err := json.Unmarshal(lineItemsJSON, &invoice.LineItems); err != nil {
				return nil, fmt.Errorf("failed to unmarshal line items: %w", err)
			}

			linkedDocs = append(linkedDocs, &invoice)
		}

	default:
		return nil, fmt.Errorf("document type %s does not support linking", docType)
	}

	return linkedDocs, nil
}