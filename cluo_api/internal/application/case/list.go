package caseService

import (
	"context"
	"fmt"
	"strings"

	caseDomain "github.com/hengadev/cluo_api/internal/domain/case"
)

// List implements the general List function with filtering and decryption
func (s *CaseService) List(ctx context.Context, r *caseDomain.ListCasesRequest) (*caseDomain.ListCasesResponse, error) {
	// Validate request
	if err := r.Valid(ctx); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Convert request to CaseFilter for repository
	caseFilter, err := r.ToCaseFilter()
	if err != nil {
		return nil, fmt.Errorf("failed to parse filters: %w", err)
	}

	// Create pagination object
	pagination := caseDomain.Pagination{
		Page:     r.Page,
		PageSize: r.PageSize,
	}

	// Get cases from repository
	caseEncxList, total, err := s.repo.List(ctx, caseFilter, pagination)
	if err != nil {
		return nil, fmt.Errorf("failed to list cases: %w", err)
	}

	// Decrypt and convert to response
	caseResponses := make([]*caseDomain.CaseResponse, 0, len(caseEncxList))
	for _, caseEncx := range caseEncxList {
		caseDecrypted, err := caseDomain.DecryptCaseEncx(ctx, s.crypto, caseEncx)
		if err != nil {
			// Skip cases that can't be decrypted
			continue
		}

		// Apply service-layer filtering for encrypted fields
		if !caseMatchesFilter(caseDecrypted, caseFilter) {
			continue
		}

		caseResponses = append(caseResponses, caseDecrypted.ToResponse())
	}

	// Adjust total count for cases filtered at service layer
	// (this is approximate since we can't get the exact count without decryption)
	if caseFilter.Status != nil || caseFilter.Search != nil ||
		caseFilter.DateUpdatedFrom != nil || caseFilter.DateUpdatedTo != nil {
		total = len(caseResponses)
	}

	// Create pagination info
	paginationInfo := caseDomain.NewPaginationInfo(r.Page, r.PageSize, total)

	return &caseDomain.ListCasesResponse{
		Cases:      caseResponses,
		Pagination: paginationInfo,
	}, nil
}

// caseMatchesFilter applies service-layer filtering for encrypted fields
func caseMatchesFilter(caseDecrypted *caseDomain.Case, filter caseDomain.CaseFilter) bool {
	// Filter by status
	if filter.Status != nil && caseDecrypted.Status != *filter.Status {
		return false
	}

	// Filter by updated date range
	if filter.DateUpdatedFrom != nil && caseDecrypted.UpdatedAt.Before(*filter.DateUpdatedFrom) {
		return false
	}
	if filter.DateUpdatedTo != nil && caseDecrypted.UpdatedAt.After(*filter.DateUpdatedTo) {
		return false
	}

	// Filter by search term (title and description)
	if filter.Search != nil {
		searchTerm := *filter.Search
		titleMatch := len(searchTerm) > 0 && len(searchTerm) <= 1000 &&
			caseDecrypted.Title != "" &&
			containsIgnoreCase(caseDecrypted.Title, searchTerm)
		descMatch := len(searchTerm) > 0 && len(searchTerm) <= 1000 &&
			caseDecrypted.Description != "" &&
			containsIgnoreCase(caseDecrypted.Description, searchTerm)

		if !titleMatch && !descMatch {
			return false
		}
	}

	return true
}

// containsIgnoreCase performs case-insensitive string search
func containsIgnoreCase(s, substr string) bool {
	if len(s) < len(substr) || substr == "" {
		return false
	}

	// Convert to lowercase for case-insensitive comparison
	sLower := strings.ToLower(s)
	substrLower := strings.ToLower(substr)

	return strings.Contains(sLower, substrLower)
}
