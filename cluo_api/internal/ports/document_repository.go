package ports

import (
	"context"

	"github.com/hengadev/cluo_api/internal/domain"
)

// DocumentRepository handles persistent storage for all document types.
type DocumentRepository interface {
	// Generic document operations
	Create(ctx context.Context, doc domain.Documentable) error
	GetByID(ctx context.Context, id string, docType domain.DocumentType) (domain.Documentable, error)
	Update(ctx context.Context, doc domain.Documentable) error
	Delete(ctx context.Context, id string, docType domain.DocumentType) error
	List(ctx context.Context, filter domain.DocumentFilter, pagination domain.Pagination) ([]domain.DocumentSummary, int, error)

	// Estimate operations
	CreateEstimate(ctx context.Context, estimate *domain.Estimate) error
	GetEstimateByID(ctx context.Context, id string) (*domain.Estimate, error)
	UpdateEstimate(ctx context.Context, estimate *domain.Estimate) error
	DeleteEstimate(ctx context.Context, id string) error
	ListEstimatesByCase(ctx context.Context, caseID string, pagination domain.Pagination) ([]*domain.Estimate, int, error)

	// Mandate operations
	CreateMandate(ctx context.Context, mandate *domain.Mandate) error
	GetMandateByID(ctx context.Context, id string) (*domain.Mandate, error)
	UpdateMandate(ctx context.Context, mandate *domain.Mandate) error
	DeleteMandate(ctx context.Context, id string) error
	ListMandatesByCase(ctx context.Context, caseID string, pagination domain.Pagination) ([]*domain.Mandate, int, error)

	// Contract operations
	CreateContract(ctx context.Context, contract *domain.Contract) error
	GetContractByID(ctx context.Context, id string) (*domain.Contract, error)
	UpdateContract(ctx context.Context, contract *domain.Contract) error
	DeleteContract(ctx context.Context, id string) error
	ListContractsByCase(ctx context.Context, caseID string, pagination domain.Pagination) ([]*domain.Contract, int, error)

	// Invoice operations
	CreateInvoice(ctx context.Context, invoice *domain.Invoice) error
	GetInvoiceByID(ctx context.Context, id string) (*domain.Invoice, error)
	UpdateInvoice(ctx context.Context, invoice *domain.Invoice) error
	DeleteInvoice(ctx context.Context, id string) error
	ListInvoicesByCase(ctx context.Context, caseID string, pagination domain.Pagination) ([]*domain.Invoice, int, error)

	// Document linking operations
	GetLinkedDocuments(ctx context.Context, documentID string, docType domain.DocumentType) ([]domain.Documentable, error)
}

// DocumentVersionRepository handles versioning and audit trail for documents.
type DocumentVersionRepository interface {
	CreateVersion(ctx context.Context, version *domain.DocumentVersion) error
	GetDocumentHistory(ctx context.Context, documentID string, docType domain.DocumentType, pagination domain.Pagination) ([]*domain.DocumentVersion, int, error)
	GetVersion(ctx context.Context, documentID string, docType domain.DocumentType, version int) (*domain.DocumentVersion, error)
	DeleteVersions(ctx context.Context, documentID string, docType domain.DocumentType) error
}