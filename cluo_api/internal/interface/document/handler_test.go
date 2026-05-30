package document

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/hengadev/cluo_api/internal/common/auth/session"
	"github.com/hengadev/cluo_api/internal/common/contracts/identity"
	"github.com/hengadev/cluo_api/internal/domain/document"
	mw "github.com/hengadev/cluo_api/internal/common/middleware"
)

// ---------------------------------------------------------------------------
// Mocks
// ---------------------------------------------------------------------------

// mockDocumentService satisfies ports.DocumentService for handler tests.
type mockDocumentService struct {
	mock.Mock
}

func (m *mockDocumentService) CreateDocument(ctx context.Context, req *document.CreateDocumentRequest) (document.Documentable, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(document.Documentable), args.Error(1)
}
func (m *mockDocumentService) UpdateDocument(ctx context.Context, id string, req *document.UpdateDocumentRequest) (document.Documentable, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(document.Documentable), args.Error(1)
}
func (m *mockDocumentService) GetDocument(ctx context.Context, id string, dt document.DocumentType) (document.Documentable, error) {
	args := m.Called(ctx, id, dt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(document.Documentable), args.Error(1)
}
func (m *mockDocumentService) DeleteDocument(ctx context.Context, id string, dt document.DocumentType) error {
	args := m.Called(ctx, id, dt)
	return args.Error(0)
}
func (m *mockDocumentService) ListDocuments(ctx context.Context, f document.DocumentFilter, p document.Pagination) ([]document.DocumentSummary, int, error) {
	args := m.Called(ctx, f, p)
	return args.Get(0).([]document.DocumentSummary), args.Int(1), args.Error(2)
}
func (m *mockDocumentService) SendDocument(ctx context.Context, id string, dt document.DocumentType, req *document.SendDocumentRequest) error {
	args := m.Called(ctx, id, dt, req)
	return args.Error(0)
}
func (m *mockDocumentService) SignDocument(ctx context.Context, id string, dt document.DocumentType, req *document.SignDocumentRequest) error {
	args := m.Called(ctx, id, dt, req)
	return args.Error(0)
}
func (m *mockDocumentService) ArchiveDocument(ctx context.Context, id string, dt document.DocumentType) error {
	args := m.Called(ctx, id, dt)
	return args.Error(0)
}
func (m *mockDocumentService) GetDocumentHistory(ctx context.Context, id string, dt document.DocumentType, p document.Pagination) ([]*document.DocumentVersion, int, error) {
	args := m.Called(ctx, id, dt, p)
	return args.Get(0).([]*document.DocumentVersion), args.Int(1), args.Error(2)
}
func (m *mockDocumentService) CreateEstimate(ctx context.Context, e *document.Estimate) (*document.Estimate, error) {
	args := m.Called(ctx, e)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*document.Estimate), args.Error(1)
}
func (m *mockDocumentService) AcceptEstimate(ctx context.Context, id string, acceptedBy string) (*document.Mandate, error) {
	args := m.Called(ctx, id, acceptedBy)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*document.Mandate), args.Error(1)
}
func (m *mockDocumentService) UpdateEstimate(ctx context.Context, id string, items []document.EstimateItem) (*document.Estimate, error) {
	args := m.Called(ctx, id, items)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*document.Estimate), args.Error(1)
}
func (m *mockDocumentService) CreateMandate(ctx context.Context, man *document.Mandate) (*document.Mandate, error) {
	args := m.Called(ctx, man)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*document.Mandate), args.Error(1)
}
func (m *mockDocumentService) SignMandate(ctx context.Context, id string, req *document.SignDocumentRequest) (*document.Mandate, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*document.Mandate), args.Error(1)
}
func (m *mockDocumentService) ActivateMandate(ctx context.Context, id string) (*document.Mandate, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*document.Mandate), args.Error(1)
}
func (m *mockDocumentService) CreateContractFromMandate(ctx context.Context, mandateID string, c *document.Contract) (*document.Contract, error) {
	args := m.Called(ctx, mandateID, c)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*document.Contract), args.Error(1)
}
func (m *mockDocumentService) CreateContract(ctx context.Context, c *document.Contract) (*document.Contract, error) {
	args := m.Called(ctx, c)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*document.Contract), args.Error(1)
}
func (m *mockDocumentService) SignContract(ctx context.Context, id string, req *document.SignDocumentRequest) (*document.Contract, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*document.Contract), args.Error(1)
}
func (m *mockDocumentService) ActivateContract(ctx context.Context, id string) (*document.Contract, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*document.Contract), args.Error(1)
}
func (m *mockDocumentService) CreateInvoiceFromContract(ctx context.Context, contractID string) (*document.Invoice, error) {
	args := m.Called(ctx, contractID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*document.Invoice), args.Error(1)
}
func (m *mockDocumentService) CreateInvoice(ctx context.Context, inv *document.Invoice) (*document.Invoice, error) {
	args := m.Called(ctx, inv)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*document.Invoice), args.Error(1)
}
func (m *mockDocumentService) ProcessPayment(ctx context.Context, id string, req *document.PaymentRequest) (*document.Invoice, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*document.Invoice), args.Error(1)
}
func (m *mockDocumentService) VoidInvoice(ctx context.Context, id string) (*document.Invoice, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*document.Invoice), args.Error(1)
}
func (m *mockDocumentService) GetOverdueInvoices(ctx context.Context, p document.Pagination) ([]*document.Invoice, int, error) {
	args := m.Called(ctx, p)
	return args.Get(0).([]*document.Invoice), args.Int(1), args.Error(2)
}
func (m *mockDocumentService) GetDocumentWorkflow(ctx context.Context, caseID string) (*document.DocumentWorkflowResponse, error) {
	args := m.Called(ctx, caseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*document.DocumentWorkflowResponse), args.Error(1)
}
func (m *mockDocumentService) ValidateDocumentTransitions(ctx context.Context, id string, dt document.DocumentType, s document.DocumentStatus) error {
	args := m.Called(ctx, id, dt, s)
	return args.Error(0)
}

