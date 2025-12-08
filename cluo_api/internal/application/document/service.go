package document

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/hengadev/cluo_api/internal/domain"
	"github.com/hengadev/cluo_api/internal/ports"
)

type Service struct {
	repo          ports.DocumentRepository
	versionRepo   ports.DocumentVersionRepository
	// TODO: Add additional dependencies like:
	// - email service for sending documents
	// - pdf generation service
	// - notification service
	// - user service for authorization
}

func New(repo ports.DocumentRepository, versionRepo ports.DocumentVersionRepository) ports.DocumentService {
	return &Service{
		repo:        repo,
		versionRepo: versionRepo,
	}
}

// Helper functions for document management

// createDocumentVersion creates a version record for a document
func (s *Service) createDocumentVersion(ctx context.Context, doc domain.Documentable, authorID *uuid.UUID, reason *string) error {
	// Serialize document data
	data, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to serialize document: %w", err)
	}

	// Get current version count
	pagination := domain.NewPagination()
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
	version := &domain.DocumentVersion{
		ID:        uuid.New(),
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

// validateDocumentTransition checks if a status transition is allowed
func (s *Service) validateDocumentTransition(doc domain.Documentable, newStatus domain.DocumentStatus) error {
	currentStatus := doc.GetStatus()

	// TODO: Implement comprehensive status transition validation:
	// - Check business rules for each document type
	// - Verify user permissions for the transition
	// - Ensure all required fields are set for the new status
	// - Check if document can be modified in current state
	// - Validate linking constraints (e.g., can't delete if linked documents exist)

	if !currentStatus.CanTransitionTo(newStatus) {
		return fmt.Errorf("cannot transition from %s to %s", currentStatus, newStatus)
	}

	// Document-specific validation
	switch d := doc.(type) {
	case *domain.Estimate:
		if newStatus == domain.DocumentStatusActive && !d.Accepted {
			return fmt.Errorf("estimate must be accepted before becoming active")
		}
		if newStatus == domain.DocumentStatusSigned && !d.Accepted {
			return fmt.Errorf("estimate must be accepted before being signed")
		}

	case *domain.Mandate:
		if newStatus == domain.DocumentStatusActive && d.ClientSignature == nil {
			return fmt.Errorf("mandate must have client signature to be active")
		}
		if d.IsExpired() && newStatus == domain.DocumentStatusActive {
			return fmt.Errorf("cannot activate expired mandate")
		}

	case *domain.Contract:
		if newStatus == domain.DocumentStatusActive && len(d.Signatures) == 0 {
			return fmt.Errorf("contract must have at least one signature to be active")
		}
		if d.IsExpired() && newStatus == domain.DocumentStatusActive {
			return fmt.Errorf("cannot activate expired contract")
		}

	case *domain.Invoice:
		if newStatus == domain.DocumentStatusActive && d.PaymentStatus != domain.PaymentStatusPaid {
			return fmt.Errorf("invoice must be paid to be active")
		}
	}

	return nil
}

// Generic document operations

func (s *Service) CreateDocument(ctx context.Context, req *domain.CreateDocumentRequest) (domain.Documentable, error) {
	// TODO: Validate request parameters
	// TODO: Verify user permissions to create documents
	// TODO: Check that case and client exist and are accessible

	var doc domain.Documentable
	var err error

	switch req.Type {
	case domain.DocumentTypeEstimate:
		// TODO: Type assert req.Data to *domain.Estimate
		// estimate := req.Data.(*domain.Estimate)
		// doc = estimate
		return nil, fmt.Errorf("estimate creation not implemented - use CreateEstimate instead")

	case domain.DocumentTypeMandate:
		// TODO: Type assert req.Data to *domain.Mandate
		return nil, fmt.Errorf("mandate creation not implemented - use CreateMandate instead")

	case domain.DocumentTypeContract:
		// TODO: Type assert req.Data to *domain.Contract
		return nil, fmt.Errorf("contract creation not implemented - use CreateContract instead")

	case domain.DocumentTypeInvoice:
		// TODO: Type assert req.Data to *domain.Invoice
		return nil, fmt.Errorf("invoice creation not implemented - use CreateInvoice instead")

	default:
		return nil, fmt.Errorf("unsupported document type: %s", req.Type)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create document: %w", err)
	}

	// Validate document
	if err := doc.Validate(); err != nil {
		return nil, fmt.Errorf("document validation failed: %w", err)
	}

	// Save to repository
	if err := s.repo.Create(ctx, doc); err != nil {
		return nil, fmt.Errorf("failed to save document: %w", err)
	}

	// Create initial version
	if err := s.createDocumentVersion(ctx, doc, nil, stringPtr("Initial creation")); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return doc, nil
}

func (s *Service) UpdateDocument(ctx context.Context, id string, req *domain.UpdateDocumentRequest) (domain.Documentable, error) {
	// TODO: Validate request parameters
	// TODO: Verify user permissions to update documents

	// Get existing document
	doc, err := s.repo.GetByID(ctx, id, domain.DocumentType(req.Data.(interface{}).(domain.Documentable).GetType()))
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	// TODO: Apply updates from req.Data to the document
	// TODO: Validate updated document

	// Save to repository
	if err := s.repo.Update(ctx, doc); err != nil {
		return nil, fmt.Errorf("failed to update document: %w", err)
	}

	// Create version record
	// TODO: Get author ID from context
	if err := s.createDocumentVersion(ctx, doc, nil, req.Reason); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return doc, nil
}

func (s *Service) GetDocument(ctx context.Context, id string, docType domain.DocumentType) (domain.Documentable, error) {
	// TODO: Verify user permissions to access documents

	doc, err := s.repo.GetByID(ctx, id, docType)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	return doc, nil
}

func (s *Service) DeleteDocument(ctx context.Context, id string, docType domain.DocumentType) error {
	// TODO: Verify user permissions to delete documents
	// TODO: Check if document can be deleted (no linked documents, not final status, etc.)

	// Get document to check status
	doc, err := s.repo.GetByID(ctx, id, docType)
	if err != nil {
		return fmt.Errorf("failed to get document: %w", err)
	}

	if doc.GetStatus().IsFinal() {
		return fmt.Errorf("cannot delete document in final status: %s", doc.GetStatus())
	}

	// Check for linked documents
	linkedDocs, err := s.repo.GetLinkedDocuments(ctx, id, docType)
	if err != nil {
		return fmt.Errorf("failed to check linked documents: %w", err)
	}

	if len(linkedDocs) > 0 {
		return fmt.Errorf("cannot delete document with linked documents")
	}

	// Delete document
	if err := s.repo.Delete(ctx, id, docType); err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}

	// Delete version history
	if err := s.versionRepo.DeleteVersions(ctx, id, docType); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return nil
}

