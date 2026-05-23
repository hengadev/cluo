package document

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// SignMandate signs a mandate.
func (s *Service) SignMandate(ctx context.Context, mandateID string, req *document.SignDocumentRequest) (*document.Mandate, error) {
	// Validate request
	if err := req.Valid(ctx); err != nil {
		return nil, errs.NewInvalidValueErr(err.Error())
	}

	// Get mandate
	mandateEncx, err := s.repo.GetMandateByID(ctx, mandateID)
	if err != nil {
		return nil, errs.NewNotFoundErr(err, "mandate")
	}

	// Decrypt mandate
	mandate, err := document.DecryptMandateEncx(ctx, s.crypto, mandateEncx)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("mandate", err)
	}

	// Create signature
	signature := document.Signature{
		ID:        uuid.New(),
		Name:      req.SignerName,
		Role:      req.SignerRole,
		Method:    req.Method,
		IPAddress: req.IPAddress,
		UserAgent: req.UserAgent,
		SignedAt:  time.Now(),
	}

	// Apply signature
	if req.SignerRole == "client" {
		if mandate.ClientSignature != nil {
			return nil, errs.NewConflictErr(fmt.Errorf("client signature already exists"))
		}
		if err := mandate.AddClientSignature(signature); err != nil {
			return nil, errs.NewInvalidValueErr(err.Error())
		}
	} else if req.SignerRole == "investigator" {
		if mandate.InvestigatorSignature != nil {
			return nil, errs.NewConflictErr(fmt.Errorf("investigator signature already exists"))
		}
		if err := mandate.AddInvestigatorSignature(signature); err != nil {
			return nil, errs.NewInvalidValueErr(err.Error())
		}
	} else {
		return nil, errs.NewInvalidValueErr("invalid signer role for mandate")
	}

	// Encrypt updated mandate
	updatedMandateEncx, err := document.ProcessMandateEncx(ctx, s.crypto, mandate)
	if err != nil {
		return nil, errs.NewNotDecryptedErr("mandate", err)
	}

	// Save updates
	if err := s.repo.UpdateMandate(ctx, updatedMandateEncx); err != nil {
		return nil, errs.NewNotUpdatedErr(err, "mandate")
	}

	// Create version record
	if err := s.createDocumentVersion(ctx, mandate, nil, stringPtr("Mandate signed")); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return mandate, nil
}