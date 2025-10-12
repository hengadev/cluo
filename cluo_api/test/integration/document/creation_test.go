package document

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hengadev/cluo_api/internal/domain"
	"github.com/hengadev/cluo_api/internal/application/document"
)

// MockDocumentRepository implements ports.DocumentRepository for testing
type MockDocumentRepository struct {
	estimates map[uuid.UUID]*domain.Estimate
	mandates  map[uuid.UUID]*domain.Mandate
	contracts map[uuid.UUID]*domain.Contract
	invoices  map[uuid.UUID]*domain.Invoice
}

func NewMockDocumentRepository() *MockDocumentRepository {
	return &MockDocumentRepository{
		estimates: make(map[uuid.UUID]*domain.Estimate),
		mandates:  make(map[uuid.UUID]*domain.Mandate),
		contracts: make(map[uuid.UUID]*domain.Contract),
		invoices:  make(map[uuid.UUID]*domain.Invoice),
	}
}

func (m *MockDocumentRepository) CreateEstimate(ctx context.Context, estimate *domain.Estimate) error {
	m.estimates[estimate.ID] = estimate
	return nil
}

func (m *MockDocumentRepository) GetEstimateByID(ctx context.Context, id string) (*domain.Estimate, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	estimate, exists := m.estimates[uid]
	if !exists {
		return nil, assert.AnError
	}
	return estimate, nil
}

func (m *MockDocumentRepository) UpdateEstimate(ctx context.Context, estimate *domain.Estimate) error {
	m.estimates[estimate.ID] = estimate
	return nil
}

func (m *MockDocumentRepository) CreateMandate(ctx context.Context, mandate *domain.Mandate) error {
	m.mandates[mandate.ID] = mandate
	return nil
}

func (m *MockDocumentRepository) GetMandateByID(ctx context.Context, id string) (*domain.Mandate, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	mandate, exists := m.mandates[uid]
	if !exists {
		return nil, assert.AnError
	}
	return mandate, nil
}

func (m *MockDocumentRepository) UpdateMandate(ctx context.Context, mandate *domain.Mandate) error {
	m.mandates[mandate.ID] = mandate
	return nil
}

func (m *MockDocumentRepository) CreateContract(ctx context.Context, contract *domain.Contract) error {
	m.contracts[contract.ID] = contract
	return nil
}

func (m *MockDocumentRepository) GetContractByID(ctx context.Context, id string) (*domain.Contract, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	contract, exists := m.contracts[uid]
	if !exists {
		return nil, assert.AnError
	}
	return contract, nil
}

func (m *MockDocumentRepository) CreateInvoice(ctx context.Context, invoice *domain.Invoice) error {
	m.invoices[invoice.ID] = invoice
	return nil
}

func (m *MockDocumentRepository) GetInvoiceByID(ctx context.Context, id string) (*domain.Invoice, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	invoice, exists := m.invoices[uid]
	if !exists {
		return nil, assert.AnError
	}
	return invoice, nil
}

// MockDocumentVersionRepository implements ports.DocumentVersionRepository for testing
type MockDocumentVersionRepository struct {
	versions map[string][]*domain.DocumentVersion // key: "documentId:docType"
}

func NewMockDocumentVersionRepository() *MockDocumentVersionRepository {
	return &MockDocumentVersionRepository{
		versions: make(map[string][]*domain.DocumentVersion),
	}
}

func (m *MockDocumentVersionRepository) CreateVersion(ctx context.Context, version *domain.DocumentVersion) error {
	key := version.DocumentID.String() + ":" + string(version.DocType)
	m.versions[key] = append(m.versions[key], version)
	return nil
}

func (m *MockDocumentVersionRepository) GetDocumentHistory(ctx context.Context, documentID string, docType domain.DocumentType, pagination domain.Pagination) ([]*domain.DocumentVersion, int, error) {
	key := documentID + ":" + string(docType)
	versions, exists := m.versions[key]
	if !exists {
		return nil, 0, nil
	}
	return versions, len(versions), nil
}

