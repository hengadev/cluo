package caseSubjectService

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/subject"
)

func (s *Service) CreateCaseSubject(ctx context.Context, req *subject.CreateCaseSubjectRequest) (*subject.CaseSubjectResponse, error) {
	if req.Lastname == "" || req.Firstname == "" {
		return nil, errs.NewInvalidValueErr("lastname and firstname are required")
	}

	plain := &subject.Subject{
		ID:         uuid.New(),
		Lastname:   req.Lastname,
		Firstname:  req.Firstname,
		Email:      req.Email,
		Phone:      req.Phone,
		City:       req.City,
		PostalCode: req.PostalCode,
		Address1:   req.Address1,
		Address2:   req.Address2,
		Occupation: req.Occupation,
		Notes:      req.Notes,
		CreatedAt:  time.Now(),
	}

	encxVal, err := subject.ProcessCaseSubjectEncx(ctx, s.crypto, plain)
	if err != nil {
		return nil, errs.NewNotEncryptedErr("case subject", err)
	}

	if err := s.repo.CreateCaseSubject(ctx, encxVal); err != nil {
		return nil, fmt.Errorf("create case subject: %w", err)
	}

	return plain.ToResponse(), nil
}
