package caseSubjectService

import (
	"context"
	"fmt"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/subject"
)

func (s *Service) UpdateCaseSubject(ctx context.Context, req *subject.UpdateCaseSubjectRequest) (*subject.CaseSubjectResponse, error) {
	encxVal, err := s.repo.GetCaseSubjectByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("get case subject for update: %w", err)
	}

	plain, err := subject.DecryptCaseSubjectEncx(ctx, s.crypto, encxVal)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("case subject for update", err)
	}

	if req.Lastname != nil {
		plain.Lastname = *req.Lastname
	}
	if req.Firstname != nil {
		plain.Firstname = *req.Firstname
	}
	if req.Email != nil {
		plain.Email = *req.Email
	}
	if req.Phone != nil {
		plain.Phone = *req.Phone
	}
	if req.City != nil {
		plain.City = *req.City
	}
	if req.PostalCode != nil {
		plain.PostalCode = *req.PostalCode
	}
	if req.Address1 != nil {
		plain.Address1 = *req.Address1
	}
	if req.Address2 != nil {
		plain.Address2 = *req.Address2
	}
	if req.Occupation != nil {
		plain.Occupation = *req.Occupation
	}
	if req.Notes != nil {
		plain.Notes = *req.Notes
	}

	updatedEncx, err := subject.ProcessCaseSubjectEncx(ctx, s.crypto, plain)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("case subject for update", err)
	}

	if err := s.repo.UpdateCaseSubject(ctx, updatedEncx); err != nil {
		return nil, fmt.Errorf("update case subject: %w", err)
	}

	return plain.ToResponse(), nil
}