// Test Estimate Creation and Acceptance Flow
func TestEstimateCreationAndAcceptance(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockRepo := NewMockDocumentRepository()
	mockVersionRepo := NewMockDocumentVersionRepository()
	service := document.New(mockRepo, mockVersionRepo)

	// Create test data
	caseID := uuid.New()
	clientID := uuid.New()
	lineItems := []domain.EstimateItem{
		{
			Description: "Surveillance services (8 hours)",
			Quantity:    8,
			UnitPrice:   50.00,
			Subtotal:    400.00,
		},
		{
			Description: "Travel expenses",
			Quantity:    1,
			UnitPrice:   25.00,
			Subtotal:    25.00,
		},
	}

	// Test estimate creation
	estimate := domain.NewEstimate(caseID, clientID, "EST-2025-001", lineItems)
	estimate.EstimateNumber = "EST-2025-001"

	createdEstimate, err := service.CreateEstimate(ctx, estimate)
	require.NoError(t, err)
	require.NotNil(t, createdEstimate)
	assert.Equal(t, caseID, createdEstimate.CaseID)
	assert.Equal(t, clientID, createdEstimate.ClientID)
	assert.Equal(t, "EST-2025-001", createdEstimate.EstimateNumber)
	assert.Equal(t, 425.00, createdEstimate.EstimatedTotal) // 400 + 25
	assert.False(t, createdEstimate.Accepted)

	// Test estimate acceptance
	acceptedBy := uuid.New().String()
	mandate, err := service.AcceptEstimate(ctx, createdEstimate.ID.String(), acceptedBy)
	require.NoError(t, err)
	require.NotNil(t, mandate)
	assert.Equal(t, caseID, mandate.CaseID)
	assert.Equal(t, clientID, mandate.ClientID)
	assert.NotNil(t, mandate.LinkedEstimateID)
	assert.Equal(t, createdEstimate.ID, *mandate.LinkedEstimateID)

	// Verify estimate was updated
	updatedEstimate, err := service.GetDocument(ctx, createdEstimate.ID.String(), domain.DocumentTypeEstimate)
	require.NoError(t, err)
	estimateTyped, ok := updatedEstimate.(*domain.Estimate)
	require.True(t, ok)
	assert.True(t, estimateTyped.Accepted)
	assert.NotNil(t, estimateTyped.AcceptedAt)
	assert.NotNil(t, estimateTyped.AcceptedBy)

	t.Logf("✅ Successfully created estimate %s and generated mandate %s",
		estimate.EstimateNumber, mandate.MandateNumber)
}

// Test Mandate Signing Flow
func TestMandateSigningFlow(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockRepo := NewMockDocumentRepository()
	mockVersionRepo := NewMockDocumentVersionRepository()
	service := document.New(mockRepo, mockVersionRepo)

	// Create test data
	caseID := uuid.New()
	clientID := uuid.New()

	// Create mandate
	mandate := domain.NewMandate(
		caseID,
		clientID,
		"MND-2025-001",
		"Foot surveillance and evidence gathering in downtown area",
		"Standard investigation terms and conditions apply",
		time.Now(),
	)
	mandate.MandateNumber = "MND-2025-001"

	createdMandate, err := service.CreateMandate(ctx, mandate)
	require.NoError(t, err)
	require.NotNil(t, createdMandate)
	assert.Equal(t, domain.DocumentStatusDraft, createdMandate.Status)

	// Sign mandate
	signReq := &domain.SignDocumentRequest{
		SignerName:       "John Doe",
		SignerRole:       "client",
		Method:           "e-sign",
		SignatureFileURL: "https://example.com/signature.pdf",
	}

	signedMandate, err := service.SignMandate(ctx, createdMandate.ID.String(), signReq)
	require.NoError(t, err)
	require.NotNil(t, signedMandate)
	assert.NotNil(t, signedMandate.ClientSignature)
	assert.Equal(t, "John Doe", signedMandate.ClientSignature.Name)
	assert.Equal(t, "client", signedMandate.ClientSignature.Role)

	// Activate mandate
	activatedMandate, err := service.ActivateMandate(ctx, signedMandate.ID.String())
	require.NoError(t, err)
	require.NotNil(t, activatedMandate)
	assert.Equal(t, domain.DocumentStatusActive, activatedMandate.Status)

	t.Logf("✅ Successfully created mandate %s, signed by client, and activated",
		mandate.MandateNumber)
}

