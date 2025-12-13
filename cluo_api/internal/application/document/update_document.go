package document

import (
	"context"
	"fmt"
	"time"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// UpdateDocument updates an existing document.
func (s *Service) UpdateDocument(ctx context.Context, id string, req *document.UpdateDocumentRequest) (document.Documentable, error) {
	// Validate request
	if err := req.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Get existing document
	doc, err := s.repo.GetByID(ctx, id, req.Type)
	if err != nil {
		return nil, errs.NewNotFoundErr(err, "document")
	}

	// Check if document can be modified
	if doc.GetStatus().IsFinal() {
		return nil, errs.NewInvalidValueErr(fmt.Sprintf("cannot modify document in final status: %s", doc.GetStatus()))
	}

	// Apply updates
	switch d := doc.(type) {
	case *document.Estimate:
		if req.Title != nil {
			d.Title = *req.Title
		}
		if req.Description != nil {
			d.Description = *req.Description
		}
		// Apply other field updates...

	case *document.Mandate:
		if req.Title != nil {
			d.Title = *req.Title
		}
		if req.Description != nil {
			d.Description = *req.Description
		}
		// Apply other field updates...

	case *document.Contract:
		if req.Title != nil {
			d.Title = *req.Title
		}
		if req.Description != nil {
			d.Description = *req.Description
		}
		// Apply other field updates...

	case *document.Invoice:
		if req.Title != nil {
			d.Title = *req.Title
		}
		if req.Description != nil {
			d.Description = *req.Description
		}
		// Apply other field updates...
	}

	// Update metadata
	doc.SetUpdatedAt(time.Now())

	// Validate updated document
	if err := doc.Validate(); err != nil {
		return nil, errs.NewInvalidValueErr(fmt.Sprintf("document validation failed: %s", err.Error()))
	}

	// Save updates
	if err := s.repo.Update(ctx, doc); err != nil {
		return nil, errs.NewNotUpdatedErr(err, "document")
	}

	// Create version record
	authorID := s.getUserIDFromContext(ctx)
	if err := s.createDocumentVersion(ctx, doc, &authorID, req.Reason); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return doc, nil
}