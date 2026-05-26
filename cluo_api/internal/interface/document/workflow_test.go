package document

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/hengadev/cluo_api/internal/domain/document"
)

// ---------------------------------------------------------------------------
// GET /cases/{caseId}/document-workflow
// ---------------------------------------------------------------------------

func TestGetDocumentWorkflow_Unauthorized(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	req := httptest.NewRequest("GET", "/cases/"+uuid.New().String()+"/document-workflow", nil)
	req.SetPathValue("caseId", uuid.New().String())
	w := httptest.NewRecorder()

	h.GetDocumentWorkflow(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	resp := mustDecodeResponse(t, w.Body.String())
	assert.False(t, resp["success"].(bool))
	assert.Contains(t, resp["error"], "Authentication required")
}

func TestGetDocumentWorkflow_MissingCaseID(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	req := authedRequest("GET", "/cases//document-workflow", "")
	req.SetPathValue("caseId", "")
	w := httptest.NewRecorder()

	h.GetDocumentWorkflow(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp := mustDecodeResponse(t, w.Body.String())
	assert.Contains(t, resp["error"], "Case ID is required")
}

func TestGetDocumentWorkflow_CaseNotFound(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	caseID := uuid.New()

	svc.On("GetDocumentWorkflow", mock.Anything, caseID.String()).
		Return(nil, errCaseNotFound(caseID))

	req := authedRequest("GET", "/cases/"+caseID.String()+"/document-workflow", "")
	req.SetPathValue("caseId", caseID.String())
	w := httptest.NewRecorder()

	h.GetDocumentWorkflow(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	resp := mustDecodeResponse(t, w.Body.String())
	assert.False(t, resp["success"].(bool))
	assert.Contains(t, resp["error"], "Case not found")
	svc.AssertExpectations(t)
}

func TestGetDocumentWorkflow_AllDocumentsPresent(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	caseID := uuid.New()
	clientID := uuid.New()

	estimate := &document.Estimate{
		DocumentBase:   document.DocumentBase{ID: uuid.New(), CaseID: caseID, ClientID: clientID, Status: document.DocumentStatusSigned, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		EstimateNumber: "EST-2025-001",
		IssueDate:      time.Now(),
		EstimatedTotal: 1000,
	}
	mandate := &document.Mandate{
		DocumentBase:     document.DocumentBase{ID: uuid.New(), CaseID: caseID, ClientID: clientID, Status: document.DocumentStatusActive, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		MandateNumber:    "MND-2025-001",
		LinkedEstimateID: &estimate.ID,
	}
	contract := &document.Contract{
		DocumentBase:    document.DocumentBase{ID: uuid.New(), CaseID: caseID, ClientID: clientID, Status: document.DocumentStatusActive, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		ContractNumber:  "CNT-2025-001",
		LinkedMandateID: &mandate.ID,
	}
	invoice := &document.Invoice{
		DocumentBase:     document.DocumentBase{ID: uuid.New(), CaseID: caseID, ClientID: clientID, Status: document.DocumentStatusSent, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		InvoiceNumber:    "INV-2025-001",
		PaymentStatus:    document.PaymentStatusUnpaid,
		LinkedContractID: &contract.ID,
		TotalAmount:      1000,
		DueDate:          time.Now().Add(30 * 24 * time.Hour),
	}

	workflow := &document.DocumentWorkflowResponse{
		Estimate: estimate,
		Mandate:  mandate,
		Contract: contract,
		Invoice:  invoice,
	}

	svc.On("GetDocumentWorkflow", mock.Anything, caseID.String()).
		Return(workflow, nil)

	req := authedRequest("GET", "/cases/"+caseID.String()+"/document-workflow", "")
	req.SetPathValue("caseId", caseID.String())
	w := httptest.NewRecorder()

	h.GetDocumentWorkflow(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := mustDecodeResponse(t, w.Body.String())
	assert.True(t, resp["success"].(bool))

	data := resp["data"].(map[string]interface{})
	assert.NotNil(t, data["estimate"])
	assert.NotNil(t, data["mandate"])
	assert.NotNil(t, data["contract"])
	assert.NotNil(t, data["invoice"])

	// Verify estimate details
	est := data["estimate"].(map[string]interface{})
	assert.Equal(t, "EST-2025-001", est["estimate_number"])
	assert.Equal(t, "signed", est["status"])

	// Verify invoice details
	inv := data["invoice"].(map[string]interface{})
	assert.Equal(t, "INV-2025-001", inv["invoice_number"])
	assert.Equal(t, "unpaid", inv["payment_status"])

	svc.AssertExpectations(t)
}

func TestGetDocumentWorkflow_OnlyEstimatePresent(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	caseID := uuid.New()
	clientID := uuid.New()

	estimate := &document.Estimate{
		DocumentBase:   document.DocumentBase{ID: uuid.New(), CaseID: caseID, ClientID: clientID, Status: document.DocumentStatusDraft, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		EstimateNumber: "EST-2025-002",
		IssueDate:      time.Now(),
		EstimatedTotal: 500,
	}

	workflow := &document.DocumentWorkflowResponse{
		Estimate: estimate,
		Mandate:  nil,
		Contract: nil,
		Invoice:  nil,
	}

	svc.On("GetDocumentWorkflow", mock.Anything, caseID.String()).
		Return(workflow, nil)

	req := authedRequest("GET", "/cases/"+caseID.String()+"/document-workflow", "")
	req.SetPathValue("caseId", caseID.String())
	w := httptest.NewRecorder()

	h.GetDocumentWorkflow(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := mustDecodeResponse(t, w.Body.String())
	assert.True(t, resp["success"].(bool))

	data := resp["data"].(map[string]interface{})
	assert.NotNil(t, data["estimate"])
	assert.Nil(t, data["mandate"])
	assert.Nil(t, data["contract"])
	assert.Nil(t, data["invoice"])

	svc.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// GET /invoices/overdue
// ---------------------------------------------------------------------------

func TestGetOverdueInvoices_Unauthorized(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	req := httptest.NewRequest("GET", "/invoices/overdue", nil)
	w := httptest.NewRecorder()

	h.GetOverdueInvoices(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	resp := mustDecodeResponse(t, w.Body.String())
	assert.False(t, resp["success"].(bool))
	assert.Contains(t, resp["error"], "Authentication required")
}

func TestGetOverdueInvoices_InvalidPagination(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	req := authedRequest("GET", "/invoices/overdue?page=-1", "")
	req = req.WithContext(req.Context())
	w := httptest.NewRecorder()

	h.GetOverdueInvoices(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resp := mustDecodeResponse(t, w.Body.String())
	assert.Contains(t, resp["error"], "page")
}

func TestGetOverdueInvoices_WithResults(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	invoice1 := &document.Invoice{
		DocumentBase:  document.NewDocumentBase(uuid.New(), uuid.New()),
		InvoiceNumber: "INV-2025-001",
		DueDate:       time.Now().Add(-48 * time.Hour), // 2 days overdue
		PaymentStatus: document.PaymentStatusUnpaid,
		TotalAmount:   1000,
		TaxRate:       0,
		TaxAmount:     0,
	}
	invoice2 := &document.Invoice{
		DocumentBase:  document.NewDocumentBase(uuid.New(), uuid.New()),
		InvoiceNumber: "INV-2025-002",
		DueDate:       time.Now().Add(-24 * time.Hour), // 1 day overdue
		PaymentStatus: document.PaymentStatusUnpaid,
		TotalAmount:   500,
		TaxRate:       0,
		TaxAmount:     0,
	}

	svc.On("GetOverdueInvoices", mock.Anything, mock.AnythingOfType("document.Pagination")).
		Return([]*document.Invoice{invoice1, invoice2}, 2, nil)

	req := authedRequest("GET", "/invoices/overdue?page=1&per_page=20", "")
	w := httptest.NewRecorder()

	h.GetOverdueInvoices(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := mustDecodeResponse(t, w.Body.String())
	assert.True(t, resp["success"].(bool))
	assert.Equal(t, float64(2), resp["total"].(float64))
	assert.Equal(t, float64(1), resp["page"].(float64))
	assert.Equal(t, float64(20), resp["per_page"].(float64))

	data := resp["data"].([]interface{})
	assert.Len(t, data, 2)

	inv1 := data[0].(map[string]interface{})
	assert.Equal(t, "INV-2025-001", inv1["invoice_number"])
	assert.Equal(t, "unpaid", inv1["payment_status"])

	svc.AssertExpectations(t)
}

func TestGetOverdueInvoices_EmptyList(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	svc.On("GetOverdueInvoices", mock.Anything, mock.AnythingOfType("document.Pagination")).
		Return([]*document.Invoice{}, 0, nil)

	req := authedRequest("GET", "/invoices/overdue", "")
	w := httptest.NewRecorder()

	h.GetOverdueInvoices(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := mustDecodeResponse(t, w.Body.String())
	assert.True(t, resp["success"].(bool))
	assert.Equal(t, float64(0), resp["total"].(float64))

	data := resp["data"].([]interface{})
	assert.Len(t, data, 0)

	svc.AssertExpectations(t)
}

func TestGetOverdueInvoices_PaginationRespected(t *testing.T) {
	svc := new(mockDocumentService)
	h := newTestHandler(svc)

	svc.On("GetOverdueInvoices", mock.Anything, document.Pagination{Page: 2, PageSize: 5}).
		Return([]*document.Invoice{}, 0, nil)

	req := authedRequest("GET", "/invoices/overdue?page=2&per_page=5", "")
	w := httptest.NewRecorder()

	h.GetOverdueInvoices(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resp := mustDecodeResponse(t, w.Body.String())
	assert.Equal(t, float64(2), resp["page"].(float64))
	assert.Equal(t, float64(5), resp["per_page"].(float64))

	svc.AssertExpectations(t)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// errCaseNotFound creates a simple error with "case not found" message.
func errCaseNotFound(caseID uuid.UUID) error {
	return &caseNotFoundError{caseID: caseID}
}

type caseNotFoundError struct {
	caseID uuid.UUID
}

func (e *caseNotFoundError) Error() string {
	return "case not found"
}
