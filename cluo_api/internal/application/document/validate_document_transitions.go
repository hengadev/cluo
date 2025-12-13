package document

import (
	"context"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// ValidateDocumentTransitions validates if a document can transition to a new status.
func (s *Service) ValidateDocumentTransitions(ctx context.Context, id string, docType document.DocumentType, newStatus document.DocumentStatus) error {
	// Get document
	doc, err := s.repo.GetByID(ctx, id, docType)
	if err != nil {
		return errs.NewNotFoundErr(err, "document")
	}

	// Validate transition
	return s.validateDocumentTransition(doc, newStatus)
}