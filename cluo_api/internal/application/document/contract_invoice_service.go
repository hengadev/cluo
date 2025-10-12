package document

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/hengadev/cluo_api/internal/domain"
	"github.com/hengadev/cluo_api/internal/ports"
)

// Contract operations

func (s *Service) CreateContract(ctx context.Context, contract *domain.Contract) (*domain.Contract, error) {
	// TODO: Validate user permissions to create contracts
	// TODO: Verify that case and client exist and are accessible

	// Validate contract
	if err := contract.Validate(); err != nil {
		return nil, fmt.Errorf("contract validation failed: %w", err)
	}

	// Generate contract number if not provided
	if contract.ContractNumber == "" {
		// TODO: Implement contract number generation
		contract.ContractNumber = fmt.Sprintf("CNT-%d", uuid.New().ID())
	}

	// Save to repository
	if err := s.repo.CreateContract(ctx, contract); err != nil {
		return nil, fmt.Errorf("failed to create contract: %w", err)
	}

	// Create initial version
	if err := s.createDocumentVersion(ctx, contract, nil, stringPtr("Initial creation")); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return contract, nil
}

func (s *Service) SignContract(ctx context.Context, contractID string, req *domain.SignDocumentRequest) (*domain.Contract, error) {
	// TODO: Validate user permissions to sign contracts
	// TODO: Verify that signer has authority to sign for the specified role

	// Get contract
	contract, err := s.repo.GetContractByID(ctx, contractID)
	if err != nil {
		return nil, fmt.Errorf("failed to get contract: %w", err)
	}

	// Check if contract can be signed by this person
	if !contract.CanBeSigned(req.SignerName, req.SignerRole) {
		return nil, fmt.Errorf("contract cannot be signed by %s as %s", req.SignerName, req.SignerRole)
	}

	// Create signature
	signature := domain.NewSignature(
		req.SignerName,
		req.SignerRole,
		req.Method,
		req.SignatureFileURL,
		nil, // TODO: Get signer ID from context
	)

	// Add additional fields if provided
	if req.IPAddress != nil {
		signature.IPAddress = req.IPAddress
	}
	if req.UserAgent != nil {
		signature.UserAgent = req.UserAgent
	}

	// Add signature to contract
	if err := contract.AddSignature(signature); err != nil {
		return nil, fmt.Errorf("failed to add signature: %w", err)
	}

	// Validate updated contract
	if err := contract.Validate(); err != nil {
		return nil, fmt.Errorf("updated contract validation failed: %w", err)
	}

	// Update contract
	if err := s.repo.UpdateContract(ctx, contract); err != nil {
		return nil, fmt.Errorf("failed to update contract: %w", err)
	}

	// Create version record
	// TODO: Get author ID from context
	signerUUID := uuid.New() // TODO: Get actual signer ID from context
	reason := fmt.Sprintf("Signed by %s as %s", req.SignerName, req.SignerRole)
	if err := s.createDocumentVersion(ctx, contract, &signerUUID, stringPtr(reason)); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return contract, nil
}

func (s *Service) ActivateContract(ctx context.Context, contractID string) (*domain.Contract, error) {
	// TODO: Validate user permissions to activate contracts

	// Get contract
	contract, err := s.repo.GetContractByID(ctx, contractID)
	if err != nil {
		return nil, fmt.Errorf("failed to get contract: %w", err)
	}

	// Activate contract
	if err := contract.Activate(); err != nil {
		return nil, fmt.Errorf("failed to activate contract: %w", err)
	}

	// Update contract
	if err := s.repo.UpdateContract(ctx, contract); err != nil {
		return nil, fmt.Errorf("failed to update contract: %w", err)
	}

	// Create version record
	// TODO: Get author ID from context
	if err := s.createDocumentVersion(ctx, contract, nil, stringPtr("Contract activated")); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return contract, nil
}

