package document

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/hengadev/cluo_api/internal/domain"
	"github.com/hengadev/cluo_api/internal/ports"
)

type Handler struct {
	service ports.DocumentService
}

func New(service ports.DocumentService) *Handler {
	return &Handler{service: service}
}

// Response types
type DocumentResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type DocumentListResponse struct {
	Success bool                   `json:"success"`
	Data    []domain.DocumentSummary `json:"data"`
	Total   int                    `json:"total"`
	Page    int                    `json:"page"`
	PerPage int                    `json:"per_page"`
}

// Helper functions

func (h *Handler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, DocumentResponse{
		Success: false,
		Error:   message,
	})
}

func (h *Handler) writeSuccess(w http.ResponseWriter, data interface{}) {
	h.writeJSON(w, http.StatusOK, DocumentResponse{
		Success: true,
		Data:    data,
	})
}

func (h *Handler) getPaginationFromRequest(r *http.Request) (domain.Pagination, error) {
	pageStr := r.URL.Query().Get("page")
	perPageStr := r.URL.Query().Get("per_page")

	page := 1
	perPage := 20

	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p < 1 {
			return domain.Pagination{}, &ValidationError{Message: "invalid page parameter"}
		}
		page = p
	}

	if perPageStr != "" {
		pp, err := strconv.Atoi(perPageStr)
		if err != nil || pp < 1 || pp > 100 {
			return domain.Pagination{}, &ValidationError{Message: "invalid per_page parameter (must be 1-100)"}
		}
		perPage = pp
	}

	return domain.Pagination{
		Page:     page,
		PageSize: perPage,
	}, nil
}

func (h *Handler) getDocumentFilterFromRequest(r *http.Request) domain.DocumentFilter {
	filter := domain.DocumentFilter{}

	// Parse query parameters
	if docType := r.URL.Query().Get("type"); docType != "" {
		dt := domain.DocumentType(docType)
		filter.Type = &dt
	}

	if status := r.URL.Query().Get("status"); status != "" {
		ds := domain.DocumentStatus(status)
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

func (h *Handler) CreateDocument(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// TODO: Validate request
	// TODO: Get user ID from context

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

func (h *Handler) GetDocument(w http.ResponseWriter, r *http.Request) {
	documentID := chi.URLParam(r, "id")
	docType := chi.URLParam(r, "type")

	if documentID == "" || docType == "" {
		h.writeError(w, http.StatusBadRequest, "Document ID and type are required")
		return
	}

	doc, err := h.service.GetDocument(r.Context(), documentID, domain.DocumentType(docType))
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

func (h *Handler) UpdateDocument(w http.ResponseWriter, r *http.Request) {
	documentID := chi.URLParam(r, "id")
	docType := chi.URLParam(r, "type")

	if documentID == "" || docType == "" {
		h.writeError(w, http.StatusBadRequest, "Document ID and type are required")
		return
	}

	var req domain.UpdateDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// TODO: Validate request
	// TODO: Get user ID from context

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

func (h *Handler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	documentID := chi.URLParam(r, "id")
	docType := chi.URLParam(r, "type")

	if documentID == "" || docType == "" {
		h.writeError(w, http.StatusBadRequest, "Document ID and type are required")
		return
	}

	err := h.service.DeleteDocument(r.Context(), documentID, domain.DocumentType(docType))
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

func (h *Handler) ListDocuments(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) SendDocument(w http.ResponseWriter, r *http.Request) {
	documentID := chi.URLParam(r, "id")
	docType := chi.URLParam(r, "type")

	if documentID == "" || docType == "" {
		h.writeError(w, http.StatusBadRequest, "Document ID and type are required")
		return
	}

	var req domain.SendDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// TODO: Validate request
	// TODO: Get user ID from context

	err := h.service.SendDocument(r.Context(), documentID, domain.DocumentType(docType), &req)
	if err != nil {
		if err.Error() == "document not found" {
			h.writeError(w, http.StatusNotFound, "Document not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to send document")
		return
	}

	h.writeSuccess(w, map[string]string{"message": "Document sent successfully"})
}

func (h *Handler) SignDocument(w http.ResponseWriter, r *http.Request) {
	documentID := chi.URLParam(r, "id")
	docType := chi.URLParam(r, "type")

	if documentID == "" || docType == "" {
		h.writeError(w, http.StatusBadRequest, "Document ID and type are required")
		return
	}

	var req domain.SignDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// TODO: Validate request
	// TODO: Get user ID from context

	doc, err := h.service.SignDocument(r.Context(), documentID, domain.DocumentType(docType), &req)
	if err != nil {
		if err.Error() == "document not found" {
			h.writeError(w, http.StatusNotFound, "Document not found")
			return
		}
		if _, ok := err.(*ValidationError); ok {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to sign document")
		return
	}

	h.writeSuccess(w, doc)
}

func (h *Handler) ArchiveDocument(w http.ResponseWriter, r *http.Request) {
	documentID := chi.URLParam(r, "id")
	docType := chi.URLParam(r, "type")

	if documentID == "" || docType == "" {
		h.writeError(w, http.StatusBadRequest, "Document ID and type are required")
		return
	}

	err := h.service.ArchiveDocument(r.Context(), documentID, domain.DocumentType(docType))
	if err != nil {
		if err.Error() == "document not found" {
			h.writeError(w, http.StatusNotFound, "Document not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to archive document")
		return
	}

	h.writeSuccess(w, map[string]string{"message": "Document archived successfully"})
}

func (h *Handler) GetDocumentHistory(w http.ResponseWriter, r *http.Request) {
	documentID := chi.URLParam(r, "id")
	docType := chi.URLParam(r, "type")

	if documentID == "" || docType == "" {
		h.writeError(w, http.StatusBadRequest, "Document ID and type are required")
		return
	}

	pagination, err := h.getPaginationFromRequest(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	versions, total, err := h.service.GetDocumentHistory(r.Context(), documentID, domain.DocumentType(docType), pagination)
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
		Data      []*domain.DocumentVersion  `json:"data"`
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

func (h *Handler) GetDocumentWorkflow(w http.ResponseWriter, r *http.Request) {
	caseID := chi.URLParam(r, "caseId")
	if caseID == "" {
		h.writeError(w, http.StatusBadRequest, "Case ID is required")
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