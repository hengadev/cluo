package document

import (
	"context"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// ArchiveDocument archives a document.
func (s *Service) ArchiveDocument(ctx context.Context, id string, docType document.DocumentType) error {
	// Get document
	doc, err := s.repo.GetByID(ctx, id, docType)
	if err != nil {
		return errs.NewNotFoundErr(err, "document")
	}

	// Check if document can be archived
	if doc.GetStatus() == document.DocumentStatusArchived {
		return errs.NewInvalidValueErr("document is already archived")
	}

	// Validate transition to archived
	if err := s.validateDocumentTransition(doc, document.DocumentStatusArchived); err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}

	// Update status
	doc.SetStatus(document.DocumentStatusArchived)

	// Save updates
	if err := s.repo.Update(ctx, doc); err != nil {
		return errs.NewNotUpdatedErr(err, "document")
	}

	// Create version record
	authorID := s.getUserIDFromContext(ctx)
	if err := s.createDocumentVersion(ctx, doc, &authorID, stringPtr("Document archived")); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return nil
}

