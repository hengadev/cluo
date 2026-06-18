package document

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// DeleteDocument deletes a document by ID.
func (s *Service) DeleteDocument(ctx context.Context, id string, docType document.DocumentType) error {
	// Get document to check status
	doc, err := s.repo.GetByID(ctx, id, docType)
	if err != nil {
		return errs.NewNotFoundErr(err, "document")
	}

	if doc.GetStatus() != document.DocumentStatusDraft {
		return errs.NewInvalidValueErr(fmt.Sprintf("cannot delete document in status: %s", doc.GetStatus()))
	}

	// Check for linked documents
	linkedDocs, err := s.repo.GetLinkedDocuments(ctx, id, docType)
	if err != nil {
		return fmt.Errorf("failed to check linked documents: %w", err)
	}

	if len(linkedDocs) > 0 {
		return errs.NewConflictErr(fmt.Errorf("cannot delete document with linked documents"))
	}

	// Delete document
	if err := s.repo.Delete(ctx, id, docType); err != nil {
		return errs.NewNotDeletedErr(err, "document")
	}

	// Delete version history
	if err := s.versionRepo.DeleteVersions(ctx, id, docType); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return nil
}