// Test Contract Creation from Mandate
func TestContractCreationFromMandate(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockRepo := NewMockDocumentRepository()
	mockVersionRepo := NewMockDocumentVersionRepository()
	service := document.New(mockRepo, mockVersionRepo)

	// Create test data
	caseID := uuid.New()
	clientID := uuid.New()

	// Create active mandate first
	mandate := domain.NewMandate(
		caseID,
		clientID,
		"MND-2025-002",
		"Background investigation services",
		"Standard investigation terms and conditions apply",
		time.Now(),
	)
	mandate.MandateNumber = "MND-2025-002"
	mandate.SetStatus(domain.DocumentStatusActive)

	createdMandate, err := service.CreateMandate(ctx, mandate)
	require.NoError(t, err)

	// Create contract from mandate
	contract := domain.NewContract(
		caseID,
		clientID,
		"CNT-2025-002",
		"Comprehensive background investigation services including employment verification, criminal record check, and reference verification",
		"50% upon signing, 50% upon completion",
		"All information obtained during investigation will be kept strictly confidential",
		"Either party may terminate with 30 days written notice",
		time.Now(),
	)
	contract.ContractNumber = "CNT-2025-002"

	createdContract, err := service.CreateContractFromMandate(ctx, createdMandate.ID.String(), contract)
	require.NoError(t, err)
	require.NotNil(t, createdContract)
	assert.Equal(t, caseID, createdContract.CaseID)
	assert.Equal(t, clientID, createdContract.ClientID)
	assert.NotNil(t, createdContract.LinkedMandateID)
	assert.Equal(t, createdMandate.ID, *createdContract.LinkedMandateID)

	t.Logf("✅ Successfully created contract %s from mandate %s",
		contract.ContractNumber, mandate.MandateNumber)
}

// Test Invoice Creation from Contract
func TestInvoiceCreationFromContract(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockRepo := NewMockDocumentRepository()
	mockVersionRepo := NewMockDocumentVersionRepository()
	service := document.New(mockRepo, mockVersionRepo)

	// Create test data
	caseID := uuid.New()
	clientID := uuid.New()

	// Create active contract first
	contract := domain.NewContract(
		caseID,
		clientID,
		"CNT-2025-003",
		"Private investigation services",
		"Net 30 days",
		"Standard confidentiality agreement",
		"30 days notice required for termination",
		time.Now(),
	)
	contract.ContractNumber = "CNT-2025-003"
	contract.SetStatus(domain.DocumentStatusActive)

	createdContract, err := service.CreateContract(ctx, contract)
	require.NoError(t, err)

	// Create invoice from contract
	lineItems := []domain.InvoiceItem{
		{
			Description: "Investigation services (40 hours)",
			Quantity:    40,
			UnitPrice:   75.00,
			Subtotal:    3000.00,
		},
		{
			Description: "Travel expenses",
			Quantity:    1,
			UnitPrice:   150.00,
			Subtotal:    150.00,
		},
	}

	invoice := domain.NewInvoice(
		caseID,
		clientID,
		"INV-2025-003",
		lineItems,
		0.08, // 8% tax
		time.Now().AddDate(0, 0, 30), // Due in 30 days
	)
	invoice.InvoiceNumber = "INV-2025-003"

	createdInvoice, err := service.CreateInvoiceFromContract(ctx, createdContract.ID.String(), invoice)
	require.NoError(t, err)
	require.NotNil(t, createdInvoice)
	assert.Equal(t, caseID, createdInvoice.CaseID)
	assert.Equal(t, clientID, createdInvoice.ClientID)
	assert.NotNil(t, createdInvoice.LinkedContractID)
	assert.Equal(t, createdContract.ID, *createdInvoice.LinkedContractID)
	assert.Equal(t, 3150.00, createdInvoice.TotalAmount) // 3000 + 150
	assert.Equal(t, 252.00, createdInvoice.TaxAmount)   // 3150 * 0.08

	t.Logf("✅ Successfully created invoice %s from contract %s",
		invoice.InvoiceNumber, contract.ContractNumber)
}

