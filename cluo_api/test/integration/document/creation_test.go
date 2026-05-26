package document

import (
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domain "github.com/hengadev/cluo_api/internal/domain/document"
	investigation "github.com/hengadev/cluo_api/internal/domain/investigation"
	client "github.com/hengadev/cluo_api/internal/domain/client"
	"github.com/hengadev/cluo_api/internal/application/document"
)

// ---------------------------------------------------------------------------
// Mocks
// ---------------------------------------------------------------------------

type mockDocumentRepo struct {
	estimates map[string]*domain.EstimateEncx
	mandates  map[string]*domain.MandateEncx
	contracts map[string]*domain.ContractEncx
	invoices  map[string]*domain.InvoiceEncx
}

func newMockDocumentRepo() *mockDocumentRepo {
	return &mockDocumentRepo{
		estimates: make(map[string]*domain.EstimateEncx),
		mandates:  make(map[string]*domain.MandateEncx),
		contracts: make(map[string]*domain.ContractEncx),
		invoices:  make(map[string]*domain.InvoiceEncx),
	}
}

func (m *mockDocumentRepo) Create(_ context.Context, _ domain.Documentable) error { return nil }
func (m *mockDocumentRepo) GetByID(_ context.Context, _ string, _ domain.DocumentType) (domain.Documentable, error) {
	return nil, fmt.Errorf("not implemented")
}
func (m *mockDocumentRepo) Update(_ context.Context, _ domain.Documentable) error { return nil }
func (m *mockDocumentRepo) Delete(_ context.Context, _ string, _ domain.DocumentType) error {
	return nil
}
func (m *mockDocumentRepo) List(_ context.Context, _ domain.DocumentFilter, _ domain.Pagination) ([]domain.DocumentSummary, int, error) {
	return nil, 0, nil
}
func (m *mockDocumentRepo) GetLinkedDocuments(_ context.Context, _ string, _ domain.DocumentType) ([]domain.Documentable, error) {
	return nil, nil
}
func (m *mockDocumentRepo) GetFirstByCaseAndType(_ context.Context, _ string, _ domain.DocumentType) (domain.Documentable, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *mockDocumentRepo) CreateEstimate(_ context.Context, e *domain.EstimateEncx) error {
	m.estimates[e.ID.String()] = e
	return nil
}
func (m *mockDocumentRepo) GetEstimateByID(_ context.Context, id string) (*domain.EstimateEncx, error) {
	e, ok := m.estimates[id]
	if !ok {
		return nil, fmt.Errorf("estimate %s not found", id)
	}
	return e, nil
}
func (m *mockDocumentRepo) UpdateEstimate(_ context.Context, e *domain.EstimateEncx) error {
	m.estimates[e.ID.String()] = e
	return nil
}
func (m *mockDocumentRepo) DeleteEstimate(_ context.Context, id string) error {
	delete(m.estimates, id)
	return nil
}
func (m *mockDocumentRepo) ListEstimatesByCase(_ context.Context, _ string, _ domain.Pagination) ([]*domain.EstimateEncx, int, error) {
	return nil, 0, nil
}

func (m *mockDocumentRepo) CreateMandate(_ context.Context, man *domain.MandateEncx) error {
	m.mandates[man.ID.String()] = man
	return nil
}
func (m *mockDocumentRepo) GetMandateByID(_ context.Context, id string) (*domain.MandateEncx, error) {
	man, ok := m.mandates[id]
	if !ok {
		return nil, fmt.Errorf("mandate %s not found", id)
	}
	return man, nil
}
func (m *mockDocumentRepo) UpdateMandate(_ context.Context, man *domain.MandateEncx) error {
	m.mandates[man.ID.String()] = man
	return nil
}
func (m *mockDocumentRepo) DeleteMandate(_ context.Context, id string) error {
	delete(m.mandates, id)
	return nil
}
func (m *mockDocumentRepo) ListMandatesByCase(_ context.Context, _ string, _ domain.Pagination) ([]*domain.MandateEncx, int, error) {
	return nil, 0, nil
}

