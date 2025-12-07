package caseService

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hengadev/cluo_api/internal/common/errs"
	caseDomain "github.com/hengadev/cluo_api/internal/domain/case"
)

func (s *CaseService) UpdateCase(ctx context.Context, request *caseDomain.UpdateCaseRequest) (*caseDomain.CaseResponse, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Get existing case from repository
	caseEncx, err := s.repo.GetCaseByID(ctx, request.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get case by ID: %w", err)
	}

	// Decrypt case data to allow field updates
	caseDecrypted, err := caseDomain.DecryptCaseEncx(ctx, s.crypto, caseEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("case for update", err)
	}

	// Update only non-nil fields from request
	if request.Title != nil {
		caseDecrypted.Title = *request.Title
	}

	if request.Description != nil {
		caseDecrypted.Description = *request.Description
	}

	if request.ClientID != nil {
		caseDecrypted.ClientID = *request.ClientID
	}

	if request.AssignedContactID != nil {
		caseDecrypted.AssignedContactID = request.AssignedContactID
	}

	if request.Status != nil {
		// Parse status, defaulting to Draft if invalid
		status := caseDomain.CaseStatus(strings.ToLower(strings.TrimSpace(*request.Status)))
		if status.IsValid() {
			caseDecrypted.Status = status
		} else {
			caseDecrypted.Status = caseDomain.CaseStatusDraft
		}
	}

	// Update the timestamp for the update
	caseDecrypted.UpdatedAt = time.Now()

	// Encrypt the case data
	updatedCaseEncx, err := caseDomain.ProcessCaseEncx(ctx, s.crypto, caseDecrypted)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("case for update", err)
	}

	// Save updated case to repository
	if err := s.repo.UpdateCase(ctx, updatedCaseEncx); err != nil {
		return nil, fmt.Errorf("failed to update case: %w", err)
	}

	return caseDecrypted.ToResponse(), nil
}