func (s *Service) CreateInvoiceFromContract(ctx context.Context, contractID string, invoice *domain.Invoice) (*domain.Invoice, error) {
	// TODO: Validate user permissions to create invoices
	// TODO: Verify that case and client exist and are accessible

	// Get contract
	contract, err := s.repo.GetContractByID(ctx, contractID)
	if err != nil {
		return nil, fmt.Errorf("failed to get contract: %w", err)
	}

	// Verify contract is in appropriate state
	if contract.Status != domain.DocumentStatusSigned && contract.Status != domain.DocumentStatusActive {
		return nil, fmt.Errorf("contract must be signed or active to create invoice")
	}

	// Link invoice to contract
	if err := invoice.LinkToContract(contract.ID); err != nil {
		return nil, fmt.Errorf("failed to link invoice to contract: %w", err)
	}

	// Ensure invoice belongs to same case and client
	if invoice.CaseID != contract.CaseID {
		return nil, fmt.Errorf("invoice must belong to same case as contract")
	}
	if invoice.ClientID != contract.ClientID {
		return nil, fmt.Errorf("invoice must belong to same client as contract")
	}

	// Generate invoice number if not provided
	if invoice.InvoiceNumber == "" {
		// TODO: Implement invoice number generation
		invoice.InvoiceNumber = fmt.Sprintf("INV-%d", uuid.New().ID())
	}

	// Set due date if not provided (default 30 days from issue date)
	if invoice.DueDate.IsZero() {
		invoice.DueDate = invoice.IssueDate.AddDate(0, 0, 30) // 30 days from issue date
	}

	// Validate invoice
	if err := invoice.Validate(); err != nil {
		return nil, fmt.Errorf("invoice validation failed: %w", err)
	}

	// Save invoice
	if err := s.repo.CreateInvoice(ctx, invoice); err != nil {
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	// Create version record
	// TODO: Get author ID from context
	reason := fmt.Sprintf("Created from contract %s", contract.ContractNumber)
	if err := s.createDocumentVersion(ctx, invoice, nil, stringPtr(reason)); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return invoice, nil
}

// Invoice operations

func (s *Service) CreateInvoice(ctx context.Context, invoice *domain.Invoice) (*domain.Invoice, error) {
	// TODO: Validate user permissions to create invoices
	// TODO: Verify that case and client exist and are accessible

	// Validate invoice
	if err := invoice.Validate(); err != nil {
		return nil, fmt.Errorf("invoice validation failed: %w", err)
	}

	// Generate invoice number if not provided
	if invoice.InvoiceNumber == "" {
		// TODO: Implement invoice number generation
		invoice.InvoiceNumber = fmt.Sprintf("INV-%d", uuid.New().ID())
	}

	// Set due date if not provided (default 30 days from issue date)
	if invoice.DueDate.IsZero() {
		invoice.DueDate = invoice.IssueDate.AddDate(0, 0, 30) // 30 days from issue date
	}

	// Save to repository
	if err := s.repo.CreateInvoice(ctx, invoice); err != nil {
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	// Create initial version
	if err := s.createDocumentVersion(ctx, invoice, nil, stringPtr("Initial creation")); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return invoice, nil
}

func (s *Service) ProcessPayment(ctx context.Context, invoiceID string, req *domain.PaymentRequest) (*domain.Invoice, error) {
	// TODO: Validate user permissions to process payments
	// TODO: Validate payment request data

	// Get invoice
	invoice, err := s.repo.GetInvoiceByID(ctx, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice: %w", err)
	}

	// Process payment
	if err := invoice.AddPayment(req.Amount, req.PaymentMethod); err != nil {
		return nil, fmt.Errorf("failed to process payment: %w", err)
	}

	// Update invoice
	if err := s.repo.UpdateInvoice(ctx, invoice); err != nil {
		return nil, fmt.Errorf("failed to update invoice: %w", err)
	}

	// Create version record
	// TODO: Get author ID from context
	payerUUID := uuid.New() // TODO: Get actual payer ID from context
	reason := fmt.Sprintf("Payment processed: %.2f via %s", req.Amount, req.PaymentMethod)
	if req.Notes != nil {
		reason += fmt.Sprintf(" (%s)", *req.Notes)
	}
	if err := s.createDocumentVersion(ctx, invoice, &payerUUID, stringPtr(reason)); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return invoice, nil
}

func (s *Service) VoidInvoice(ctx context.Context, invoiceID string) (*domain.Invoice, error) {
	// TODO: Validate user permissions to void invoices

	// Get invoice
	invoice, err := s.repo.GetInvoiceByID(ctx, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get invoice: %w", err)
	}

	// Void invoice
	if err := invoice.Void(); err != nil {
		return nil, fmt.Errorf("failed to void invoice: %w", err)
	}

	// Update invoice
	if err := s.repo.UpdateInvoice(ctx, invoice); err != nil {
		return nil, fmt.Errorf("failed to update invoice: %w", err)
	}

	// Create version record
	// TODO: Get author ID from context
	if err := s.createDocumentVersion(ctx, invoice, nil, stringPtr("Invoice voided")); err != nil {
		// Log error but don't fail the operation
		// TODO: Add proper logging
	}

	return invoice, nil
}

func (s *Service) GetOverdueInvoices(ctx context.Context, pagination domain.Pagination) ([]*domain.Invoice, int, error) {
	// TODO: Validate user permissions to access invoices

	if err := pagination.Validate(); err != nil {
		return nil, 0, fmt.Errorf("invalid pagination: %w", err)
	}

	// TODO: This would require a new repository method or query optimization
	// For now, return empty results
	// invoices, total, err := s.repo.GetOverdueInvoices(ctx, pagination)
	// if err != nil {
	//     return nil, 0, fmt.Errorf("failed to get overdue invoices: %w", err)
	// }

	return []*domain.Invoice{}, 0, fmt.Errorf("get overdue invoices not implemented")
}

// Document workflow operations

func (s *Service) GetDocumentWorkflow(ctx context.Context, caseID string) ([]domain.DocumentSummary, error) {
	// TODO: Validate user permissions to access case documents

	// Get all documents for the case
	filter := domain.DocumentFilter{
		CaseID: stringPtr(caseID),
	}
	pagination := domain.NewPagination()
	pagination.PageSize = 100 // Get all documents

	documents, _, err := s.repo.List(ctx, filter, pagination)
	if err != nil {
		return nil, fmt.Errorf("failed to get case documents: %w", err)
	}

	return documents, nil
}

func (s *Service) ValidateDocumentTransitions(ctx context.Context, id string, docType domain.DocumentType, newStatus domain.DocumentStatus) error {
	// TODO: Validate user permissions to validate transitions

	// Get document
	doc, err := s.repo.GetByID(ctx, id, docType)
	if err != nil {
		return fmt.Errorf("failed to get document: %w", err)
	}

	// Validate transition
	return s.validateDocumentTransition(doc, newStatus)
}

// Helper function for string pointers
func stringPtr(s string) *string {
	return &s
}