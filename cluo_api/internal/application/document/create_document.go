package document

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// CreateDocument creates a new document of the specified type.
func (s *Service) CreateDocument(ctx context.Context, req *document.CreateDocumentRequest) (document.Documentable, error) {
	// Validate request
	if err := req.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Verify case exists
	caseExists, err := s.caseRepo.ExistsCase(ctx, req.CaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to check case existence: %w", err)
	}
	if !caseExists {
		return nil, errs.NewNotFoundErr(fmt.Errorf("case not found"), "case")
	}

	// Verify client exists
	clientExists, err := s.clientRepo.ExistsClient(ctx, req.ClientID)
	if err != nil {
		return nil, fmt.Errorf("failed to check client existence: %w", err)
	}
	if !clientExists {
		return nil, errs.NewNotFoundErr(fmt.Errorf("client not found"), "client")
	}

	var doc document.Documentable

	switch req.Type {
	case document.DocumentTypeEstimate:
		estimate, ok := req.Data.(*document.Estimate)
		if !ok {
			return nil, errs.NewInvalidValueErr("invalid data type for estimate")
		}
		doc = estimate

	case document.DocumentTypeMandate:
		mandate, ok := req.Data.(*document.Mandate)
		if !ok {
			return nil, errs.NewInvalidValueErr("invalid data type for mandate")
		}
		doc = mandate

	case document.DocumentTypeContract:
		contract, ok := req.Data.(*document.Contract)
		if !ok {
			return nil, errs.NewInvalidValueErr("invalid data type for contract")
		}
		doc = contract

	case document.DocumentTypeInvoice:
		invoice, ok := req.Data.(*document.Invoice)
		if !ok {
			return nil, errs.NewInvalidValueErr("invalid data type for invoice")
		}
		doc = invoice

	default:
		return nil, errs.NewInvalidValueErr(fmt.Sprintf("unsupported document type: %s", req.Type))
	}

	// Set base fields
	doc.SetCaseID(req.CaseID)
	doc.SetClientID(req.ClientID)
	if req.Status != nil {
		doc.SetStatus(*req.Status)
	}

	// Validate document
	if err := doc.Validate(); err != nil {
		return nil, errs.NewInvalidValueErr(fmt.Sprintf("document validation failed: %s", err.Error()))
	}

	// Encrypt and save document
	if err := s.repo.Create(ctx, doc); err != nil {
		return nil, errs.NewNotCreatedErr(err, "document")
	}

	// Create initial version
	if err := s.createDocumentVersion(ctx, doc, nil, stringPtr("Initial creation")); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return doc, nil
}