func (m *mockDocumentRepo) CreateContract(_ context.Context, c *domain.ContractEncx) error {
	m.contracts[c.ID.String()] = c
	return nil
}
func (m *mockDocumentRepo) GetContractByID(_ context.Context, id string) (*domain.ContractEncx, error) {
	c, ok := m.contracts[id]
	if !ok {
		return nil, fmt.Errorf("contract %s not found", id)
	}
	return c, nil
}
func (m *mockDocumentRepo) UpdateContract(_ context.Context, c *domain.ContractEncx) error {
	m.contracts[c.ID.String()] = c
	return nil
}
func (m *mockDocumentRepo) DeleteContract(_ context.Context, id string) error {
	delete(m.contracts, id)
	return nil
}
func (m *mockDocumentRepo) ListContractsByCase(_ context.Context, _ string, _ domain.Pagination) ([]*domain.ContractEncx, int, error) {
	return nil, 0, nil
}

func (m *mockDocumentRepo) CreateInvoice(_ context.Context, inv *domain.InvoiceEncx) error {
	m.invoices[inv.ID.String()] = inv
	return nil
}
func (m *mockDocumentRepo) GetInvoiceByID(_ context.Context, id string) (*domain.InvoiceEncx, error) {
	inv, ok := m.invoices[id]
	if !ok {
		return nil, fmt.Errorf("invoice %s not found", id)
	}
	return inv, nil
}
func (m *mockDocumentRepo) UpdateInvoice(_ context.Context, inv *domain.InvoiceEncx) error {
	m.invoices[inv.ID.String()] = inv
	return nil
}
func (m *mockDocumentRepo) DeleteInvoice(_ context.Context, id string) error {
	delete(m.invoices, id)
	return nil
}
func (m *mockDocumentRepo) ListInvoicesByCase(_ context.Context, _ string, _ domain.Pagination) ([]*domain.InvoiceEncx, int, error) {
	return nil, 0, nil
}
func (m *mockDocumentRepo) ListOverdueInvoices(_ context.Context, _ domain.Pagination) ([]*domain.InvoiceEncx, int, error) {
	return nil, 0, nil
}

type mockVersionRepo struct{}

func (m *mockVersionRepo) CreateVersion(_ context.Context, _ *domain.DocumentVersion) error {
	return nil
}
func (m *mockVersionRepo) GetDocumentHistory(_ context.Context, _ string, _ domain.DocumentType, _ domain.Pagination) ([]*domain.DocumentVersion, int, error) {
	return nil, 0, nil
}
func (m *mockVersionRepo) GetVersion(_ context.Context, _ string, _ domain.DocumentType, _ int) (*domain.DocumentVersion, error) {
	return nil, nil
}
func (m *mockVersionRepo) DeleteVersions(_ context.Context, _ string, _ domain.DocumentType) error {
	return nil
}

type mockCaseRepo struct{}

func (m *mockCaseRepo) CreateCase(_ context.Context, _ *investigation.InvestigationEncx) error {
	return nil
}
func (m *mockCaseRepo) GetCaseByID(_ context.Context, _ uuid.UUID) (*investigation.InvestigationEncx, error) {
	return nil, nil
}
func (m *mockCaseRepo) UpdateCase(_ context.Context, _ *investigation.InvestigationEncx) error {
	return nil
}
func (m *mockCaseRepo) DeleteCase(_ context.Context, _ uuid.UUID) error { return nil }
func (m *mockCaseRepo) List(_ context.Context, _ investigation.Filter, _ investigation.Pagination) ([]*investigation.InvestigationEncx, int, error) {
	return nil, 0, nil
}
func (m *mockCaseRepo) ListByClient(_ context.Context, _ uuid.UUID, _ investigation.Pagination) ([]*investigation.InvestigationEncx, int, error) {
	return nil, 0, nil
}
func (m *mockCaseRepo) ExistsCase(_ context.Context, _ uuid.UUID) (bool, error) { return true, nil }

type mockClientRepo struct{}

