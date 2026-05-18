package investigationService

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/investigation"
)

// ListByClient implements the ListByClient function for a specific client
func (s *CaseService) ListByClient(ctx context.Context, r *investigation.ListByClientRequest) (*investigation.ListCasesResponse, error) {
	// Validate request
	if err := r.Valid(ctx); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Parse client ID
	clientID, err := uuid.Parse(r.ClientID)
	if err != nil {
		return nil, fmt.Errorf("invalid client ID: %w", err)
	}

	// Create pagination object
	pagination := investigation.Pagination{
		Page:     r.Page,
		PageSize: r.PageSize,
	}

	// Get cases for the specific client from repository
	caseEncxList, total, err := s.repo.ListByClient(ctx, clientID, pagination)
	if err != nil {
		return nil, fmt.Errorf("failed to list cases for client: %w", err)
	}

	// Decrypt and convert to response
	caseResponses := make([]*investigation.CaseResponse, 0, len(caseEncxList))
	for _, caseEncx := range caseEncxList {
		caseDecrypted, err := investigation.DecryptInvestigationEncx(ctx, s.crypto, caseEncx)
		if err != nil {
			// Skip cases that can't be decrypted
			continue
		}

		caseResponses = append(caseResponses, caseDecrypted.ToResponse())
	}

	// Create pagination info
	paginationInfo := investigation.NewPaginationInfo(r.Page, r.PageSize, total)

	return &investigation.ListCasesResponse{
		Cases:      caseResponses,
		Pagination: paginationInfo,
	}, nil
}