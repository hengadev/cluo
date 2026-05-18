package investigationService

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/investigation"
)

func (s *CaseService) CreateCase(ctx context.Context, r *investigation.CreateCaseRequest) (*investigation.CaseResponse, error) {
	if err := r.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Check if client exists in database
	clientUUID, err := uuid.Parse(r.ClientID)
	if err != nil {
		return nil, errs.NewInvalidValueErr("invalid client ID format")
	}

	exists, err := s.clientRepo.ExistsClient(ctx, clientUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to check client existence: %w", err)
	}

	if !exists {
		return nil, errs.NewRepositoryNotFoundErr(fmt.Errorf("client with ID %s not found", r.ClientID), "client")
	}

	// Check if assigned contact exists (if provided)
	if r.AssignedContactID != nil {
		contactUUID, err := uuid.Parse(*r.AssignedContactID)
		if err != nil {
			return nil, errs.NewInvalidValueErr("invalid assigned contact ID format")
		}

		exists, err = s.clientRepo.ExistsContact(ctx, contactUUID)
		if err != nil {
			return nil, fmt.Errorf("failed to check contact existence: %w", err)
		}

		if !exists {
			return nil, errs.NewRepositoryNotFoundErr(fmt.Errorf("contact with ID %s not found", *r.AssignedContactID), "contact")
		}
	}

	// Check if case subject exists (if provided)
	if r.CaseSubjectID != nil {
		caseSubjectUUID, err := uuid.Parse(*r.CaseSubjectID)
		if err != nil {
			return nil, errs.NewInvalidValueErr("invalid case subject ID format")
		}

		exists, err = s.caseSubjectRepo.ExistsCaseSubject(ctx, caseSubjectUUID)
		if err != nil {
			return nil, fmt.Errorf("failed to check case subject existence: %w", err)
		}

		if !exists {
			return nil, errs.NewRepositoryNotFoundErr(fmt.Errorf("case subject with ID %s not found", *r.CaseSubjectID), "case_subject")
		}
	}

	c := investigation.New(r)
	cEncx, err := investigation.ProcessInvestigationEncx(ctx, s.crypto, c)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("case", err)
	}

	if err := s.repo.CreateCase(ctx, cEncx); err != nil {
		return nil, fmt.Errorf("failed to create case: %w", err)
	}
	return c.ToResponse(), nil
}