func (m *mockClientRepo) CreateContact(_ context.Context, _ *client.ContactEncx) error { return nil }
func (m *mockClientRepo) ExistsClient(_ context.Context, _ uuid.UUID) (bool, error)    { return true, nil }
func (m *mockClientRepo) DeleteContact(_ context.Context, _ uuid.UUID) error           { return nil }
func (m *mockClientRepo) UpdateContact(_ context.Context, _ *client.ContactEncx) error { return nil }
func (m *mockClientRepo) GetContactByID(_ context.Context, _ uuid.UUID) (*client.ContactEncx, error) {
	return nil, nil
}
func (m *mockClientRepo) GetAllContactsByClientID(_ context.Context, _ uuid.UUID) ([]*client.ContactEncx, error) {
	return nil, nil
}
func (m *mockClientRepo) GetContactIDsForClient(_ context.Context, _ uuid.UUID) ([]uuid.UUID, error) {
	return nil, nil
}
func (m *mockClientRepo) CreateClient(_ context.Context, _ *client.ClientEncx) error { return nil }
func (m *mockClientRepo) ExistsContact(_ context.Context, _ uuid.UUID) (bool, error) {
	return true, nil
}
func (m *mockClientRepo) DeleteClient(_ context.Context, _ uuid.UUID) error           { return nil }
func (m *mockClientRepo) UpdateClient(_ context.Context, _ *client.ClientEncx) error  { return nil }
func (m *mockClientRepo) GetClientByID(_ context.Context, _ uuid.UUID) (*client.ClientEncx, error) {
	return nil, nil
}
func (m *mockClientRepo) GetAllClients(_ context.Context) ([]*client.ClientEncx, error) {
	return nil, nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func stringPtr(s string) *string { return &s }

// noOpEmailSvc is a test double for ports.EmailService.
type noOpEmailSvc struct{}

func (noOpEmailSvc) Send(_ context.Context, _, _, _ string) error { return nil }

// ---------------------------------------------------------------------------
// Tests
// TODO: These tests need a test crypto service (encx.NewTestCrypto) wired in
// before they can run. Skipped until the setup is added.
// ---------------------------------------------------------------------------

func TestEstimateCreationAndAcceptance(t *testing.T) {
	t.Skip("TODO: wire encx.NewTestCrypto — service requires encryption")

	ctx := context.Background()
	repo := newMockDocumentRepo()
	svc := document.New(repo, &mockVersionRepo{}, &mockCaseRepo{}, &mockClientRepo{}, nil, noOpEmailSvc{}, slog.Default())

	caseID := uuid.New()
	clientID := uuid.New()
	lineItems := []domain.EstimateItem{
		{Description: "Surveillance services (8 hours)", Quantity: 8, UnitPrice: 50.00, Subtotal: 400.00},
		{Description: "Travel expenses", Quantity: 1, UnitPrice: 25.00, Subtotal: 25.00},
	}

	estimate := domain.NewEstimate(caseID, clientID, "EST-2025-001", lineItems)
	estimate.EstimateNumber = "EST-2025-001"

	createdEstimate, err := svc.CreateEstimate(ctx, estimate)
	require.NoError(t, err)
	require.NotNil(t, createdEstimate)
	assert.Equal(t, caseID, createdEstimate.CaseID)
	assert.Equal(t, clientID, createdEstimate.ClientID)
	assert.Equal(t, "EST-2025-001", createdEstimate.EstimateNumber)
	assert.Equal(t, 425.00, createdEstimate.EstimatedTotal)
	assert.False(t, createdEstimate.Accepted)

	acceptedBy := uuid.New().String()
	mandate, err := svc.AcceptEstimate(ctx, createdEstimate.ID.String(), acceptedBy)
	require.NoError(t, err)
	require.NotNil(t, mandate)
	assert.Equal(t, caseID, mandate.CaseID)
	assert.Equal(t, clientID, mandate.ClientID)
	assert.NotNil(t, mandate.LinkedEstimateID)
	assert.Equal(t, createdEstimate.ID, *mandate.LinkedEstimateID)

	t.Logf("✅ Successfully created estimate %s and generated mandate %s",
		estimate.EstimateNumber, mandate.MandateNumber)
}

func TestMandateSigningFlow(t *testing.T) {
	t.Skip("TODO: wire encx.NewTestCrypto — service requires encryption")

	ctx := context.Background()
	repo := newMockDocumentRepo()
	svc := document.New(repo, &mockVersionRepo{}, &mockCaseRepo{}, &mockClientRepo{}, nil, noOpEmailSvc{}, slog.Default())

	caseID := uuid.New()
	clientID := uuid.New()

	mandate := domain.NewMandate(
		caseID, clientID, "MND-2025-001",
		"Foot surveillance and evidence gathering in downtown area",
		"Standard investigation terms and conditions apply",
		time.Now(),
	)
	mandate.MandateNumber = "MND-2025-001"

	createdMandate, err := svc.CreateMandate(ctx, mandate)
	require.NoError(t, err)
	require.NotNil(t, createdMandate)
	assert.Equal(t, domain.DocumentStatusDraft, createdMandate.Status)

	signReq := &domain.SignDocumentRequest{
		SignerName:       "John Doe",
		SignerRole:       "client",
		Method:           "e-sign",
		SignatureFileURL: "https://example.com/signature.pdf",
	}

	signedMandate, err := svc.SignMandate(ctx, createdMandate.ID.String(), signReq)
	require.NoError(t, err)
	require.NotNil(t, signedMandate)
	assert.NotNil(t, signedMandate.ClientSignature)
	assert.Equal(t, "John Doe", signedMandate.ClientSignature.Name)

	activatedMandate, err := svc.ActivateMandate(ctx, signedMandate.ID.String())
	require.NoError(t, err)
	require.NotNil(t, activatedMandate)
	assert.Equal(t, domain.DocumentStatusActive, activatedMandate.Status)

	t.Logf("✅ Successfully created mandate %s, signed by client, and activated",
		mandate.MandateNumber)
}

func TestContractCreationFromMandate(t *testing.T) {
	t.Skip("TODO: wire encx.NewTestCrypto — service requires encryption")

	ctx := context.Background()
	repo := newMockDocumentRepo()
	svc := document.New(repo, &mockVersionRepo{}, &mockCaseRepo{}, &mockClientRepo{}, nil, noOpEmailSvc{}, slog.Default())

	caseID := uuid.New()
	clientID := uuid.New()

	mandate := domain.NewMandate(
		caseID, clientID, "MND-2025-002",
		"Background investigation services",
		"Standard investigation terms and conditions apply",
		time.Now(),
	)
	mandate.MandateNumber = "MND-2025-002"
	mandate.SetStatus(domain.DocumentStatusActive)

	createdMandate, err := svc.CreateMandate(ctx, mandate)
	require.NoError(t, err)

	contract := domain.NewContract(
		caseID, clientID, "CNT-2025-002",
		"Comprehensive background investigation services",
		"50% upon signing, 50% upon completion",
		"All information kept strictly confidential",
		"Either party may terminate with 30 days written notice",
		time.Now(),
	)
	contract.ContractNumber = "CNT-2025-002"

	createdContract, err := svc.CreateContractFromMandate(ctx, createdMandate.ID.String(), contract)
	require.NoError(t, err)
	require.NotNil(t, createdContract)
	assert.Equal(t, caseID, createdContract.CaseID)
	assert.NotNil(t, createdContract.LinkedMandateID)
	assert.Equal(t, createdMandate.ID, *createdContract.LinkedMandateID)

	t.Logf("✅ Successfully created contract %s from mandate %s",
		contract.ContractNumber, mandate.MandateNumber)
}

func TestInvoiceCreationFromContract(t *testing.T) {
	t.Skip("TODO: wire encx.NewTestCrypto — service requires encryption")

	ctx := context.Background()
	repo := newMockDocumentRepo()
	svc := document.New(repo, &mockVersionRepo{}, &mockCaseRepo{}, &mockClientRepo{}, nil, noOpEmailSvc{}, slog.Default())

	caseID := uuid.New()
	clientID := uuid.New()

	contract := domain.NewContract(
		caseID, clientID, "CNT-2025-003",
		"Private investigation services",
		"Net 30 days", "Standard confidentiality agreement",
		"30 days notice required for termination",
		time.Now(),
	)
	contract.ContractNumber = "CNT-2025-003"
	contract.SetStatus(domain.DocumentStatusActive)

	createdContract, err := svc.CreateContract(ctx, contract)
	require.NoError(t, err)

	lineItems := []domain.InvoiceItem{
		{Description: "Investigation services (40 hours)", Quantity: 40, UnitPrice: 75.00, Subtotal: 3000.00},
		{Description: "Travel expenses", Quantity: 1, UnitPrice: 150.00, Subtotal: 150.00},
	}
	invoice := domain.NewInvoice(caseID, clientID, "INV-2025-003", lineItems, 0.08, time.Now().AddDate(0, 0, 30))
	invoice.InvoiceNumber = "INV-2025-003"

	createdInvoice, err := svc.CreateInvoiceFromContract(ctx, createdContract.ID.String(), invoice)
	require.NoError(t, err)
	require.NotNil(t, createdInvoice)
	assert.Equal(t, caseID, createdInvoice.CaseID)
	assert.NotNil(t, createdInvoice.LinkedContractID)
	assert.Equal(t, createdContract.ID, *createdInvoice.LinkedContractID)

	t.Logf("✅ Successfully created invoice %s from contract %s",
		invoice.InvoiceNumber, contract.ContractNumber)
}

func TestCompleteDocumentWorkflow(t *testing.T) {
	t.Skip("TODO: wire encx.NewTestCrypto — service requires encryption")

	ctx := context.Background()
	repo := newMockDocumentRepo()
	svc := document.New(repo, &mockVersionRepo{}, &mockCaseRepo{}, &mockClientRepo{}, nil, noOpEmailSvc{}, slog.Default())

	caseID := uuid.New()
	clientID := uuid.New()

	t.Log("🔄 Starting complete document workflow...")

	estimateLineItems := []domain.EstimateItem{
		{Description: "Investigation services", Quantity: 20, UnitPrice: 100.00, Subtotal: 2000.00},
	}
	estimate := domain.NewEstimate(caseID, clientID, "EST-2025-WF-001", estimateLineItems)
	estimate.EstimateNumber = "EST-2025-WF-001"

	createdEstimate, err := svc.CreateEstimate(ctx, estimate)
	require.NoError(t, err)
	t.Logf("✅ Step 1: Created estimate %s", createdEstimate.EstimateNumber)

	acceptedBy := uuid.New().String()
	mandate, err := svc.AcceptEstimate(ctx, createdEstimate.ID.String(), acceptedBy)
	require.NoError(t, err)
	t.Logf("✅ Step 2: Accepted estimate and created mandate %s", mandate.MandateNumber)

	signReq := &domain.SignDocumentRequest{
		SignerName: "Client Name", SignerRole: "client",
		Method: "e-sign", SignatureFileURL: "https://example.com/client-signature.pdf",
	}
	signedMandate, err := svc.SignMandate(ctx, mandate.ID.String(), signReq)
	require.NoError(t, err)
	t.Log("✅ Step 3: Signed mandate by client")

	activatedMandate, err := svc.ActivateMandate(ctx, signedMandate.ID.String())
	require.NoError(t, err)
	t.Log("✅ Step 4: Activated mandate")

	contract := domain.NewContract(
		caseID, clientID, "CNT-2025-WF-001",
		"Private investigation services as outlined in estimate",
		"Net 15 days", "Standard confidentiality",
		"30 days notice for termination", time.Now(),
	)
	contract.ContractNumber = "CNT-2025-WF-001"

	createdContract, err := svc.CreateContractFromMandate(ctx, activatedMandate.ID.String(), contract)
	require.NoError(t, err)
	t.Logf("✅ Step 5: Created contract %s from mandate", createdContract.ContractNumber)

	contractSignReq := &domain.SignDocumentRequest{
		SignerName: "Investigator Name", SignerRole: "investigator",
		Method: "e-sign", SignatureFileURL: "https://example.com/investigator-signature.pdf",
	}
	signedContract, err := svc.SignContract(ctx, createdContract.ID.String(), contractSignReq)
	require.NoError(t, err)
	t.Log("✅ Step 6: Signed contract by investigator")

	invoiceLineItems := []domain.InvoiceItem{
		{Description: "Investigation services", Quantity: 20, UnitPrice: 100.00, Subtotal: 2000.00},
	}
	invoice := domain.NewInvoice(caseID, clientID, "INV-2025-WF-001", invoiceLineItems, 0.08, time.Now().AddDate(0, 0, 15))
	invoice.InvoiceNumber = "INV-2025-WF-001"

	createdInvoice, err := svc.CreateInvoiceFromContract(ctx, signedContract.ID.String(), invoice)
	require.NoError(t, err)
	t.Logf("✅ Step 7: Created invoice %s from contract", createdInvoice.InvoiceNumber)

	paymentReq := &domain.PaymentRequest{
		Amount: 2160.00, PaymentMethod: "wire_transfer",
		Notes: stringPtr("Payment for investigation services"),
	}
	paidInvoice, err := svc.ProcessPayment(ctx, createdInvoice.ID.String(), paymentReq)
	require.NoError(t, err)
	assert.Equal(t, domain.PaymentStatusPaid, paidInvoice.PaymentStatus)
	t.Log("✅ Step 8: Processed payment for invoice")

	workflow, err := svc.GetDocumentWorkflow(ctx, caseID.String())
	require.NoError(t, err)
	require.NotNil(t, workflow)
	assert.NotNil(t, workflow.Estimate)
	assert.NotNil(t, workflow.Mandate)
	assert.NotNil(t, workflow.Contract)
	assert.NotNil(t, workflow.Invoice)
	t.Logf("✅ Step 9: Retrieved workflow summary - estimate=%v mandate=%v contract=%v invoice=%v",
		workflow.Estimate != nil, workflow.Mandate != nil, workflow.Contract != nil, workflow.Invoice != nil)

	t.Log("🎉 Complete document workflow test passed successfully!")
}
