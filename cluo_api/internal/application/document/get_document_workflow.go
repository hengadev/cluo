package document

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// GetDocumentWorkflow retrieves the full financial document chain for a case.
// It returns a structured response with Estimate, Mandate, Contract, and Invoice.
// Documents that don't exist yet are returned as nil.
func (s *Service) GetDocumentWorkflow(ctx context.Context, caseID string) (*document.DocumentWorkflowResponse, error) {
	parsedCaseID, err := uuid.Parse(caseID)
	if err != nil {
		return nil, fmt.Errorf("invalid case ID: %w", err)
	}

	// Verify case exists
	caseExists, err := s.caseRepo.ExistsCase(ctx, parsedCaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to check case existence: %w", err)
	}
	if !caseExists {
		return nil, errors.New("case not found")
	}

	resp := &document.DocumentWorkflowResponse{}

	// Fetch estimate
	estRaw, err := s.repo.GetFirstByCaseAndType(ctx, caseID, document.DocumentTypeEstimate)
	if err != nil && !errors.Is(err, errs.ErrRepositoryNotFound) {
		return nil, fmt.Errorf("failed to fetch estimate: %w", err)
	}
	if err == nil && estRaw != nil {
		encx, ok := estRaw.(*document.EstimateEncx)
		if !ok {
			return nil, fmt.Errorf("unexpected type for estimate document")
		}
		estimate, decErr := document.DecryptEstimateEncx(ctx, s.crypto, encx)
		if decErr != nil {
			return nil, fmt.Errorf("failed to decrypt estimate: %w", decErr)
		}
		resp.Estimate = estimate
	}

	// Fetch mandate
	mandRaw, err := s.repo.GetFirstByCaseAndType(ctx, caseID, document.DocumentTypeMandate)
	if err != nil && !errors.Is(err, errs.ErrRepositoryNotFound) {
		return nil, fmt.Errorf("failed to fetch mandate: %w", err)
	}
	if err == nil && mandRaw != nil {
		encx, ok := mandRaw.(*document.MandateEncx)
		if !ok {
			return nil, fmt.Errorf("unexpected type for mandate document")
		}
		mandate, decErr := document.DecryptMandateEncx(ctx, s.crypto, encx)
		if decErr != nil {
			return nil, fmt.Errorf("failed to decrypt mandate: %w", decErr)
		}
		resp.Mandate = mandate
	}

	// Fetch contract
	contRaw, err := s.repo.GetFirstByCaseAndType(ctx, caseID, document.DocumentTypeContract)
	if err != nil && !errors.Is(err, errs.ErrRepositoryNotFound) {
		return nil, fmt.Errorf("failed to fetch contract: %w", err)
	}
	if err == nil && contRaw != nil {
		encx, ok := contRaw.(*document.ContractEncx)
		if !ok {
			return nil, fmt.Errorf("unexpected type for contract document")
		}
		contract, decErr := document.DecryptContractEncx(ctx, s.crypto, encx)
		if decErr != nil {
			return nil, fmt.Errorf("failed to decrypt contract: %w", decErr)
		}
		resp.Contract = contract
	}

	// Fetch invoice
	invRaw, err := s.repo.GetFirstByCaseAndType(ctx, caseID, document.DocumentTypeInvoice)
	if err != nil && !errors.Is(err, errs.ErrRepositoryNotFound) {
		return nil, fmt.Errorf("failed to fetch invoice: %w", err)
	}
	if err == nil && invRaw != nil {
		encx, ok := invRaw.(*document.InvoiceEncx)
		if !ok {
			return nil, fmt.Errorf("unexpected type for invoice document")
		}
		invoice, decErr := document.DecryptInvoiceEncx(ctx, s.crypto, encx)
		if decErr != nil {
			return nil, fmt.Errorf("failed to decrypt invoice: %w", decErr)
		}
		resp.Invoice = invoice
	}

	return resp, nil
}
