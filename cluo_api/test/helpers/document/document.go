package documentHelpers

import (
	"fmt"
	"testing"
	"time"

	"github.com/hengadev/cluo_api/internal/domain/document"

	"github.com/google/uuid"
)

// NewTestEstimate creates a basic Estimate for testing
func NewTestEstimate(t *testing.T) *document.Estimate {
	t.Helper()

	now := time.Now()
	validUntil := now.AddDate(0, 3, 0) // 3 months from now
	notes := "Test estimate notes"

	return &document.Estimate{
		DocumentBase: document.DocumentBase{
			ID:        uuid.New(),
			CaseID:    uuid.New(),
			ClientID:  uuid.New(),
			Status:    document.DocumentStatusDraft,
			CreatedAt: now,
			UpdatedAt: now,
		},
		EstimateNumber: "EST-2024-001",
		IssueDate:      now,
		ValidUntil:     &validUntil,
		LineItems: []document.EstimateItem{
			{
				Description: "Test Service 1",
				Quantity:    10,
				UnitPrice:   100.00,
				Subtotal:    1000.00,
			},
			{
				Description: "Test Service 2",
				Quantity:    5,
				UnitPrice:   50.00,
				Subtotal:    250.00,
			},
		},
		EstimatedTotal: 1250.00,
		Notes:          &notes,
		Accepted:       false,
		AcceptedAt:     nil,
		AcceptedBy:     nil,
	}
}

// NewTestEstimateWithCaseID creates an Estimate with specific case and client IDs
func NewTestEstimateWithCaseID(t *testing.T, caseID, clientID uuid.UUID) *document.Estimate {
	t.Helper()

	estimate := NewTestEstimate(t)
	estimate.CaseID = caseID
	estimate.ClientID = clientID
	return estimate
}

// NewTestAcceptedEstimate creates an Estimate that has been accepted
func NewTestAcceptedEstimate(t *testing.T) *document.Estimate {
	t.Helper()

	estimate := NewTestEstimate(t)
	estimate.Status = document.DocumentStatusSigned
	estimate.Accepted = true
	estimate.AcceptedAt = &[]time.Time{time.Now()}[0]
	estimate.AcceptedBy = &[]uuid.UUID{uuid.New()}[0]
	return estimate
}

// NewTestMandate creates a basic Mandate for testing
func NewTestMandate(t *testing.T) *document.Mandate {
	t.Helper()

	now := time.Now()
	validUntil := now.AddDate(1, 0, 0) // 1 year from now
	specialInstructions := "Handle with care"
	jurisdiction := "California"

	return &document.Mandate{
		DocumentBase: document.DocumentBase{
			ID:        uuid.New(),
			CaseID:    uuid.New(),
			ClientID:  uuid.New(),
			Status:    document.DocumentStatusDraft,
			CreatedAt: now,
			UpdatedAt: now,
		},
		MandateNumber:   "MAND-2024-001",
		IssueDate:       now,
		ScopeOfWork:     "Investigation of financial documents",
		ValidFrom:       now,
		ValidUntil:      &validUntil,
		TermsConditions: "Standard terms and conditions apply",
		ClientSignature: &document.Signature{
			Name:      "John Doe",
			Role:      "Client",
			Signature: "signature-base64",
			SignedAt:  now,
		},
		InvestigatorSignature: nil,
		LinkedEstimateID:      nil,
		SpecialInstructions:   &specialInstructions,
		Jurisdiction:          &jurisdiction,
	}
}

// NewTestMandateWithCaseID creates a Mandate with specific case and client IDs
func NewTestMandateWithCaseID(t *testing.T, caseID, clientID uuid.UUID) *document.Mandate {
	t.Helper()

	mandate := NewTestMandate(t)
	mandate.CaseID = caseID
	mandate.ClientID = clientID
	return mandate
}

// NewTestMandateWithLinkedEstimate creates a Mandate linked to an estimate
func NewTestMandateWithLinkedEstimate(t *testing.T, estimateID uuid.UUID) *document.Mandate {
	t.Helper()

	mandate := NewTestMandate(t)
	mandate.LinkedEstimateID = &estimateID
	return mandate
}

