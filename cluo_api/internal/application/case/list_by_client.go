package caseService

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	caseDomain "github.com/hengadev/cluo_api/internal/domain/case"
)

// ListByClient implements the ListByClient function for a specific client
func (s *CaseService) ListByClient(ctx context.Context, r *caseDomain.ListByClientRequest) (*caseDomain.ListCasesResponse, error) {
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
	pagination := caseDomain.Pagination{
		Page:     r.Page,
		PageSize: r.PageSize,
	}

	// Get cases for the specific client from repository
	caseEncxList, total, err := s.repo.ListByClient(ctx, clientID, pagination)
	if err != nil {
		return nil, fmt.Errorf("failed to list cases for client: %w", err)
	}

	// Decrypt and convert to response
	caseResponses := make([]*caseDomain.CaseResponse, 0, len(caseEncxList))
	for _, caseEncx := range caseEncxList {
		caseDecrypted, err := caseDomain.DecryptCaseEncx(ctx, s.crypto, caseEncx)
		if err != nil {
			// Skip cases that can't be decrypted
			continue
		}

		caseResponses = append(caseResponses, caseDecrypted.ToResponse())
	}

	// Create pagination info
	paginationInfo := caseDomain.NewPaginationInfo(r.Page, r.PageSize, total)

	return &caseDomain.ListCasesResponse{
		Cases:      caseResponses,
		Pagination: paginationInfo,
	}, nil
}