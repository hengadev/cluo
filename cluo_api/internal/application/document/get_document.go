package document

import (
	"context"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// GetDocument retrieves a document by ID.
func (s *Service) GetDocument(ctx context.Context, id string, docType document.DocumentType) (document.Documentable, error) {
	doc, err := s.repo.GetByID(ctx, id, docType)
	if err != nil {
		return nil, errs.NewNotFoundErr(err, "document")
	}

	return doc, nil
}