func (s *Service) ListDocuments(ctx context.Context, filter domain.DocumentFilter, pagination domain.Pagination) ([]domain.DocumentSummary, int, error) {
	// TODO: Verify user permissions to list documents
	// TODO: Apply user-specific filtering (case access, client access, etc.)

	if err := pagination.Validate(); err != nil {
		return nil, 0, fmt.Errorf("invalid pagination: %w", err)
	}

	documents, total, err := s.repo.List(ctx, filter, pagination)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list documents: %w", err)
	}

	return documents, total, nil
}

func (s *Service) SendDocument(ctx context.Context, id string, docType domain.DocumentType, req *domain.SendDocumentRequest) error {
	// TODO: Validate request parameters
	// TODO: Verify user permissions to send documents

	// Get document
	doc, err := s.repo.GetByID(ctx, id, docType)
	if err != nil {
		return fmt.Errorf("failed to get document: %w", err)
	}

	// Check if document can be sent
	if doc.GetStatus() != domain.DocumentStatusDraft {
		return fmt.Errorf("only draft documents can be sent")
	}

	// TODO: Implement actual sending logic:
	// - Generate PDF
	// - Send emails
	// - Send SMS notifications
	// - Create notification events
	// - Update document status

	// Update document status
	// TODO: This would need to be done through a type-specific update method
	// For now, we'll return an error indicating the operation isn't fully implemented

	return fmt.Errorf("document sending not fully implemented")
}

func (s *Service) SignDocument(ctx context.Context, id string, docType domain.DocumentType, req *domain.SignDocumentRequest) error {
	// TODO: Validate request parameters
	// TODO: Verify user permissions to sign documents
	// TODO: Validate signature data

	// Get document
	doc, err := s.repo.GetByID(ctx, id, docType)
	if err != nil {
		return fmt.Errorf("failed to get document: %w", err)
	}

	// TODO: Implement document-specific signing logic
	// This would involve calling the appropriate method on the document
	// and then updating it in the repository

	return fmt.Errorf("document signing not fully implemented - use document-specific methods")
}

func (s *Service) ArchiveDocument(ctx context.Context, id string, docType domain.DocumentType) error {
	// TODO: Verify user permissions to archive documents

	// Get document
	doc, err := s.repo.GetByID(ctx, id, docType)
	if err != nil {
		return fmt.Errorf("failed to get document: %w", err)
	}

	// Check if document can be archived
	if doc.GetStatus() == domain.DocumentStatusArchived {
		return fmt.Errorf("document is already archived")
	}

	// TODO: Implement archiving logic
	// This would involve setting the status to archived and updating the document

	return fmt.Errorf("document archiving not fully implemented")
}

func (s *Service) GetDocumentHistory(ctx context.Context, id string, docType domain.DocumentType, pagination domain.Pagination) ([]*domain.DocumentVersion, int, error) {
	// TODO: Verify user permissions to access document history

	if err := pagination.Validate(); err != nil {
		return nil, 0, fmt.Errorf("invalid pagination: %w", err)
	}

	versions, total, err := s.versionRepo.GetDocumentHistory(ctx, id, docType, pagination)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get document history: %w", err)
	}

	return versions, total, nil
}

// Helper function for string pointers
func stringPtr(s string) *string {
	return &s
}
