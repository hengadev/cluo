package document

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/hengadev/cluo_api/internal/common/auth/session"
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// createDocumentVersion creates a version record for a document.
func (s *Service) createDocumentVersion(ctx context.Context, doc document.Documentable, authorID *uuid.UUID, reason *string) error {
	// Serialize document data
	data, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to serialize document: %w", err)
	}

	// Get current version count
	pagination := document.NewPagination()
	pagination.PageSize = 1 // We only need the latest version
	history, _, err := s.versionRepo.GetDocumentHistory(ctx, doc.GetID().String(), doc.GetType(), pagination)
	if err != nil {
		return fmt.Errorf("failed to get document history: %w", err)
	}

	// Determine next version number
	nextVersion := 1
	if len(history) > 0 {
		nextVersion = history[0].Version + 1
	}

	// Create version record
	version := &document.DocumentVersion{
		DocumentID: doc.GetID(),
		DocType:   doc.GetType(),
		Version:   nextVersion,
		AuthorID:  authorID,
		Data:      data,
		Reason:    reason,
		CreatedAt: time.Now(),
	}

	return s.versionRepo.CreateVersion(ctx, version)
}

// validateDocumentTransition checks if a status transition is allowed.
func (s *Service) validateDocumentTransition(doc document.Documentable, newStatus document.DocumentStatus) error {
	currentStatus := doc.GetStatus()

	if !currentStatus.CanTransitionTo(newStatus) {
		return errs.NewConflictErr(fmt.Errorf("cannot transition document from %s to %s", currentStatus, newStatus))
	}

	// Document-specific validation
	switch d := doc.(type) {
	case *document.Estimate:
		if newStatus == document.DocumentStatusActive && !d.Accepted {
			return fmt.Errorf("estimate must be accepted before becoming active")
		}

	case *document.Mandate:
		if newStatus == document.DocumentStatusActive && d.ClientSignature == nil {
			return fmt.Errorf("mandate must have client signature to be active")
		}
		if d.IsExpired() && newStatus == document.DocumentStatusActive {
			return fmt.Errorf("cannot activate expired mandate")
		}

	case *document.Contract:
		if newStatus == document.DocumentStatusActive && len(d.Signatures) == 0 {
			return fmt.Errorf("contract must have at least one signature to be active")
		}
		if d.IsExpired() && newStatus == document.DocumentStatusActive {
			return fmt.Errorf("cannot activate expired contract")
		}

	case *document.Invoice:
		if newStatus == document.DocumentStatusActive && d.PaymentStatus != document.PaymentStatusPaid {
			return fmt.Errorf("invoice must be paid to be active")
		}
	}

	return nil
}

// stringPtr returns a pointer to a string.
func stringPtr(s string) *string {
	return &s
}

// getUserIDFromContext extracts the authenticated user's UUID from the session stored in context.
func (s *Service) getUserIDFromContext(ctx context.Context) uuid.UUID {
	info, ok := session.SessionInfoFromContext(ctx)
	if !ok || info == nil {
		return uuid.Nil
	}
	return info.UserID
}