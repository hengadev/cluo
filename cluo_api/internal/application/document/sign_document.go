package document

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// SignDocument signs a document.
func (s *Service) SignDocument(ctx context.Context, id string, docType document.DocumentType, req *document.SignDocumentRequest) error {
	// Validate request
	if err := req.Valid(ctx); err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}

	// Get document
	doc, err := s.repo.GetByID(ctx, id, docType)
	if err != nil {
		return errs.NewNotFoundErr(err, "document")
	}

	// Create signature
	signature := &document.Signature{
		ID:        uuid.New(),
		Name:      req.SignerName,
		Role:      req.SignerRole,
		SignedAt:  time.Now(),
		Method:    req.Method,
		IPAddress: req.IPAddress,
		UserAgent: req.UserAgent,
	}

	// Apply signature based on document type
	var updated bool
	switch d := doc.(type) {
	case *document.Mandate:
		if req.SignerRole == "client" && d.ClientSignature == nil {
			d.AddClientSignature(signature)
			updated = true
		} else if req.SignerRole == "investigator" && d.InvestigatorSignature == nil {
			d.AddInvestigatorSignature(signature)
			updated = true
		}

	case *document.Contract:
		d.AddSignature(signature)
		updated = true

	default:
		return errs.NewInvalidValueErr("signing not supported for this document type")
	}

	if !updated {
		return errs.NewInvalidValueErr("could not apply signature - role may already be signed")
	}

	// Check if document should transition to signed status
	if err := s.validateDocumentTransition(doc, document.DocumentStatusSigned); err == nil {
		doc.SetStatus(document.DocumentStatusSigned)
	}

	doc.SetUpdatedAt(time.Now())

	// Save updates
	if err := s.repo.Update(ctx, doc); err != nil {
		return errs.NewNotUpdatedErr(err, "document")
	}

	// Create version record
	if err := s.createDocumentVersion(ctx, doc, nil, stringPtr("Document signed")); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return nil
}