// Test Complete Document Workflow
func TestCompleteDocumentWorkflow(t *testing.T) {
	// Setup
	ctx := context.Background()
	mockRepo := NewMockDocumentRepository()
	mockVersionRepo := NewMockDocumentVersionRepository()
	service := document.New(mockRepo, mockVersionRepo)

	// Create test data
	caseID := uuid.New()
	clientID := uuid.New()

	t.Log("🔄 Starting complete document workflow...")

	// 1. Create Estimate
	estimateLineItems := []domain.EstimateItem{
		{Description: "Investigation services", Quantity: 20, UnitPrice: 100.00, Subtotal: 2000.00},
	}
	estimate := domain.NewEstimate(caseID, clientID, "EST-2025-WF-001", estimateLineItems)
	estimate.EstimateNumber = "EST-2025-WF-001"

	createdEstimate, err := service.CreateEstimate(ctx, estimate)
	require.NoError(t, err)
	t.Logf("✅ Step 1: Created estimate %s", createdEstimate.EstimateNumber)

	// 2. Accept Estimate → Create Mandate
	acceptedBy := uuid.New().String()
	mandate, err := service.AcceptEstimate(ctx, createdEstimate.ID.String(), acceptedBy)
	require.NoError(t, err)
	t.Logf("✅ Step 2: Accepted estimate and created mandate %s", mandate.MandateNumber)

	// 3. Sign Mandate
	signReq := &domain.SignDocumentRequest{
		SignerName:       "Client Name",
		SignerRole:       "client",
		Method:           "e-sign",
		SignatureFileURL: "https://example.com/client-signature.pdf",
	}
	signedMandate, err := service.SignMandate(ctx, mandate.ID.String(), signReq)
	require.NoError(t, err)
	t.Logf("✅ Step 3: Signed mandate by client")

	// 4. Activate Mandate
	activatedMandate, err := service.ActivateMandate(ctx, signedMandate.ID.String())
	require.NoError(t, err)
	t.Logf("✅ Step 4: Activated mandate")

	// 5. Create Contract from Mandate
	contract := domain.NewContract(
		caseID,
		clientID,
		"CNT-2025-WF-001",
		"Private investigation services as outlined in estimate",
		"Net 15 days",
		"Standard confidentiality",
		"30 days notice for termination",
		time.Now(),
	)
	contract.ContractNumber = "CNT-2025-WF-001"

	createdContract, err := service.CreateContractFromMandate(ctx, activatedMandate.ID.String(), contract)
	require.NoError(t, err)
	t.Logf("✅ Step 5: Created contract %s from mandate", createdContract.ContractNumber)

	// 6. Sign Contract
	contractSignReq := &domain.SignDocumentRequest{
		SignerName:       "Investigator Name",
		SignerRole:       "investigator",
		Method:           "e-sign",
		SignatureFileURL: "https://example.com/investigator-signature.pdf",
	}
	signedContract, err := service.SignContract(ctx, createdContract.ID.String(), contractSignReq)
	require.NoError(t, err)
	t.Logf("✅ Step 6: Signed contract by investigator")

	// 7. Create Invoice from Contract
	invoiceLineItems := []domain.InvoiceItem{
		{Description: "Investigation services", Quantity: 20, UnitPrice: 100.00, Subtotal: 2000.00},
	}
	invoice := domain.NewInvoice(
		caseID,
		clientID,
		"INV-2025-WF-001",
		invoiceLineItems,
		0.08, // 8% tax
		time.Now().AddDate(0, 0, 15), // Due in 15 days
	)
	invoice.InvoiceNumber = "INV-2025-WF-001"

	createdInvoice, err := service.CreateInvoiceFromContract(ctx, signedContract.ID.String(), invoice)
	require.NoError(t, err)
	t.Logf("✅ Step 7: Created invoice %s from contract", createdInvoice.InvoiceNumber)

	// 8. Process Payment
	paymentReq := &domain.PaymentRequest{
		Amount:        2160.00, // Total amount + tax
		PaymentMethod: "wire_transfer",
		Notes:         stringPtr("Payment for investigation services"),
	}
	paidInvoice, err := service.ProcessPayment(ctx, createdInvoice.ID.String(), paymentReq)
	require.NoError(t, err)
	assert.Equal(t, domain.PaymentStatusPaid, paidInvoice.PaymentStatus)
	t.Logf("✅ Step 8: Processed payment for invoice")

	// 9. Get Document Workflow Summary
	workflow, err := service.GetDocumentWorkflow(ctx, caseID.String())
	require.NoError(t, err)
	assert.Len(t, workflow, 4) // estimate, mandate, contract, invoice
	t.Logf("✅ Step 9: Retrieved workflow summary - %d documents", len(workflow))

	t.Log("🎉 Complete document workflow test passed successfully!")
	t.Logf("📋 Workflow Summary:")
	t.Logf("   - Estimate: %s (Status: %s, Accepted: %t)",
		createdEstimate.EstimateNumber, createdEstimate.Status, createdEstimate.Accepted)
	t.Logf("   - Mandate: %s (Status: %s)",
		mandate.MandateNumber, activatedMandate.Status)
	t.Logf("   - Contract: %s (Status: %s)",
		createdContract.ContractNumber, signedContract.Status)
	t.Logf("   - Invoice: %s (Status: %s, Payment: %s)",
		createdInvoice.InvoiceNumber, paidInvoice.Status, paidInvoice.PaymentStatus)
}