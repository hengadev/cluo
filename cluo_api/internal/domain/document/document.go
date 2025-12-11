package document

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Documentable is the shared interface — all concrete document types implement this
type Documentable interface {
	GetID() uuid.UUID
	GetType() DocumentType
	GetCaseID() uuid.UUID
	Validate() error
	GetStatus() DocumentStatus
}

// Document represents any document in the system.
// This is a base interface that all document types implement.
type Document interface {
	Documentable

	// Common operations that all documents should support
	UpdateTimestamp()
	SetStatus(status DocumentStatus)
}

// DocumentVersion represents a version of a document for audit purposes.
type DocumentVersion struct {
	ID         int64        `json:"id" db:"id"`
	DocumentID uuid.UUID    `json:"document_id" db:"document_id"`
	DocType    DocumentType `json:"doc_type" db:"doc_type"`
	Version    int          `json:"version" db:"version"`
	AuthorID   *uuid.UUID   `json:"author_id,omitempty" db:"author_id"`
	Data       []byte       `json:"data" db:"data"` // Serialized document data (JSONB)
	CreatedAt  time.Time    `json:"created_at" db:"created_at"`
	Reason     *string      `json:"reason,omitempty" db:"reason"`
}

// DocumentSummary represents a summary view of a document for listings.
type DocumentSummary struct {
	ID          uuid.UUID      `json:"id"`
	CaseID      uuid.UUID      `json:"case_id"`
	ClientID    uuid.UUID      `json:"client_id"`
	Type        DocumentType   `json:"type"`
	Status      DocumentStatus `json:"status"`
	DocumentRef string         `json:"document_ref"` // EstimateNumber, MandateNumber, etc.
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	// TODO: Add summary fields specific to each document type as needed
}

// CreateDocumentRequest represents a request to create a new document.
type CreateDocumentRequest struct {
	Type     DocumentType `json:"type" validate:"required"`
	CaseID   uuid.UUID    `json:"case_id" validate:"required"`
	ClientID uuid.UUID    `json:"client_id" validate:"required"`
	Data     interface{}  `json:"data" validate:"required"`
}

// UpdateDocumentRequest represents a request to update a document.
type UpdateDocumentRequest struct {
	Data   interface{} `json:"data" validate:"required"`
	Reason *string     `json:"reason,omitempty"`
}

// SendDocumentRequest represents a request to send a document.
type SendDocumentRequest struct {
	Recipients []string `json:"recipients" validate:"required,min=1"`
	Subject    *string  `json:"subject,omitempty"`
	Message    *string  `json:"message,omitempty"`
	SendEmail  bool     `json:"send_email"`
	SendSMS    bool     `json:"send_sms"`
}

// SignDocumentRequest represents a request to sign a document.
type SignDocumentRequest struct {
	SignerName       string  `json:"signer_name" validate:"required"`
	SignerRole       string  `json:"signer_role" validate:"required"`
	SignatureFileURL string  `json:"signature_file_url"`
	Method           string  `json:"method" validate:"required,oneof=e-sign wet pdf-stamp third-party"`
	IPAddress        *string `json:"ip_address,omitempty"`
	UserAgent        *string `json:"user_agent,omitempty"`
}

// PaymentRequest represents a request to make a payment on an invoice.
type PaymentRequest struct {
	Amount        float64 `json:"amount" validate:"required,gt=0"`
	PaymentMethod string  `json:"payment_method" validate:"required"`
	TransactionID *string `json:"transaction_id,omitempty"`
	Notes         *string `json:"notes,omitempty"`
}

// DocumentFilter represents filtering options for document queries.
type DocumentFilter struct {
	Type     *DocumentType   `json:"type,omitempty"`
	Status   *DocumentStatus `json:"status,omitempty"`
	CaseID   *uuid.UUID      `json:"case_id,omitempty"`
	ClientID *uuid.UUID      `json:"client_id,omitempty"`
	DateFrom *time.Time      `json:"date_from,omitempty"`
	DateTo   *time.Time      `json:"date_to,omitempty"`
	Search   *string         `json:"search,omitempty"`
}

// Pagination represents pagination parameters.
type Pagination struct {
	Page     int `json:"page" validate:"min=1"`
	PageSize int `json:"page_size" validate:"min=1,max=100"`
}

// NewPagination creates a new pagination with default values.
func NewPagination() Pagination {
	return Pagination{
		Page:     1,
		PageSize: 20,
	}
}

// GetOffset calculates the offset for database queries.
func (p Pagination) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// Validate performs validation on the pagination parameters.
// TODO: Implement pagination validation:
// - Page: must be positive integer
// - PageSize: must be positive integer, reasonable upper limit (e.g., 100)
// - Add validation for extremely large page numbers that might cause performance issues
func (p Pagination) Validate() error {
	// TODO: Add pagination validation implementation
	if p.Page < 1 {
		return fmt.Errorf("page must be at least 1")
	}
	if p.PageSize < 1 {
		return fmt.Errorf("page size must be at least 1")
	}
	if p.PageSize > 100 {
		return fmt.Errorf("page size cannot exceed 100")
	}
	return nil
}
