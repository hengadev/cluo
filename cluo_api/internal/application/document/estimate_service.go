package document

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/hengadev/cluo_api/internal/domain"
	"github.com/hengadev/cluo_api/internal/ports"
)

// Estimate operations

func (s *Service) CreateEstimate(ctx context.Context, estimate *domain.Estimate) (*domain.Estimate, error) {
	// TODO: Validate user permissions to create estimates
	// TODO: Verify that case and client exist and are accessible

	// Validate estimate
	if err := estimate.Validate(); err != nil {
		return nil, fmt.Errorf("estimate validation failed: %w", err)
	}

	// Generate estimate number if not provided
	if estimate.EstimateNumber == "" {
		// TODO: Implement estimate number generation
		estimate.EstimateNumber = fmt.Sprintf("EST-%d", uuid.New().ID())
	}

	// Save to repository
	if err := s.repo.CreateEstimate(ctx, estimate); err != nil {
		return nil, fmt.Errorf("failed to create estimate: %w", err)
	}

	// Create initial version
	if err := s.createDocumentVersion(ctx, estimate, nil, stringPtr("Initial creation")); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return estimate, nil
}

func (s *Service) AcceptEstimate(ctx context.Context, estimateID string, acceptedBy string) (*domain.Mandate, error) {
	// TODO: Validate user permissions to accept estimates
	// TODO: Verify that acceptedBy is a valid user with appropriate permissions

	// Get estimate
	estimate, err := s.repo.GetEstimateByID(ctx, estimateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get estimate: %w", err)
	}

	// Check if estimate can be accepted
	if !estimate.CanBeAccepted() {
		return nil, fmt.Errorf("estimate cannot be accepted")
	}

	// Accept the estimate
	acceptedByUUID := uuid.MustParse(acceptedBy) // TODO: Handle parsing error properly
	if err := estimate.Accept(acceptedByUUID); err != nil {
		return nil, fmt.Errorf("failed to accept estimate: %w", err)
	}

	// Update estimate
	if err := s.repo.UpdateEstimate(ctx, estimate); err != nil {
		return nil, fmt.Errorf("failed to update estimate: %w", err)
	}

	// Create version record
	if err := s.createDocumentVersion(ctx, estimate, &acceptedByUUID, stringPtr("Estimate accepted")); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	// Create mandate from estimate
	mandate := &domain.Mandate{
		DocumentBase:     domain.NewDocumentBase(estimate.CaseID, estimate.ClientID),
		MandateNumber:    fmt.Sprintf("MND-%d", uuid.New().ID()), // TODO: Generate proper mandate number
		IssueDate:        estimate.IssueDate,
		ScopeOfWork:      "Investigation services as outlined in estimate " + estimate.EstimateNumber,
		ValidFrom:        estimate.IssueDate,
		ValidUntil:       estimate.ValidUntil,
		TermsConditions:  "Standard terms and conditions apply as per estimate agreement",
		LinkedEstimateID: &estimate.ID,
	}

	// Validate mandate
	if err := mandate.Validate(); err != nil {
		return nil, fmt.Errorf("created mandate validation failed: %w", err)
	}

	// Save mandate
	if err := s.repo.CreateMandate(ctx, mandate); err != nil {
		return nil, fmt.Errorf("failed to create mandate: %w", err)
	}

	// Create version record for mandate
	if err := s.createDocumentVersion(ctx, mandate, &acceptedByUUID, stringPtr("Created from accepted estimate")); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return mandate, nil
}

func (s *Service) UpdateEstimate(ctx context.Context, estimateID string, lineItems []domain.EstimateItem) (*domain.Estimate, error) {
	// TODO: Validate user permissions to update estimates

	// Get estimate
	estimate, err := s.repo.GetEstimateByID(ctx, estimateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get estimate: %w", err)
	}

	// Check if estimate can be modified
	if estimate.Accepted {
		return nil, fmt.Errorf("cannot modify accepted estimate")
	}

	// Update line items
	estimate.LineItems = lineItems
	estimate.CalculateTotal()

	// Validate updated estimate
	if err := estimate.Validate(); err != nil {
		return nil, fmt.Errorf("updated estimate validation failed: %w", err)
	}

	// Update estimate
	if err := s.repo.UpdateEstimate(ctx, estimate); err != nil {
		return nil, fmt.Errorf("failed to update estimate: %w", err)
	}

	// Create version record
	// TODO: Get author ID from context
	if err := s.createDocumentVersion(ctx, estimate, nil, stringPtr("Updated line items")); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return estimate, nil
}