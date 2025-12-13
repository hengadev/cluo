package document

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// GetDocumentHistory retrieves the version history of a document.
func (s *Service) GetDocumentHistory(ctx context.Context, id string, docType document.DocumentType, pagination document.Pagination) ([]*document.DocumentVersion, int, error) {
	if err := pagination.Validate(); err != nil {
		return nil, 0, errs.NewInvalidValueErr(fmt.Sprintf("invalid pagination: %s", err.Error()))
	}

	versions, total, err := s.versionRepo.GetDocumentHistory(ctx, id, docType, pagination)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get document history: %w", err)
	}

	return versions, total, nil
}