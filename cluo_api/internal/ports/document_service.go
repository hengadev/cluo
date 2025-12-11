package ports

import (
	"context"

	"github.com/hengadev/cluo_api/internal/domain"
)

// DocumentService handles business logic for document operations.
type DocumentService interface {
	// Generic document operations
	CreateDocument(ctx context.Context, req *domain.CreateDocumentRequest) (domain.Documentable, error)
	UpdateDocument(ctx context.Context, id string, req *domain.UpdateDocumentRequest) (domain.Documentable, error)
	GetDocument(ctx context.Context, id string, docType domain.DocumentType) (domain.Documentable, error)
	DeleteDocument(ctx context.Context, id string, docType domain.DocumentType) error
	ListDocuments(ctx context.Context, filter domain.DocumentFilter, pagination domain.Pagination) ([]domain.DocumentSummary, int, error)
	SendDocument(ctx context.Context, id string, docType domain.DocumentType, req *domain.SendDocumentRequest) error
	SignDocument(ctx context.Context, id string, docType domain.DocumentType, req *domain.SignDocumentRequest) error
	ArchiveDocument(ctx context.Context, id string, docType domain.DocumentType) error
	GetDocumentHistory(ctx context.Context, id string, docType domain.DocumentType, pagination domain.Pagination) ([]*domain.DocumentVersion, int, error)

	// Estimate operations
	CreateEstimate(ctx context.Context, estimate *domain.Estimate) (*domain.Estimate, error)
	AcceptEstimate(ctx context.Context, estimateID string, acceptedBy string) (*domain.Mandate, error)
	UpdateEstimate(ctx context.Context, estimateID string, lineItems []domain.EstimateItem) (*domain.Estimate, error)

	// Mandate operations
	CreateMandate(ctx context.Context, mandate *domain.Mandate) (*domain.Mandate, error)
	SignMandate(ctx context.Context, mandateID string, req *domain.SignDocumentRequest) (*domain.Mandate, error)
	ActivateMandate(ctx context.Context, mandateID string) (*domain.Mandate, error)
	CreateContractFromMandate(ctx context.Context, mandateID string, contract *domain.Contract) (*domain.Contract, error)

	// Contract operations
	CreateContract(ctx context.Context, contract *domain.Contract) (*domain.Contract, error)
	SignContract(ctx context.Context, contractID string, req *domain.SignDocumentRequest) (*domain.Contract, error)
	ActivateContract(ctx context.Context, contractID string) (*domain.Contract, error)
	CreateInvoiceFromContract(ctx context.Context, contractID string, invoice *domain.Invoice) (*domain.Invoice, error)

	// Invoice operations
	CreateInvoice(ctx context.Context, invoice *domain.Invoice) (*domain.Invoice, error)
	ProcessPayment(ctx context.Context, invoiceID string, req *domain.PaymentRequest) (*domain.Invoice, error)
	VoidInvoice(ctx context.Context, invoiceID string) (*domain.Invoice, error)
	GetOverdueInvoices(ctx context.Context, pagination domain.Pagination) ([]*domain.Invoice, int, error)

	// Document linking and workflow
	GetDocumentWorkflow(ctx context.Context, caseID string) ([]domain.DocumentSummary, error)
	ValidateDocumentTransitions(ctx context.Context, id string, docType domain.DocumentType, newStatus domain.DocumentStatus) error
}