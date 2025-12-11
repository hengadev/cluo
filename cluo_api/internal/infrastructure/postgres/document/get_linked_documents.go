package documentRepository

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

	var linkedDocs []document.Documentable

	switch docType {
	case document.DocumentTypeEstimate:
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
			var mandate document.Mandate
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

	case document.DocumentTypeMandate:
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
			var contract document.Contract
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

	case document.DocumentTypeContract:
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
