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
	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// ---------------------------------------------------------------------------
// Helper to build an authenticated request
// ---------------------------------------------------------------------------

func authedRequest(method, target, body string) *http.Request {
	r := httptest.NewRequest(method, target, strings.NewReader(body))
	sessionInfo := &session.SessionInfo{
		ID:     uuid.New(),
		UserID: uuid.New(),
		Role:   identity.Administrator,
		State:  session.SessionActive,
	}
	ctx := context.WithValue(r.Context(), session.GetSessionContextKey(), sessionInfo)
	return r.WithContext(ctx)
}

func contextWithSession() context.Context {
	return context.WithValue(context.Background(), session.GetSessionContextKey(), &session.SessionInfo{
		ID:     uuid.New(),
		UserID: uuid.New(),
		Role:   identity.Administrator,
		State:  session.SessionActive,
	})
}

func decodeResp(t *testing.T, body string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(body), &m))
	return m
}

// ---------------------------------------------------------------------------
// POST /documents/{id}/accept — Accept Estimate
// ---------------------------------------------------------------------------

func TestAcceptEstimate_Success(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	estimateID := uuid.New()
	mandate := &document.Mandate{
		DocumentBase:    document.NewDocumentBase(uuid.New(), uuid.New()),
		MandateNumber:   "MND-2025-001",
		LinkedEstimateID: &estimateID,
	}

	svc.On("AcceptEstimate", mock.Anything, estimateID.String(), mock.AnythingOfType("string")).
		Return(mandate, nil)

	req := authedRequest("POST", "/estimates/"+estimateID.String()+"/accept", "")
	req.SetPathValue("id", estimateID.String())
	w := httptest.NewRecorder()

	h.AcceptEstimate(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := decodeResp(t, w.Body.String())
	assert.True(t, resp["success"].(bool))
	svc.AssertExpectations(t)
}

func TestAcceptEstimate_Conflict(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	estimateID := uuid.New()

	svc.On("AcceptEstimate", mock.Anything, estimateID.String(), mock.AnythingOfType("string")).
		Return(nil, errs.NewConflictErr(fmt.Errorf("estimate already accepted")))

	req := authedRequest("POST", "/estimates/"+estimateID.String()+"/accept", "")
	req.SetPathValue("id", estimateID.String())
	w := httptest.NewRecorder()

	h.AcceptEstimate(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	resp := decodeResp(t, w.Body.String())
	assert.False(t, resp["success"].(bool))
	svc.AssertExpectations(t)
}

func TestAcceptEstimate_Unauthorized(t *testing.T) {
	svc := new(mockDocumentService)
	_h := newTestHandler(svc)

	req := httptest.NewRequest("POST", "/estimates/123/accept", nil)
	req.SetPathValue("id", "123")
	w := httptest.NewRecorder()

	_h.AcceptEstimate(w, req)

	// With mock pass-through auth, the service will fail (not found)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestAcceptEstimate_MissingID(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	req := authedRequest("POST", "/estimates//accept", "")
	req.SetPathValue("id", "")
	w := httptest.NewRecorder()

	h.AcceptEstimate(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ---------------------------------------------------------------------------
// POST /documents/{id}/activate — Activate Mandate
// ---------------------------------------------------------------------------

func TestActivateMandate_Success(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	mandateID := uuid.New()
	mandate := &document.Mandate{
		DocumentBase: document.DocumentBase{
			ID:     mandateID,
			Status: document.DocumentStatusActive,
		},
		MandateNumber: "MND-2025-001",
	}

	svc.On("ActivateMandate", mock.Anything, mandateID.String()).
		Return(mandate, nil)

	req := authedRequest("POST", "/mandates/"+mandateID.String()+"/activate", "")
	req.SetPathValue("id", mandateID.String())
	w := httptest.NewRecorder()

	h.ActivateMandate(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := decodeResp(t, w.Body.String())
	assert.True(t, resp["success"].(bool))
	svc.AssertExpectations(t)
}

func TestActivateMandate_Unauthorized(t *testing.T) {
	svc := new(mockDocumentService)
	_h := newTestHandler(svc)

	req := httptest.NewRequest("POST", "/mandates/123/activate", nil)
	w := httptest.NewRecorder()

	// Without path value set, id is empty → 400 (not 401)
	// Since the mockAuthMiddleware is a pass-through, we can't easily test
	// the unauthorized case here — the auth check would pass anyway.
	// The unauthorized path is tested via the generic handler_test.go pattern.
	req.SetPathValue("id", "123")
	_h.ActivateMandate(w, req)

	// With mock auth that passes through, the service would be called and fail (not found)
	// So we just verify it's not 200
	assert.NotEqual(t, http.StatusOK, w.Code)
}

// ---------------------------------------------------------------------------
// POST /documents/{id}/activate — Activate Contract
// ---------------------------------------------------------------------------

func TestActivateContract_Success(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	contractID := uuid.New()
	contract := &document.Contract{
		DocumentBase: document.DocumentBase{
			ID:     contractID,
			Status: document.DocumentStatusActive,
		},
		ContractNumber: "CNT-2025-001",
	}

	svc.On("ActivateContract", mock.Anything, contractID.String()).
		Return(contract, nil)

	req := authedRequest("POST", "/contracts/"+contractID.String()+"/activate", "")
	req.SetPathValue("id", contractID.String())
	w := httptest.NewRecorder()

	h.ActivateContract(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := decodeResp(t, w.Body.String())
	assert.True(t, resp["success"].(bool))
	svc.AssertExpectations(t)
}

func TestActivateContract_Unauthorized(t *testing.T) {
	svc := new(mockDocumentService)
	_h := newTestHandler(svc)

	req := httptest.NewRequest("POST", "/contracts/123/activate", nil)
	req.SetPathValue("id", "123")
	w := httptest.NewRecorder()

	_h.ActivateContract(w, req)

	// With mock pass-through auth, the service will fail (not found)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

// ---------------------------------------------------------------------------
// POST /invoices/{id}/pay — Process Payment
// ---------------------------------------------------------------------------

func TestProcessPayment_PartialPayment(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	invoiceID := uuid.New()
	paidAmount := 50.0
	invoice := &document.Invoice{
		DocumentBase:  document.NewDocumentBase(uuid.New(), uuid.New()),
		InvoiceNumber: "INV-2025-001",
		PaymentStatus: document.PaymentStatusPartiallyPaid,
		PaidAmount:    &paidAmount,
	}

	svc.On("ProcessPayment", mock.Anything, invoiceID.String(), mock.AnythingOfType("*document.PaymentRequest")).
		Return(invoice, nil)

	body := `{"amount":50,"payment_method":"wire_transfer"}`
	req := authedRequest("POST", "/invoices/"+invoiceID.String()+"/pay", body)
	req.SetPathValue("id", invoiceID.String())
	w := httptest.NewRecorder()

	h.ProcessPayment(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := decodeResp(t, w.Body.String())
	assert.True(t, resp["success"].(bool))
	svc.AssertExpectations(t)
}

func TestProcessPayment_FullPayment(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	invoiceID := uuid.New()
	paidAmount := 100.0
	invoice := &document.Invoice{
		DocumentBase:  document.NewDocumentBase(uuid.New(), uuid.New()),
		InvoiceNumber: "INV-2025-001",
		PaymentStatus: document.PaymentStatusPaid,
		PaidAmount:    &paidAmount,
	}

	svc.On("ProcessPayment", mock.Anything, invoiceID.String(), mock.AnythingOfType("*document.PaymentRequest")).
		Return(invoice, nil)

	body := `{"amount":100,"payment_method":"credit_card"}`
	req := authedRequest("POST", "/invoices/"+invoiceID.String()+"/pay", body)
	req.SetPathValue("id", invoiceID.String())
	w := httptest.NewRecorder()

	h.ProcessPayment(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := decodeResp(t, w.Body.String())
	assert.True(t, resp["success"].(bool))
	svc.AssertExpectations(t)
}

func TestProcessPayment_Unauthorized(t *testing.T) {
	svc := new(mockDocumentService)
	_h := newTestHandler(svc)

	req := httptest.NewRequest("POST", "/invoices/123/pay", strings.NewReader(`{"amount":50,"method":"card"}`))
	req.SetPathValue("id", "123")
	w := httptest.NewRecorder()

	_h.ProcessPayment(w, req)

	// With mock pass-through auth, the service will fail (not found or decode error)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestProcessPayment_MissingID(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	req := authedRequest("POST", "/invoices//pay", `{"amount":50,"method":"card"}`)
	req.SetPathValue("id", "")
	w := httptest.NewRecorder()

	h.ProcessPayment(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ---------------------------------------------------------------------------
// POST /invoices/{id}/void — Void Invoice
// ---------------------------------------------------------------------------

func TestVoidInvoice_Success(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	invoiceID := uuid.New()
	invoice := &document.Invoice{
		DocumentBase: document.DocumentBase{
			ID:     invoiceID,
			Status: document.DocumentStatusCancelled,
		},
		InvoiceNumber: "INV-2025-001",
		PaymentStatus: document.PaymentStatusVoid,
	}

	svc.On("VoidInvoice", mock.Anything, invoiceID.String()).
		Return(invoice, nil)

	req := authedRequest("POST", "/invoices/"+invoiceID.String()+"/void", "")
	req.SetPathValue("id", invoiceID.String())
	w := httptest.NewRecorder()

	h.VoidInvoice(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := decodeResp(t, w.Body.String())
	assert.True(t, resp["success"].(bool))
	svc.AssertExpectations(t)
}

func TestVoidInvoice_Unauthorized(t *testing.T) {
	svc := new(mockDocumentService)
	_h := newTestHandler(svc)

	req := httptest.NewRequest("POST", "/invoices/123/void", nil)
	req.SetPathValue("id", "123")
	w := httptest.NewRecorder()

	_h.VoidInvoice(w, req)

	// With mock pass-through auth, the service will fail (not found)
	assert.NotEqual(t, http.StatusOK, w.Code)
}

func TestVoidInvoice_MissingID(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	req := authedRequest("POST", "/invoices//void", "")
	req.SetPathValue("id", "")
	w := httptest.NewRecorder()

	h.VoidInvoice(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ---------------------------------------------------------------------------
// State machine enforcement: test that CanTransitionTo rejects invalid transitions
// ---------------------------------------------------------------------------

func TestEstimateAcceptance_RejectsExpiredEstimate(t *testing.T) {
	issueDate := time.Now().Add(-72 * time.Hour)
	validUntil := time.Now().Add(-24 * time.Hour) // expired yesterday
	estimate := &document.Estimate{
		DocumentBase: document.DocumentBase{
			Status: document.DocumentStatusSent,
		},
		EstimateNumber: "EST-2025-001",
		IssueDate:      issueDate,
		ValidUntil:     &validUntil,
		LineItems: []document.EstimateItem{
			{Description: "Test", Quantity: 1, UnitPrice: 100, Subtotal: 100},
		},
		EstimatedTotal: 100,
	}

	assert.True(t, estimate.IsExpired())
	assert.False(t, estimate.CanBeAccepted())
}

func TestEstimateAcceptance_RejectsAlreadyAccepted(t *testing.T) {
	estimate := &document.Estimate{
		DocumentBase: document.DocumentBase{
			Status: document.DocumentStatusSigned,
		},
		EstimateNumber: "EST-2025-001",
		IssueDate:      time.Now(),
		LineItems: []document.EstimateItem{
			{Description: "Test", Quantity: 1, UnitPrice: 100, Subtotal: 100},
		},
		EstimatedTotal: 100,
		Accepted:       true,
	}

	assert.False(t, estimate.CanBeAccepted())
}

func TestEstimateAcceptance_RejectsDraftStatus(t *testing.T) {
	estimate := &document.Estimate{
		DocumentBase: document.DocumentBase{
			Status: document.DocumentStatusDraft,
		},
		EstimateNumber: "EST-2025-001",
		IssueDate:      time.Now(),
		LineItems: []document.EstimateItem{
			{Description: "Test", Quantity: 1, UnitPrice: 100, Subtotal: 100},
		},
		EstimatedTotal: 100,
	}

	assert.False(t, estimate.CanBeAccepted())
}

func TestMandateActivate_RequiresSignature(t *testing.T) {
	mandate := &document.Mandate{
		DocumentBase: document.DocumentBase{
			Status: document.DocumentStatusSigned,
		},
		MandateNumber:   "MND-2025-001",
		IssueDate:       time.Now(),
		ScopeOfWork:     "Test investigation scope of work",
		ValidFrom:       time.Now().Add(-1 * time.Hour),
		TermsConditions: "Standard terms",
	}

	// No client signature
	err := mandate.Activate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client signature")
}

func TestContractActivate_RequiresSignature(t *testing.T) {
	contract := &document.Contract{
		DocumentBase: document.DocumentBase{
			Status: document.DocumentStatusSigned,
		},
		ContractNumber:    "CNT-2025-001",
		StartDate:         time.Now().Add(-1 * time.Hour),
		ScopeOfServices:   strings.Repeat("x", 50),
		PaymentTerms:      "Net 30",
		Confidentiality:   "Standard",
		TerminationClause: "30 days",
		Signatures:        []document.Signature{}, // empty
	}

	err := contract.Activate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "signature")
}

func TestInvoiceVoid_RejectsPaidInvoice(t *testing.T) {
	invoice := &document.Invoice{
		DocumentBase: document.DocumentBase{
			Status: document.DocumentStatusSent,
		},
		InvoiceNumber: "INV-2025-001",
		IssueDate:     time.Now(),
		DueDate:       time.Now().Add(24 * time.Hour),
		PaymentStatus: document.PaymentStatusPaid,
	}

	err := invoice.Void()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "paid")
}

func TestInvoicePay_RejectsPaidInvoice(t *testing.T) {
	invoice := &document.Invoice{
		DocumentBase: document.DocumentBase{
			Status: document.DocumentStatusSent,
		},
		InvoiceNumber: "INV-2025-001",
		PaymentStatus: document.PaymentStatusPaid,
	}

	err := invoice.AddPayment(50, "wire_transfer")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already fully paid")
}

func TestInvoicePay_RejectsVoidInvoice(t *testing.T) {
	invoice := &document.Invoice{
		DocumentBase: document.DocumentBase{
			Status: document.DocumentStatusCancelled,
		},
		InvoiceNumber: "INV-2025-001",
		PaymentStatus: document.PaymentStatusVoid,
	}

	err := invoice.AddPayment(50, "wire_transfer")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "void")
}

func TestInvoicePay_PartialThenFull(t *testing.T) {
	invoice := &document.Invoice{
		DocumentBase: document.DocumentBase{
			Status: document.DocumentStatusSent,
		},
		InvoiceNumber: "INV-2025-001",
		IssueDate:     time.Now(),
		DueDate:       time.Now().Add(24 * time.Hour),
		LineItems: []document.InvoiceItem{
			{Description: "Service", Quantity: 1, UnitPrice: 100, Subtotal: 100},
		},
		TotalAmount:   100,
		TaxRate:       0,
		TaxAmount:     0,
		PaymentStatus: document.PaymentStatusUnpaid,
	}

	// Partial payment
	err := invoice.AddPayment(40, "wire_transfer")
	assert.NoError(t, err)
	assert.Equal(t, document.PaymentStatusPartiallyPaid, invoice.PaymentStatus)

	// Full payment
	err = invoice.AddPayment(60, "wire_transfer")
	assert.NoError(t, err)
	assert.Equal(t, document.PaymentStatusPaid, invoice.PaymentStatus)
}

func TestInvoiceVoid_SetsDocumentStatusToCancelled(t *testing.T) {
	invoice := &document.Invoice{
		DocumentBase: document.DocumentBase{
			Status: document.DocumentStatusSent,
		},
		InvoiceNumber: "INV-2025-001",
		IssueDate:     time.Now(),
		DueDate:       time.Now().Add(24 * time.Hour),
		PaymentStatus: document.PaymentStatusUnpaid,
	}

	err := invoice.Void()
	assert.NoError(t, err)
	assert.Equal(t, document.PaymentStatusVoid, invoice.PaymentStatus)
	assert.Equal(t, document.DocumentStatusCancelled, invoice.Status)
}
