package document

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// ListDocuments lists documents with filtering and pagination.
func (s *Service) ListDocuments(ctx context.Context, filter document.DocumentFilter, pagination document.Pagination) ([]document.DocumentSummary, int, error) {
	if err := pagination.Validate(); err != nil {
		return nil, 0, errs.NewInvalidValueErr(fmt.Sprintf("invalid pagination: %s", err.Error()))
	}

	documents, total, err := s.repo.List(ctx, filter, pagination)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list documents: %w", err)
	}

	return documents, total, nil
}