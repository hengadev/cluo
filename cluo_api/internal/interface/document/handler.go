package document

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"github.com/hengadev/cluo_api/internal/common/errs"
	"github.com/hengadev/cluo_api/internal/common/middleware/auth"
	"github.com/hengadev/cluo_api/internal/domain/document"
	"github.com/hengadev/cluo_api/internal/ports"
)

type Handler interface {
	RegisterRoutes(router *http.ServeMux)
}

type handler struct {
	service ports.DocumentService
	authmw  auth.AuthMiddleware
}

func New(service ports.DocumentService, authmw auth.AuthMiddleware) Handler {
	return &handler{service: service, authmw: authmw}
}

// Response types
type DocumentResponse struct {
	Success bool        `json:"success"`
	Data    any `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type DocumentListResponse struct {
	Success bool                   `json:"success"`
	Data    []document.DocumentSummary `json:"data"`
	Total   int                    `json:"total"`
	Page    int                    `json:"page"`
	PerPage int                    `json:"per_page"`
}

// Helper functions

func (h *handler) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *handler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, DocumentResponse{
		Success: false,
		Error:   message,
	})
}

func (h *handler) writeSuccess(w http.ResponseWriter, data any) {
	h.writeJSON(w, http.StatusOK, DocumentResponse{
		Success: true,
		Data:    data,
	})
}

func (h *handler) writeCreated(w http.ResponseWriter, data any) {
	h.writeJSON(w, http.StatusCreated, DocumentResponse{
		Success: true,
		Data:    data,
	})
}

func (h *handler) getPaginationFromRequest(r *http.Request) (document.Pagination, error) {
	pageStr := r.URL.Query().Get("page")
	perPageStr := r.URL.Query().Get("per_page")

	page := 1
	perPage := 20

	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p < 1 {
			return document.Pagination{}, &ValidationError{Message: "invalid page parameter"}
		}
		page = p
	}

	if perPageStr != "" {
		pp, err := strconv.Atoi(perPageStr)
		if err != nil || pp < 1 || pp > 100 {
			return document.Pagination{}, &ValidationError{Message: "invalid per_page parameter (must be 1-100)"}
		}
		perPage = pp
	}

	return document.Pagination{
		Page:     page,
		PageSize: perPage,
	}, nil
}

func (h *handler) getDocumentFilterFromRequest(r *http.Request) document.DocumentFilter {
	filter := document.DocumentFilter{}

	// Parse query parameters
	if docType := r.URL.Query().Get("type"); docType != "" {
		dt := document.DocumentType(docType)
		filter.Type = &dt
	}

	if status := r.URL.Query().Get("status"); status != "" {
		ds := document.DocumentStatus(status)
		filter.Status = &ds
	}

	if caseID := r.URL.Query().Get("case_id"); caseID != "" {
		if uid, err := uuid.Parse(caseID); err == nil {
			filter.CaseID = &uid
		}
	}

	if clientID := r.URL.Query().Get("client_id"); clientID != "" {
		if uid, err := uuid.Parse(clientID); err == nil {
			filter.ClientID = &uid
		}
	}

	if search := r.URL.Query().Get("search"); search != "" {
		filter.Search = &search
	}

	// TODO: Parse date ranges from query parameters
	// dateFrom := r.URL.Query().Get("date_from")
	// dateTo := r.URL.Query().Get("date_to")

	return filter
}

// Generic document handlers

func (h *handler) CreateDocument(w http.ResponseWriter, r *http.Request) {
	if _, err := h.getUserID(r); err != nil {
		h.writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	var req document.CreateDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := req.Valid(r.Context()); err != nil {
		h.writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	doc, err := h.service.CreateDocument(r.Context(), &req)
	if err != nil {
		if _, ok := err.(*ValidationError); ok {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to create document")
		return
	}

	h.writeSuccess(w, doc)
}

func (h *handler) GetDocument(w http.ResponseWriter, r *http.Request) {
	documentID := r.PathValue("id")
	docType := r.PathValue("type")

	if documentID == "" || docType == "" {
		h.writeError(w, http.StatusBadRequest, "Document ID and type are required")
		return
	}

	doc, err := h.service.GetDocument(r.Context(), documentID, document.DocumentType(docType))
	if err != nil {
		if err.Error() == "document not found" {
			h.writeError(w, http.StatusNotFound, "Document not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to get document")
		return
	}

	h.writeSuccess(w, doc)
}

func (h *handler) UpdateDocument(w http.ResponseWriter, r *http.Request) {
	documentID := r.PathValue("id")
	docType := r.PathValue("type")

	if documentID == "" || docType == "" {
		h.writeError(w, http.StatusBadRequest, "Document ID and type are required")
		return
	}

	if _, err := h.getUserID(r); err != nil {
		h.writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	var req document.UpdateDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := req.Valid(r.Context()); err != nil {
		h.writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	doc, err := h.service.UpdateDocument(r.Context(), documentID, &req)
	if err != nil {
		if err.Error() == "document not found" {
			h.writeError(w, http.StatusNotFound, "Document not found")
			return
		}
		if _, ok := err.(*ValidationError); ok {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to update document")
		return
	}

	h.writeSuccess(w, doc)
}

func (h *handler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	documentID := r.PathValue("id")
	docType := r.PathValue("type")

	if documentID == "" || docType == "" {
		h.writeError(w, http.StatusBadRequest, "Document ID and type are required")
		return
	}

	err := h.service.DeleteDocument(r.Context(), documentID, document.DocumentType(docType))
	if err != nil {
		if err.Error() == "document not found" {
			h.writeError(w, http.StatusNotFound, "Document not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to delete document")
		return
	}

	h.writeSuccess(w, map[string]string{"message": "Document deleted successfully"})
}

func (h *handler) ListDocuments(w http.ResponseWriter, r *http.Request) {
	pagination, err := h.getPaginationFromRequest(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	filter := h.getDocumentFilterFromRequest(r)

	documents, total, err := h.service.ListDocuments(r.Context(), filter, pagination)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to list documents")
		return
	}

	h.writeJSON(w, http.StatusOK, DocumentListResponse{
		Success: true,
		Data:    documents,
		Total:   total,
		Page:    pagination.Page,
		PerPage: pagination.PageSize,
	})
}

func (h *handler) SendDocument(w http.ResponseWriter, r *http.Request) {
	documentID := r.PathValue("id")
	docType := r.PathValue("type")

	if documentID == "" || docType == "" {
		h.writeError(w, http.StatusBadRequest, "Document ID and type are required")
		return
	}

	if _, err := h.getUserID(r); err != nil {
		h.writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	var req document.SendDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := req.Valid(r.Context()); err != nil {
		h.writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	err := h.service.SendDocument(r.Context(), documentID, document.DocumentType(docType), &req)
	if err != nil {
		if err.Error() == "document not found" {
			h.writeError(w, http.StatusNotFound, "Document not found")
			return
		}
		if isConflictError(err) {
			h.writeError(w, http.StatusConflict, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to send document")
		return
	}

	h.writeSuccess(w, map[string]string{"message": "Document sent successfully"})
}

func (h *handler) SignDocument(w http.ResponseWriter, r *http.Request) {
	documentID := r.PathValue("id")
	docType := r.PathValue("type")

	if documentID == "" || docType == "" {
		h.writeError(w, http.StatusBadRequest, "Document ID and type are required")
		return
	}

	if _, err := h.getUserID(r); err != nil {
		h.writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	var req document.SignDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := req.Valid(r.Context()); err != nil {
		h.writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	err := h.service.SignDocument(r.Context(), documentID, document.DocumentType(docType), &req)
	if err != nil {
		if err.Error() == "document not found" {
			h.writeError(w, http.StatusNotFound, "Document not found")
			return
		}
		if isConflictError(err) {
			h.writeError(w, http.StatusConflict, err.Error())
			return
		}
		if _, ok := err.(*ValidationError); ok {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to sign document")
		return
	}

	h.writeSuccess(w, map[string]string{"message": "Document signed successfully"})
}

func (h *handler) ArchiveDocument(w http.ResponseWriter, r *http.Request) {
	documentID := r.PathValue("id")
	docType := r.PathValue("type")

	if documentID == "" || docType == "" {
		h.writeError(w, http.StatusBadRequest, "Document ID and type are required")
		return
	}

	err := h.service.ArchiveDocument(r.Context(), documentID, document.DocumentType(docType))
	if err != nil {
		if err.Error() == "document not found" {
			h.writeError(w, http.StatusNotFound, "Document not found")
			return
		}
		if isConflictError(err) {
			h.writeError(w, http.StatusConflict, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to archive document")
		return
	}

	h.writeSuccess(w, map[string]string{"message": "Document archived successfully"})
}

func (h *handler) GetDocumentHistory(w http.ResponseWriter, r *http.Request) {
	documentID := r.PathValue("id")
	docType := r.PathValue("type")

	if documentID == "" || docType == "" {
		h.writeError(w, http.StatusBadRequest, "Document ID and type are required")
		return
	}

	pagination, err := h.getPaginationFromRequest(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	versions, total, err := h.service.GetDocumentHistory(r.Context(), documentID, document.DocumentType(docType), pagination)
	if err != nil {
		if err.Error() == "document not found" {
			h.writeError(w, http.StatusNotFound, "Document not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to get document history")
		return
	}

	response := struct {
		Success   bool                      `json:"success"`
		Data      []*document.DocumentVersion  `json:"data"`
		Total     int                       `json:"total"`
		Page      int                       `json:"page"`
		PerPage   int                       `json:"per_page"`
	}{
		Success: true,
		Data:    versions,
		Total:   total,
		Page:    pagination.Page,
		PerPage: pagination.PageSize,
	}

	h.writeJSON(w, http.StatusOK, response)
}

func (h *handler) GetDocumentWorkflow(w http.ResponseWriter, r *http.Request) {
	caseID := r.PathValue("caseId")
	if caseID == "" {
		h.writeError(w, http.StatusBadRequest, "Case ID is required")
		return
	}

	if _, err := h.getUserID(r); err != nil {
		h.writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	documents, err := h.service.GetDocumentWorkflow(r.Context(), caseID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get document workflow")
		return
	}

	h.writeSuccess(w, documents)
}

// Custom error type for validation errors
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// isConflictError returns true if the error wraps errs.ErrConflict.
func isConflictError(err error) bool {
	return errors.Is(err, errs.ErrConflict)
}

// ---------------------------------------------------------------------------
// Generic lifecycle action handlers (type inferred from stored document)
// ---------------------------------------------------------------------------

// acceptRequest is the body for POST /documents/{id}/accept (empty for now).
type acceptRequest struct{}

func (h *handler) AcceptDocument(w http.ResponseWriter, r *http.Request) {
	documentID := r.PathValue("id")
	if documentID == "" {
		h.writeError(w, http.StatusBadRequest, "Document ID is required")
		return
	}

	userID, err := h.getUserID(r)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	// Resolve document to determine its type
	result, err := h.resolveDocument(r, documentID)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "Document not found")
		return
	}

	switch result.docType {
	case document.DocumentTypeEstimate:
		mandate, err := h.service.AcceptEstimate(r.Context(), documentID, userID.String())
		if err != nil {
			if isConflictError(err) {
				h.writeError(w, http.StatusConflict, err.Error())
				return
			}
			h.writeError(w, http.StatusInternalServerError, "Failed to accept estimate")
			return
		}
		h.writeSuccess(w, mandate)
	default:
		h.writeError(w, http.StatusBadRequest, fmt.Sprintf("accept is not supported for document type %s", result.docType))
	}
}

func (h *handler) ActivateDocument(w http.ResponseWriter, r *http.Request) {
	documentID := r.PathValue("id")
	if documentID == "" {
		h.writeError(w, http.StatusBadRequest, "Document ID is required")
		return
	}

	if _, err := h.getUserID(r); err != nil {
		h.writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	result, err := h.resolveDocument(r, documentID)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "Document not found")
		return
	}

	switch result.docType {
	case document.DocumentTypeMandate:
		mandate, err := h.service.ActivateMandate(r.Context(), documentID)
		if err != nil {
			if isConflictError(err) {
				h.writeError(w, http.StatusConflict, err.Error())
				return
			}
			h.writeError(w, http.StatusInternalServerError, "Failed to activate mandate")
			return
		}
		h.writeSuccess(w, mandate)
	case document.DocumentTypeContract:
		contract, err := h.service.ActivateContract(r.Context(), documentID)
		if err != nil {
			if isConflictError(err) {
				h.writeError(w, http.StatusConflict, err.Error())
				return
			}
			h.writeError(w, http.StatusInternalServerError, "Failed to activate contract")
			return
		}
		h.writeSuccess(w, contract)
	default:
		h.writeError(w, http.StatusBadRequest, fmt.Sprintf("activate is not supported for document type %s", result.docType))
	}
}

type payRequest struct {
	Amount float64 `json:"amount"`
	Method string  `json:"method"`
}

func (h *handler) PayDocument(w http.ResponseWriter, r *http.Request) {
	documentID := r.PathValue("id")
	if documentID == "" {
		h.writeError(w, http.StatusBadRequest, "Document ID is required")
		return
	}

	if _, err := h.getUserID(r); err != nil {
		h.writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	var req payRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if req.Amount <= 0 {
		h.writeError(w, http.StatusBadRequest, "amount must be greater than 0")
		return
	}
	if req.Method == "" {
		h.writeError(w, http.StatusBadRequest, "method is required")
		return
	}

	result, err := h.resolveDocument(r, documentID)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "Document not found")
		return
	}

	if result.docType != document.DocumentTypeInvoice {
		h.writeError(w, http.StatusBadRequest, fmt.Sprintf("pay is not supported for document type %s", result.docType))
		return
	}

	paymentReq := &document.PaymentRequest{
		Amount:        req.Amount,
		PaymentMethod: req.Method,
	}

	invoice, err := h.service.ProcessPayment(r.Context(), documentID, paymentReq)
	if err != nil {
		if isConflictError(err) {
			h.writeError(w, http.StatusConflict, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to process payment")
		return
	}

	h.writeSuccess(w, invoice)
}

func (h *handler) VoidDocument(w http.ResponseWriter, r *http.Request) {
	documentID := r.PathValue("id")
	if documentID == "" {
		h.writeError(w, http.StatusBadRequest, "Document ID is required")
		return
	}

	if _, err := h.getUserID(r); err != nil {
		h.writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	result, err := h.resolveDocument(r, documentID)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "Document not found")
		return
	}

	if result.docType != document.DocumentTypeInvoice {
		h.writeError(w, http.StatusBadRequest, fmt.Sprintf("void is not supported for document type %s", result.docType))
		return
	}

	invoice, err := h.service.VoidInvoice(r.Context(), documentID)
	if err != nil {
		if isConflictError(err) {
			h.writeError(w, http.StatusConflict, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to void invoice")
		return
	}

	h.writeSuccess(w, invoice)
}

// resolvedDoc holds a document lookup result with its type.
type resolvedDoc struct {
	docType document.DocumentType
}

// resolveDocument attempts to find a document by trying all known types.
func (h *handler) resolveDocument(r *http.Request, id string) (*resolvedDoc, error) {
	for _, dt := range []document.DocumentType{
		document.DocumentTypeEstimate,
		document.DocumentTypeMandate,
		document.DocumentTypeContract,
		document.DocumentTypeInvoice,
	} {
		_, err := h.service.GetDocument(r.Context(), id, dt)
		if err == nil {
			return &resolvedDoc{docType: dt}, nil
		}
	}
	return nil, fmt.Errorf("document %s not found", id)
}