// NewTestSignedMandate creates a Mandate that is fully signed
func NewTestSignedMandate(t *testing.T) *document.Mandate {
	t.Helper()

	mandate := NewTestMandate(t)
	mandate.Status = document.DocumentStatusSigned
	mandate.InvestigatorSignature = &document.Signature{
		Name:      "Jane Smith",
		Role:      "Lead Investigator",
		Signature: "investigator-signature-base64",
		SignedAt:  time.Now(),
	}
	return mandate
}

// NewTestContract creates a basic Contract for testing
func NewTestContract(t *testing.T) *document.Contract {
	t.Helper()

	now := time.Now()
	startDate := now
	endDate := now.AddDate(1, 0, 0) // 1 year from now
	paymentTerms := "Net 30 days"
	confidentiality := "Strict confidentiality required"
	terminationClause := "30 days notice required"

	return &document.Contract{
		DocumentBase: document.DocumentBase{
			ID:        uuid.New(),
			CaseID:    uuid.New(),
			ClientID:  uuid.New(),
			Status:    document.DocumentStatusDraft,
			CreatedAt: now,
			UpdatedAt: now,
		},
		ContractNumber:    "CTR-2024-001",
		StartDate:         startDate,
		EndDate:           &endDate,
		ScopeOfServices:   "Complete investigation as per mandate",
		PaymentTerms:      paymentTerms,
		Confidentiality:   confidentiality,
		TerminationClause: terminationClause,
		Signatures: []document.Signature{
			{
				Name:      "John Doe",
				Role:      "Client",
				Signature: "client-signature-base64",
				SignedAt:  now,
			},
		},
	}
}

// NewTestContractWithCaseID creates a Contract with specific case and client IDs
func NewTestContractWithCaseID(t *testing.T, caseID, clientID uuid.UUID) *document.Contract {
	t.Helper()

	contract := NewTestContract(t)
	contract.CaseID = caseID
	contract.ClientID = clientID
	return contract
}

// NewTestActiveContract creates a Contract that is active
func NewTestActiveContract(t *testing.T) *document.Contract {
	t.Helper()

	contract := NewTestContract(t)
	contract.Status = document.DocumentStatusActive

	// Add investigator signature
	contract.Signatures = append(contract.Signatures, document.Signature{
		Name:      "Jane Smith",
		Role:      "Lead Investigator",
		Signature: "investigator-signature-base64",
		SignedAt:  time.Now(),
	})

	return contract
}

// NewTestInvoice creates a basic Invoice for testing
func NewTestInvoice(t *testing.T) *document.Invoice {
	t.Helper()

	now := time.Now()
	dueDate := now.AddDate(0, 1, 0) // 1 month from now
	notes := "Payment due within 30 days"
	currency := "USD"
	paymentTerms := "Net 30 days"

	return &document.Invoice{
		DocumentBase: document.DocumentBase{
			ID:        uuid.New(),
			CaseID:    uuid.New(),
			ClientID:  uuid.New(),
			Status:    document.DocumentStatusDraft,
			CreatedAt: now,
			UpdatedAt: now,
		},
		InvoiceNumber:    "INV-2024-001",
		IssueDate:        now,
		DueDate:          dueDate,
		LineItems: []document.InvoiceItem{
			{
				Description: "Investigation Services",
				Quantity:    40,
				UnitPrice:   150.00,
				Subtotal:    6000.00,
			},
			{
				Description: "Travel Expenses",
				Quantity:    1,
				UnitPrice:   500.00,
				Subtotal:    500.00,
			},
		},
		TotalAmount:       6500.00,
		TaxRate:           0.10,
		TaxAmount:         650.00,
		Notes:             &notes,
		PaymentStatus:     document.PaymentStatusUnpaid,
		PaidAt:            nil,
		PaidAmount:        nil,
		PaymentMethod:     nil,
		PaymentReference:  nil,
		LinkedContractID:  nil,
		Currency:          &currency,
		PaymentTerms:      &paymentTerms,
		LateFee:           nil,
		LateFeeRate:       nil,
	}
}