// mockAuthMiddleware is a pass-through auth middleware for handler tests.
type mockAuthMiddleware struct{}

func (m *mockAuthMiddleware) RequireAccessToken(next mw.Handler) mw.Handler { return next }
func (m *mockAuthMiddleware) RequireRefreshToken(next mw.Handler) mw.Handler { return next }
func (m *mockAuthMiddleware) RequireMinimumRole(_ identity.Role) func(mw.Handler) mw.Handler {
	return func(next mw.Handler) mw.Handler { return next }
}
func (m *mockAuthMiddleware) RequireAnyRole(_ ...identity.Role) func(mw.Handler) mw.Handler {
	return func(next mw.Handler) mw.Handler { return next }
}
func (m *mockAuthMiddleware) RequireAdmin(next mw.Handler) mw.Handler { return next }
func (m *mockAuthMiddleware) RequireServiceAuth(next mw.Handler) mw.Handler { return next }

// ---------------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------------

func newTestHandler(svc *mockDocumentService) *handler {
	return &handler{service: svc, authmw: &mockAuthMiddleware{}}
}

func authenticatedContext() context.Context {
	return context.WithValue(context.Background(), session.GetSessionContextKey(), &session.SessionInfo{
		ID:     uuid.New(),
		UserID: uuid.New(),
		Role:   identity.Administrator,
		State:  session.SessionActive,
	})
}

func mustDecodeResponse(t *testing.T, body string) map[string]interface{} {
	t.Helper()
	var result map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(body), &result))
	return result
}

// rfc3339 returns a time formatted as RFC3339
func rfc3339(t time.Time) string {
	return t.Format(time.RFC3339)
}

// ---------------------------------------------------------------------------
// Tests: Estimate
// ---------------------------------------------------------------------------

