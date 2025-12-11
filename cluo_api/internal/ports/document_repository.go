package ports

import (
	"context"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// DocumentRepository handles persistent storage for all document types.
type DocumentRepository interface {
	// Generic document operations
	Create(ctx context.Context, doc document.Documentable) error
	GetByID(ctx context.Context, id string, docType document.DocumentType) (document.Documentable, error)
	Update(ctx context.Context, doc document.Documentable) error
	Delete(ctx context.Context, id string, docType document.DocumentType) error
	List(ctx context.Context, filter document.DocumentFilter, pagination document.Pagination) ([]document.DocumentSummary, int, error)

	// Estimate operations
	CreateEstimate(ctx context.Context, estimate *document.Estimate) error
	GetEstimateByID(ctx context.Context, id string) (*document.Estimate, error)
	UpdateEstimate(ctx context.Context, estimate *document.Estimate) error
	DeleteEstimate(ctx context.Context, id string) error
	ListEstimatesByCase(ctx context.Context, caseID string, pagination document.Pagination) ([]*document.Estimate, int, error)

	// Mandate operations
	CreateMandate(ctx context.Context, mandate *document.Mandate) error
	GetMandateByID(ctx context.Context, id string) (*document.Mandate, error)
	UpdateMandate(ctx context.Context, mandate *document.Mandate) error
	DeleteMandate(ctx context.Context, id string) error
	ListMandatesByCase(ctx context.Context, caseID string, pagination document.Pagination) ([]*document.Mandate, int, error)

	// Contract operations
	CreateContract(ctx context.Context, contract *document.Contract) error
	GetContractByID(ctx context.Context, id string) (*document.Contract, error)
	UpdateContract(ctx context.Context, contract *document.Contract) error
	DeleteContract(ctx context.Context, id string) error
	ListContractsByCase(ctx context.Context, caseID string, pagination document.Pagination) ([]*document.Contract, int, error)

	// Invoice operations
	CreateInvoice(ctx context.Context, invoice *document.Invoice) error
	GetInvoiceByID(ctx context.Context, id string) (*document.Invoice, error)
	UpdateInvoice(ctx context.Context, invoice *document.Invoice) error
	DeleteInvoice(ctx context.Context, id string) error
	ListInvoicesByCase(ctx context.Context, caseID string, pagination document.Pagination) ([]*document.Invoice, int, error)

	// Document linking operations
	GetLinkedDocuments(ctx context.Context, documentID string, docType document.DocumentType) ([]document.Documentable, error)
}

// DocumentVersionRepository handles versioning and audit trail for documents.
type DocumentVersionRepository interface {
	CreateVersion(ctx context.Context, version *document.DocumentVersion) error
	GetDocumentHistory(ctx context.Context, documentID string, docType document.DocumentType, pagination document.Pagination) ([]*document.DocumentVersion, int, error)
	GetVersion(ctx context.Context, documentID string, docType document.DocumentType, version int) (*document.DocumentVersion, error)
	DeleteVersions(ctx context.Context, documentID string, docType document.DocumentType) error
}

