package ports

import (
	"context"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// DocumentService handles business logic for document operations.
type DocumentService interface {
	// Generic document operations
	CreateDocument(ctx context.Context, req *document.CreateDocumentRequest) (document.Documentable, error)
	UpdateDocument(ctx context.Context, id string, req *document.UpdateDocumentRequest) (document.Documentable, error)
	GetDocument(ctx context.Context, id string, docType document.DocumentType) (document.Documentable, error)
	DeleteDocument(ctx context.Context, id string, docType document.DocumentType) error
	ListDocuments(ctx context.Context, filter document.DocumentFilter, pagination document.Pagination) ([]document.DocumentSummary, int, error)
	SendDocument(ctx context.Context, id string, docType document.DocumentType, req *document.SendDocumentRequest) error
	SignDocument(ctx context.Context, id string, docType document.DocumentType, req *document.SignDocumentRequest) error
	ArchiveDocument(ctx context.Context, id string, docType document.DocumentType) error
	GetDocumentHistory(ctx context.Context, id string, docType document.DocumentType, pagination document.Pagination) ([]*document.DocumentVersion, int, error)
	RenderDocumentPDF(ctx context.Context, id string, docType document.DocumentType) ([]byte, error)

	// Estimate operations
	CreateEstimate(ctx context.Context, estimate *document.Estimate) (*document.Estimate, error)
	AcceptEstimate(ctx context.Context, estimateID string, acceptedBy string) (*document.Mandate, error)
	UpdateEstimate(ctx context.Context, estimateID string, lineItems []document.EstimateItem) (*document.Estimate, error)

	// Mandate operations
	CreateMandate(ctx context.Context, mandate *document.Mandate) (*document.Mandate, error)
	SignMandate(ctx context.Context, mandateID string, req *document.SignDocumentRequest) (*document.Mandate, error)
	ActivateMandate(ctx context.Context, mandateID string) (*document.Mandate, error)
	CreateContractFromMandate(ctx context.Context, mandateID string, contract *document.Contract) (*document.Contract, error)

	// Contract operations
	CreateContract(ctx context.Context, contract *document.Contract) (*document.Contract, error)
	SignContract(ctx context.Context, contractID string, req *document.SignDocumentRequest) (*document.Contract, error)
	ActivateContract(ctx context.Context, contractID string) (*document.Contract, error)
	CreateInvoiceFromContract(ctx context.Context, contractID string) (*document.Invoice, error)

	// Invoice operations
	CreateInvoice(ctx context.Context, invoice *document.Invoice) (*document.Invoice, error)
	ProcessPayment(ctx context.Context, invoiceID string, req *document.PaymentRequest) (*document.Invoice, error)
	VoidInvoice(ctx context.Context, invoiceID string) (*document.Invoice, error)
	GetOverdueInvoices(ctx context.Context, pagination document.Pagination) ([]*document.Invoice, int, error)

	// Document linking and workflow
	GetDocumentWorkflow(ctx context.Context, caseID string) (*document.DocumentWorkflowResponse, error)
	ValidateDocumentTransitions(ctx context.Context, id string, docType document.DocumentType, newStatus document.DocumentStatus) error
}