func TestCreateEstimate_Unauthorized(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	req := httptest.NewRequest("POST", "/estimates", strings.NewReader("{}"))
	w := httptest.NewRecorder()

	h.CreateEstimate(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	resp := mustDecodeResponse(t, w.Body.String())
	assert.False(t, resp["success"].(bool))
	assert.Contains(t, resp["error"], "Authentication required")
}

func TestCreateEstimate_MalformedDate(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	body := `{
		"case_id": "` + uuid.New().String() + `",
		"client_id": "` + uuid.New().String() + `",
		"estimate_number": "EST-2025-001",
		"issue_date": "not-a-date",
		"line_items": [{"description":"test","quantity":1,"unit_price":10,"subtotal":10}],
		"estimated_total": 10
	}`

	req := httptest.NewRequest("POST", "/estimates", strings.NewReader(body)).WithContext(authenticatedContext())
	w := httptest.NewRecorder()

	h.CreateEstimate(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp := mustDecodeResponse(t, w.Body.String())
	assert.Contains(t, resp["error"], "issue_date")
}

func TestCreateEstimate_ValidationError(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	// Missing estimate_number — will fail validation
	now := time.Now()
	body := `{
		"case_id": "` + uuid.New().String() + `",
		"client_id": "` + uuid.New().String() + `",
		"estimate_number": "",
		"issue_date": "` + rfc3339(now) + `",
		"line_items": [],
		"estimated_total": 0
	}`

	req := httptest.NewRequest("POST", "/estimates", strings.NewReader(body)).WithContext(authenticatedContext())
	w := httptest.NewRecorder()

	h.CreateEstimate(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	resp := mustDecodeResponse(t, w.Body.String())
	assert.Contains(t, resp["error"], "estimate number")
}

func TestCreateEstimate_Success(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	caseID := uuid.New()
	clientID := uuid.New()
	now := time.Now()

	svc.On("CreateEstimate", mock.Anything, mock.AnythingOfType("*document.Estimate")).
		Return(&document.Estimate{
			DocumentBase:   document.NewDocumentBase(caseID, clientID),
			EstimateNumber: "EST-2025-001",
			IssueDate:      now,
			LineItems: []document.EstimateItem{
				{Description: "Investigation", Quantity: 1, UnitPrice: 100, Subtotal: 100},
			},
			EstimatedTotal: 100,
		}, nil)

	body := `{
		"case_id": "` + caseID.String() + `",
		"client_id": "` + clientID.String() + `",
		"estimate_number": "EST-2025-001",
		"issue_date": "` + rfc3339(now) + `",
		"line_items": [{"description":"Investigation","quantity":1,"unit_price":100,"subtotal":100}],
		"estimated_total": 100
	}`

	req := httptest.NewRequest("POST", "/estimates", strings.NewReader(body)).WithContext(authenticatedContext())
	w := httptest.NewRecorder()

	h.CreateEstimate(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := mustDecodeResponse(t, w.Body.String())
	assert.True(t, resp["success"].(bool))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "EST-2025-001", data["estimate_number"])
	svc.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// Tests: Mandate
// ---------------------------------------------------------------------------

func TestCreateMandate_Unauthorized(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	req := httptest.NewRequest("POST", "/mandates", strings.NewReader("{}"))
	w := httptest.NewRecorder()

	h.CreateMandate(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreateMandate_MalformedDate(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	body := `{
		"case_id": "` + uuid.New().String() + `",
		"client_id": "` + uuid.New().String() + `",
		"mandate_number": "MND-2025-001",
		"issue_date": "bad-date",
		"scope_of_work": "Investigation and surveillance services",
		"valid_from": "` + rfc3339(time.Now()) + `",
		"terms_conditions": "Standard terms apply"
	}`

	req := httptest.NewRequest("POST", "/mandates", strings.NewReader(body)).WithContext(authenticatedContext())
	w := httptest.NewRecorder()

	h.CreateMandate(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp := mustDecodeResponse(t, w.Body.String())
	assert.Contains(t, resp["error"], "issue_date")
}

func TestCreateMandate_ValidationError(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	// scope_of_work too short (<20 chars)
	body := `{
		"case_id": "` + uuid.New().String() + `",
		"client_id": "` + uuid.New().String() + `",
		"mandate_number": "MND-2025-001",
		"issue_date": "` + rfc3339(time.Now()) + `",
		"scope_of_work": "too short",
		"valid_from": "` + rfc3339(time.Now()) + `",
		"terms_conditions": "Standard terms apply"
	}`

	req := httptest.NewRequest("POST", "/mandates", strings.NewReader(body)).WithContext(authenticatedContext())
	w := httptest.NewRecorder()

	h.CreateMandate(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	resp := mustDecodeResponse(t, w.Body.String())
	assert.Contains(t, resp["error"], "20 characters")
}

func TestCreateMandate_Success(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	caseID := uuid.New()
	clientID := uuid.New()
	now := time.Now()

	svc.On("CreateMandate", mock.Anything, mock.AnythingOfType("*document.Mandate")).
		Return(&document.Mandate{
			DocumentBase:    document.NewDocumentBase(caseID, clientID),
			MandateNumber:   "MND-2025-001",
			IssueDate:       now,
			ScopeOfWork:     "Investigation and surveillance services downtown",
			ValidFrom:       now,
			TermsConditions: "Standard terms apply",
		}, nil)

	body := `{
		"case_id": "` + caseID.String() + `",
		"client_id": "` + clientID.String() + `",
		"mandate_number": "MND-2025-001",
		"issue_date": "` + rfc3339(now) + `",
		"scope_of_work": "Investigation and surveillance services downtown",
		"valid_from": "` + rfc3339(now) + `",
		"terms_conditions": "Standard terms apply"
	}`

	req := httptest.NewRequest("POST", "/mandates", strings.NewReader(body)).WithContext(authenticatedContext())
	w := httptest.NewRecorder()

	h.CreateMandate(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := mustDecodeResponse(t, w.Body.String())
	assert.True(t, resp["success"].(bool))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "MND-2025-001", data["mandate_number"])
	svc.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// Tests: Contract
// ---------------------------------------------------------------------------

func TestCreateContract_Unauthorized(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	req := httptest.NewRequest("POST", "/contracts", strings.NewReader("{}"))
	w := httptest.NewRecorder()

	h.CreateContract(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreateContract_MalformedDate(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	body := `{
		"case_id": "` + uuid.New().String() + `",
		"client_id": "` + uuid.New().String() + `",
		"contract_number": "CNT-2025-001",
		"start_date": "not-valid",
		"scope_of_services": "` + strings.Repeat("x", 50) + `",
		"payment_terms": "Net 30",
		"confidentiality": "Standard NDA",
		"termination_clause": "30 days notice"
	}`

	req := httptest.NewRequest("POST", "/contracts", strings.NewReader(body)).WithContext(authenticatedContext())
	w := httptest.NewRecorder()

	h.CreateContract(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp := mustDecodeResponse(t, w.Body.String())
	assert.Contains(t, resp["error"], "start_date")
}

func TestCreateContract_ValidationError(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	// scope_of_services too short (<50 chars)
	body := `{
		"case_id": "` + uuid.New().String() + `",
		"client_id": "` + uuid.New().String() + `",
		"contract_number": "CNT-2025-001",
		"start_date": "` + rfc3339(time.Now()) + `",
		"scope_of_services": "too short",
		"payment_terms": "Net 30",
		"confidentiality": "Standard NDA",
		"termination_clause": "30 days notice"
	}`

	req := httptest.NewRequest("POST", "/contracts", strings.NewReader(body)).WithContext(authenticatedContext())
	w := httptest.NewRecorder()

	h.CreateContract(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	resp := mustDecodeResponse(t, w.Body.String())
	assert.Contains(t, resp["error"], "50 characters")
}

func TestCreateContract_Success(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	caseID := uuid.New()
	clientID := uuid.New()
	now := time.Now()
	scope := strings.Repeat("Comprehensive background investigation services. ", 2) // >50 chars

	svc.On("CreateContract", mock.Anything, mock.AnythingOfType("*document.Contract")).
		Return(&document.Contract{
			DocumentBase:      document.NewDocumentBase(caseID, clientID),
			ContractNumber:    "CNT-2025-001",
			StartDate:         now,
			ScopeOfServices:   scope,
			PaymentTerms:      "Net 30",
			Confidentiality:   "Standard NDA",
			TerminationClause: "30 days notice",
		}, nil)

	body := fmt.Sprintf(`{
		"case_id": "%s",
		"client_id": "%s",
		"contract_number": "CNT-2025-001",
		"start_date": "%s",
		"scope_of_services": "%s",
		"payment_terms": "Net 30",
		"confidentiality": "Standard NDA",
		"termination_clause": "30 days notice"
	}`, caseID, clientID, rfc3339(now), scope)

	req := httptest.NewRequest("POST", "/contracts", strings.NewReader(body)).WithContext(authenticatedContext())
	w := httptest.NewRecorder()

	h.CreateContract(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := mustDecodeResponse(t, w.Body.String())
	assert.True(t, resp["success"].(bool))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "CNT-2025-001", data["contract_number"])
	svc.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// Tests: Invoice
// ---------------------------------------------------------------------------

func TestCreateInvoice_Unauthorized(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	req := httptest.NewRequest("POST", "/invoices", strings.NewReader("{}"))
	w := httptest.NewRecorder()

	h.CreateInvoice(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreateInvoice_MalformedDate(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	body := `{
		"case_id": "` + uuid.New().String() + `",
		"client_id": "` + uuid.New().String() + `",
		"invoice_number": "INV-2025-001",
		"issue_date": "garbage",
		"due_date": "` + rfc3339(time.Now().Add(24*time.Hour)) + `",
		"line_items": [{"description":"Service","quantity":1,"unit_price":100,"subtotal":100}],
		"total_amount": 108,
		"tax_rate": 0.08,
		"tax_amount": 8
	}`

	req := httptest.NewRequest("POST", "/invoices", strings.NewReader(body)).WithContext(authenticatedContext())
	w := httptest.NewRecorder()

	h.CreateInvoice(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp := mustDecodeResponse(t, w.Body.String())
	assert.Contains(t, resp["error"], "issue_date")
}

func TestCreateInvoice_ValidationError(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	now := time.Now()
	// Due date before issue date — will fail validation
	body := `{
		"case_id": "` + uuid.New().String() + `",
		"client_id": "` + uuid.New().String() + `",
		"invoice_number": "INV-2025-001",
		"issue_date": "` + rfc3339(now) + `",
		"due_date": "` + rfc3339(now.Add(-24*time.Hour)) + `",
		"line_items": [{"description":"Service","quantity":1,"unit_price":100,"subtotal":100}],
		"total_amount": 108,
		"tax_rate": 0.08,
		"tax_amount": 8
	}`

	req := httptest.NewRequest("POST", "/invoices", strings.NewReader(body)).WithContext(authenticatedContext())
	w := httptest.NewRecorder()

	h.CreateInvoice(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	resp := mustDecodeResponse(t, w.Body.String())
	assert.Contains(t, resp["error"], "due date")
}

func TestCreateInvoice_Success(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	caseID := uuid.New()
	clientID := uuid.New()
	now := time.Now()
	dueDate := now.Add(30 * 24 * time.Hour)

	svc.On("CreateInvoice", mock.Anything, mock.AnythingOfType("*document.Invoice")).
		Return(&document.Invoice{
			DocumentBase:  document.NewDocumentBase(caseID, clientID),
			InvoiceNumber: "INV-2025-001",
			IssueDate:     now,
			DueDate:       dueDate,
			LineItems: []document.InvoiceItem{
				{Description: "Service", Quantity: 1, UnitPrice: 100, Subtotal: 100},
			},
			TotalAmount:   108,
			TaxRate:       0.08,
			TaxAmount:     8,
			PaymentStatus: document.PaymentStatusUnpaid,
		}, nil)

	body := fmt.Sprintf(`{
		"case_id": "%s",
		"client_id": "%s",
		"invoice_number": "INV-2025-001",
		"issue_date": "%s",
		"due_date": "%s",
		"line_items": [{"description":"Service","quantity":1,"unit_price":100,"subtotal":100}],
		"total_amount": 108,
		"tax_rate": 0.08,
		"tax_amount": 8
	}`, caseID, clientID, rfc3339(now), rfc3339(dueDate))

	req := httptest.NewRequest("POST", "/invoices", strings.NewReader(body)).WithContext(authenticatedContext())
	w := httptest.NewRecorder()

	h.CreateInvoice(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	resp := mustDecodeResponse(t, w.Body.String())
	assert.True(t, resp["success"].(bool))
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "INV-2025-001", data["invoice_number"])
	svc.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// Tests: parseTimeField helper
// ---------------------------------------------------------------------------

func TestParseTimeField_Valid(t *testing.T) {
	raw := json.RawMessage(`"2025-06-15T10:30:00Z"`)
	parsed, err := parseTimeField(raw, "test_field")
	require.NoError(t, err)
	assert.Equal(t, "2025-06-15T10:30:00Z", parsed.Format(time.RFC3339))
}

func TestParseTimeField_InvalidString(t *testing.T) {
	raw := json.RawMessage(`"not-a-date"`)
	_, err := parseTimeField(raw, "test_field")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "test_field")
}

func TestParseTimeField_NotString(t *testing.T) {
	raw := json.RawMessage(`12345`)
	_, err := parseTimeField(raw, "test_field")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "test_field")
}

func TestParseTimeField_Empty(t *testing.T) {
	raw := json.RawMessage{}
	parsed, err := parseTimeField(raw, "test_field")
	assert.NoError(t, err)
	assert.True(t, parsed.IsZero())
}

// ---------------------------------------------------------------------------
// Tests: getUserID helper
// ---------------------------------------------------------------------------

func TestGetUserID_WithSession(t *testing.T) {
	expectedUserID := uuid.New()
	ctx := context.WithValue(context.Background(), session.GetSessionContextKey(), &session.SessionInfo{
		ID:     uuid.New(),
		UserID: expectedUserID,
		Role:   identity.Administrator,
		State:  session.SessionActive,
	})

	h := &handler{}
	req := httptest.NewRequest("GET", "/", nil).WithContext(ctx)

	userID, err := h.getUserID(req)
	assert.NoError(t, err)
	assert.Equal(t, expectedUserID, userID)
}

func TestGetUserID_WithoutSession(t *testing.T) {
	h := &handler{}
	req := httptest.NewRequest("GET", "/", nil)

	userID, err := h.getUserID(req)
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, userID)
}
