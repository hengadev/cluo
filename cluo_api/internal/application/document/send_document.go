package document

import (
	"context"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// SendDocument sends a document to recipients.
func (s *Service) SendDocument(ctx context.Context, id string, docType document.DocumentType, req *document.SendDocumentRequest) error {
	// Validate request
	if err := req.Valid(ctx); err != nil {
		return errs.NewInvalidValueErr(err.Error())
	}

	// Get document
	doc, err := s.repo.GetByID(ctx, id, docType)
	if err != nil {
		return errs.NewNotFoundErr(err, "document")
	}

	// Check if document can be sent
	if doc.GetStatus() != document.DocumentStatusDraft {
		return errs.NewInvalidValueErr("only draft documents can be sent")
	}

	// TODO: Implement actual sending logic:
	// - Generate PDF
	// - Send emails
	// - Send SMS notifications
	// - Create notification events

	// Update document status to Sent
	doc.SetStatus(document.DocumentStatusSent)
	doc.UpdateTimestamp()

	if err := s.repo.Update(ctx, doc); err != nil {
		return errs.NewNotUpdatedErr(err, "document")
	}

	// Create version record
	authorID := s.getUserIDFromContext(ctx)
	if err := s.createDocumentVersion(ctx, doc, &authorID, stringPtr("Document sent")); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return nil
}