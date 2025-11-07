package ports

// import (
// 	"context"
//
// 	"github.com/hengadev/cluo_api/internal/domain"
// )
//
// // DocumentRepository handles persistent storage for all document types.
// type DocumentRepository interface {
// 	// Generic document operations
// 	Create(ctx context.Context, doc domain.Documentable) error
// 	GetByID(ctx context.Context, id string, docType domain.DocumentType) (domain.Documentable, error)
// 	Update(ctx context.Context, doc domain.Documentable) error
// 	Delete(ctx context.Context, id string, docType domain.DocumentType) error
// 	List(ctx context.Context, filter domain.DocumentFilter, pagination domain.Pagination) ([]domain.DocumentSummary, int, error)
//
// 	// Estimate operations
// 	CreateEstimate(ctx context.Context, estimate *domain.Estimate) error
// 	GetEstimateByID(ctx context.Context, id string) (*domain.Estimate, error)
// 	UpdateEstimate(ctx context.Context, estimate *domain.Estimate) error
// 	DeleteEstimate(ctx context.Context, id string) error
// 	ListEstimatesByCase(ctx context.Context, caseID string, pagination domain.Pagination) ([]*domain.Estimate, int, error)
//
// 	// Mandate operations
// 	CreateMandate(ctx context.Context, mandate *domain.Mandate) error
// 	GetMandateByID(ctx context.Context, id string) (*domain.Mandate, error)
// 	UpdateMandate(ctx context.Context, mandate *domain.Mandate) error
// 	DeleteMandate(ctx context.Context, id string) error
// 	ListMandatesByCase(ctx context.Context, caseID string, pagination domain.Pagination) ([]*domain.Mandate, int, error)
//
// 	// Contract operations
// 	CreateContract(ctx context.Context, contract *domain.Contract) error
// 	GetContractByID(ctx context.Context, id string) (*domain.Contract, error)
// 	UpdateContract(ctx context.Context, contract *domain.Contract) error
// 	DeleteContract(ctx context.Context, id string) error
// 	ListContractsByCase(ctx context.Context, caseID string, pagination domain.Pagination) ([]*domain.Contract, int, error)
//
// 	// Invoice operations
// 	CreateInvoice(ctx context.Context, invoice *domain.Invoice) error
// 	GetInvoiceByID(ctx context.Context, id string) (*domain.Invoice, error)
// 	UpdateInvoice(ctx context.Context, invoice *domain.Invoice) error
// 	DeleteInvoice(ctx context.Context, id string) error
// 	ListInvoicesByCase(ctx context.Context, caseID string, pagination domain.Pagination) ([]*domain.Invoice, int, error)
//
// 	// Document linking operations
// 	GetLinkedDocuments(ctx context.Context, documentID string, docType domain.DocumentType) ([]domain.Documentable, error)
// }
//
// // DocumentVersionRepository handles versioning and audit trail for documents.
// type DocumentVersionRepository interface {
// 	CreateVersion(ctx context.Context, version *domain.DocumentVersion) error
// 	GetDocumentHistory(ctx context.Context, documentID string, docType domain.DocumentType, pagination domain.Pagination) ([]*domain.DocumentVersion, int, error)
// 	GetVersion(ctx context.Context, documentID string, docType domain.DocumentType, version int) (*domain.DocumentVersion, error)
// 	DeleteVersions(ctx context.Context, documentID string, docType domain.DocumentType) error
// }
//
// // DocumentService handles business logic for document operations.
// type DocumentService interface {
// 	// Generic document operations
// 	CreateDocument(ctx context.Context, req *domain.CreateDocumentRequest) (domain.Documentable, error)
// 	UpdateDocument(ctx context.Context, id string, req *domain.UpdateDocumentRequest) (domain.Documentable, error)
// 	GetDocument(ctx context.Context, id string, docType domain.DocumentType) (domain.Documentable, error)
// 	DeleteDocument(ctx context.Context, id string, docType domain.DocumentType) error
// 	ListDocuments(ctx context.Context, filter domain.DocumentFilter, pagination domain.Pagination) ([]domain.DocumentSummary, int, error)
// 	SendDocument(ctx context.Context, id string, docType domain.DocumentType, req *domain.SendDocumentRequest) error
// 	SignDocument(ctx context.Context, id string, docType domain.DocumentType, req *domain.SignDocumentRequest) error
// 	ArchiveDocument(ctx context.Context, id string, docType domain.DocumentType) error
// 	GetDocumentHistory(ctx context.Context, id string, docType domain.DocumentType, pagination domain.Pagination) ([]*domain.DocumentVersion, int, error)
//
// 	// Estimate operations
// 	CreateEstimate(ctx context.Context, estimate *domain.Estimate) (*domain.Estimate, error)
// 	AcceptEstimate(ctx context.Context, estimateID string, acceptedBy string) (*domain.Mandate, error)
// 	UpdateEstimate(ctx context.Context, estimateID string, lineItems []domain.EstimateItem) (*domain.Estimate, error)
//
// 	// Mandate operations
// 	CreateMandate(ctx context.Context, mandate *domain.Mandate) (*domain.Mandate, error)
// 	SignMandate(ctx context.Context, mandateID string, req *domain.SignDocumentRequest) (*domain.Mandate, error)
// 	ActivateMandate(ctx context.Context, mandateID string) (*domain.Mandate, error)
// 	CreateContractFromMandate(ctx context.Context, mandateID string, contract *domain.Contract) (*domain.Contract, error)
//
// 	// Contract operations
// 	CreateContract(ctx context.Context, contract *domain.Contract) (*domain.Contract, error)
// 	SignContract(ctx context.Context, contractID string, req *domain.SignDocumentRequest) (*domain.Contract, error)
// 	ActivateContract(ctx context.Context, contractID string) (*domain.Contract, error)
// 	CreateInvoiceFromContract(ctx context.Context, contractID string, invoice *domain.Invoice) (*domain.Invoice, error)
//
// 	// Invoice operations
// 	CreateInvoice(ctx context.Context, invoice *domain.Invoice) (*domain.Invoice, error)
// 	ProcessPayment(ctx context.Context, invoiceID string, req *domain.PaymentRequest) (*domain.Invoice, error)
// 	VoidInvoice(ctx context.Context, invoiceID string) (*domain.Invoice, error)
// 	GetOverdueInvoices(ctx context.Context, pagination domain.Pagination) ([]*domain.Invoice, int, error)
//
// 	// Document linking and workflow
// 	GetDocumentWorkflow(ctx context.Context, caseID string) ([]domain.DocumentSummary, error)
// 	ValidateDocumentTransitions(ctx context.Context, id string, docType domain.DocumentType, newStatus domain.DocumentStatus) error
// }
