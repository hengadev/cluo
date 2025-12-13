package document

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// GetDocumentWorkflow retrieves all documents in a case workflow.
func (s *Service) GetDocumentWorkflow(ctx context.Context, caseID string) ([]document.DocumentSummary, error) {
	// Verify case exists
	caseExists, err := s.caseRepo.ExistsCase(ctx, caseID)
	if err != nil {
		return nil, fmt.Errorf("failed to check case existence: %w", err)
	}
	if !caseExists {
		return nil, fmt.Errorf("case not found")
	}

	// Get all documents for the case
	filter := document.DocumentFilter{
		CaseID: &caseID,
	}

	pagination := document.NewPagination()
	pagination.PageSize = 100 // Set a reasonable limit

	documents, _, err := s.repo.List(ctx, filter, pagination)
	if err != nil {
		return nil, fmt.Errorf("failed to list case documents: %w", err)
	}

	return documents, nil
}