// NewTestInvoiceWithCaseID creates an Invoice with specific case and client IDs
func NewTestInvoiceWithCaseID(t *testing.T, caseID, clientID uuid.UUID) *document.Invoice {
	t.Helper()

	invoice := NewTestInvoice(t)
	invoice.CaseID = caseID
	invoice.ClientID = clientID
	return invoice
}

// NewTestPaidInvoice creates an Invoice that has been paid
func NewTestPaidInvoice(t *testing.T) *document.Invoice {
	t.Helper()

	invoice := NewTestInvoice(t)
	invoice.Status = document.DocumentStatusArchived
	invoice.PaymentStatus = document.PaymentStatusPaid
	paidAt := time.Now()
	paidAmount := 7150.00 // Total amount + tax
	paymentMethod := "Bank Transfer"
	paymentReference := "PAY-2024-001"

	invoice.PaidAt = &paidAt
	invoice.PaidAmount = &paidAmount
	invoice.PaymentMethod = &paymentMethod
	invoice.PaymentReference = &paymentReference

	return invoice
}

// NewTestDocumentVersion creates a DocumentVersion for testing
func NewTestDocumentVersion(t *testing.T) *document.DocumentVersion {
	t.Helper()

	version := 1
	authorID := uuid.New()
	data := `{"id":"test-id","status":"draft"}`
	reason := "Initial version"

	return &document.DocumentVersion{
		ID:         uuid.New(),
		DocumentID: uuid.New(),
		DocType:    document.DocumentTypeEstimate,
		Version:    version,
		AuthorID:   authorID,
		Data:       []byte(data),
		Reason:     reason,
		CreatedAt:  time.Now(),
	}
}

// NewTestDocumentVersionForDocument creates a DocumentVersion for a specific document
func NewTestDocumentVersionForDocument(t *testing.T, documentID uuid.UUID, docType document.DocumentType, version int) *document.DocumentVersion {
	t.Helper()

	docVersion := NewTestDocumentVersion(t)
	docVersion.DocumentID = documentID
	docVersion.DocType = docType
	docVersion.Version = version
	docVersion.Reason = fmt.Sprintf("Version %d of document", version)

	return docVersion
}

// CreateDocumentWorkflow creates a complete document workflow:
// Estimate -> Mandate -> Contract -> Invoice
func CreateDocumentWorkflow(t *testing.T, caseID, clientID uuid.UUID) (*document.Estimate, *document.Mandate, *document.Contract, *document.Invoice) {
	t.Helper()

	// Create estimate
	estimate := NewTestEstimateWithCaseID(t, caseID, clientID)
	estimate.Status = document.DocumentStatusSigned
	estimate.Accepted = true
	estimate.AcceptedAt = &[]time.Time{time.Now()}[0]
	estimate.AcceptedBy = &[]uuid.UUID{uuid.New()}[0]

	// Create mandate linked to estimate
	mandate := NewTestMandateWithLinkedEstimate(t, estimate.ID)
	mandate.CaseID = caseID
	mandate.ClientID = clientID
	mandate.Status = document.DocumentStatusSigned
	mandate.ClientSignature = &document.Signature{
		Name:      "John Doe",
		Role:      "Client",
		Signature: "client-sig",
		SignedAt:  time.Now(),
	}
	mandate.InvestigatorSignature = &document.Signature{
		Name:      "Jane Smith",
		Role:      "Investigator",
		Signature: "inv-sig",
		SignedAt:  time.Now(),
	}

	// Create contract from mandate
	contract := NewTestContractWithCaseID(t, caseID, clientID)
	contract.Status = document.DocumentStatusActive
	contract.Signatures = []document.Signature{
		*mandate.ClientSignature,
		*mandate.InvestigatorSignature,
	}

	// Create invoice from contract
	invoice := NewTestInvoiceWithCaseID(t, caseID, clientID)
	invoice.Status = document.DocumentStatusSent

	return estimate, mandate, contract, invoice
}