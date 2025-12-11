package documentRepository

import (
	"context"
	"fmt"
	"strings"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// List retrieves documents with optional filtering and pagination.
func (r *Repository) List(ctx context.Context, filter document.DocumentFilter, pagination document.Pagination) ([]document.DocumentSummary, int, error) {
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

	var documents []document.DocumentSummary
	for rows.Next() {
		var doc document.DocumentSummary
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