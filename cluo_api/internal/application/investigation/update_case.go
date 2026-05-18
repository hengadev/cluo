package investigationService

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/investigation"
)

func (s *CaseService) UpdateCase(ctx context.Context, request *investigation.UpdateCaseRequest) (*investigation.CaseResponse, error) {
	if err := request.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Get existing case from repository
	caseEncx, err := s.repo.GetCaseByID(ctx, request.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get case by ID: %w", err)
	}

	// Decrypt case data to allow field updates
	caseDecrypted, err := investigation.DecryptInvestigationEncx(ctx, s.crypto, caseEncx)
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

	// Update CaseSubjectID with validation (if provided)
	if request.CaseSubjectID != nil {
		// Validate existence if not nil
		if *request.CaseSubjectID != uuid.Nil {
			exists, err := s.caseSubjectRepo.ExistsCaseSubject(ctx, *request.CaseSubjectID)
			if err != nil {
				return nil, fmt.Errorf("failed to check case subject existence: %w", err)
			}
			if !exists {
				return nil, errs.NewRepositoryNotFoundErr(fmt.Errorf("case subject not found"), "case_subject")
			}
		}
		caseDecrypted.CaseSubjectID = request.CaseSubjectID
	}

	// Update location fields (all optional)
	if request.Placename != nil {
		caseDecrypted.Placename = *request.Placename
	}

	if request.Address1 != nil {
		caseDecrypted.Address1 = *request.Address1
	}

	if request.Address2 != nil {
		caseDecrypted.Address2 = *request.Address2
	}

	if request.City != nil {
		caseDecrypted.City = *request.City
	}

	if request.PostalCode != nil {
		caseDecrypted.PostalCode = *request.PostalCode
	}

	if request.Country != nil {
		caseDecrypted.Country = *request.Country
	}

	if request.Latitude != nil {
		caseDecrypted.Latitude = request.Latitude
	}

	if request.Longitude != nil {
		caseDecrypted.Longitude = request.Longitude
	}

	if request.LocationType != nil {
		caseDecrypted.LocationType = *request.LocationType
	}

	if request.LocationNotes != nil {
		caseDecrypted.LocationNotes = *request.LocationNotes
	}

	if request.ExternalReference != nil {
		caseDecrypted.ExternalReference = request.ExternalReference
	}

	if request.CaseType != nil {
		caseDecrypted.CaseType = *request.CaseType
	}

	if request.Status != nil {
		// Parse status, defaulting to Draft if invalid
		status := investigation.Status(strings.ToLower(strings.TrimSpace(*request.Status)))
		if status.IsValid() {
			caseDecrypted.Status = status
		} else {
			caseDecrypted.Status = investigation.StatusDraft
		}
	}

	// Update the timestamp for the update
	caseDecrypted.UpdatedAt = time.Now()

	// Encrypt the case data
	updatedCaseEncx, err := investigation.ProcessInvestigationEncx(ctx, s.crypto, caseDecrypted)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("case for update", err)
	}

	// Save updated case to repository
	if err := s.repo.UpdateCase(ctx, updatedCaseEncx); err != nil {
		return nil, fmt.Errorf("failed to update case: %w", err)
	}

	return caseDecrypted.ToResponse(), nil
}

