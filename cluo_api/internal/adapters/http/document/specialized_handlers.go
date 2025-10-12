package document

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/hengadev/cluo_api/internal/domain"
	"github.com/hengadev/cluo_api/internal/ports"
)

// Estimate handlers

func (h *Handler) CreateEstimate(w http.ResponseWriter, r *http.Request) {
	var estimate domain.Estimate
	if err := json.NewDecoder(r.Body).Decode(&estimate); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// TODO: Validate request
	// TODO: Get user ID from context

	doc, err := h.service.CreateEstimate(r.Context(), &estimate)
	if err != nil {
		if _, ok := err.(*ValidationError); ok {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to create estimate")
		return
	}

	h.writeSuccess(w, doc)
}

func (h *Handler) AcceptEstimate(w http.ResponseWriter, r *http.Request) {
	estimateID := chi.URLParam(r, "id")
	if estimateID == "" {
		h.writeError(w, http.StatusBadRequest, "Estimate ID is required")
		return
	}

	var req struct {
		AcceptedBy string `json:"accepted_by"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.AcceptedBy == "" {
		h.writeError(w, http.StatusBadRequest, "Accepted by is required")
		return
	}

	// TODO: Validate request
	// TODO: Get user ID from context

	mandate, err := h.service.AcceptEstimate(r.Context(), estimateID, req.AcceptedBy)
	if err != nil {
		if err.Error() == "estimate not found" {
			h.writeError(w, http.StatusNotFound, "Estimate not found")
			return
		}
		if _, ok := err.(*ValidationError); ok {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to accept estimate")
		return
	}

	h.writeSuccess(w, mandate)
}

func (h *Handler) UpdateEstimate(w http.ResponseWriter, r *http.Request) {
	estimateID := chi.URLParam(r, "id")
	if estimateID == "" {
		h.writeError(w, http.StatusBadRequest, "Estimate ID is required")
		return
	}

	var req struct {
		LineItems []domain.EstimateItem `json:"line_items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(req.LineItems) == 0 {
		h.writeError(w, http.StatusBadRequest, "At least one line item is required")
		return
	}

	// TODO: Validate request
	// TODO: Get user ID from context

	estimate, err := h.service.UpdateEstimate(r.Context(), estimateID, req.LineItems)
	if err != nil {
		if err.Error() == "estimate not found" {
			h.writeError(w, http.StatusNotFound, "Estimate not found")
			return
		}
		if _, ok := err.(*ValidationError); ok {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to update estimate")
		return
	}

	h.writeSuccess(w, estimate)
}

// Mandate handlers

func (h *Handler) CreateMandate(w http.ResponseWriter, r *http.Request) {
	var mandate domain.Mandate
	if err := json.NewDecoder(r.Body).Decode(&mandate); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// TODO: Validate request
	// TODO: Get user ID from context

	doc, err := h.service.CreateMandate(r.Context(), &mandate)
	if err != nil {
		if _, ok := err.(*ValidationError); ok {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to create mandate")
		return
	}

	h.writeSuccess(w, doc)
}

func (h *Handler) SignMandate(w http.ResponseWriter, r *http.Request) {
	mandateID := chi.URLParam(r, "id")
	if mandateID == "" {
		h.writeError(w, http.StatusBadRequest, "Mandate ID is required")
		return
	}

	var req domain.SignDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// TODO: Validate request
	// TODO: Get user ID from context

	mandate, err := h.service.SignMandate(r.Context(), mandateID, &req)
	if err != nil {
		if err.Error() == "mandate not found" {
			h.writeError(w, http.StatusNotFound, "Mandate not found")
			return
		}
		if _, ok := err.(*ValidationError); ok {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to sign mandate")
		return
	}

	h.writeSuccess(w, mandate)
}

func (h *Handler) ActivateMandate(w http.ResponseWriter, r *http.Request) {
	mandateID := chi.URLParam(r, "id")
	if mandateID == "" {
		h.writeError(w, http.StatusBadRequest, "Mandate ID is required")
		return
	}

	// TODO: Get user ID from context

	mandate, err := h.service.ActivateMandate(r.Context(), mandateID)
	if err != nil {
		if err.Error() == "mandate not found" {
			h.writeError(w, http.StatusNotFound, "Mandate not found")
			return
		}
		if _, ok := err.(*ValidationError); ok {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to activate mandate")
		return
	}

	h.writeSuccess(w, mandate)
}

func (h *Handler) CreateContractFromMandate(w http.ResponseWriter, r *http.Request) {
	mandateID := chi.URLParam(r, "id")
	if mandateID == "" {
		h.writeError(w, http.StatusBadRequest, "Mandate ID is required")
		return
	}

	var contract domain.Contract
	if err := json.NewDecoder(r.Body).Decode(&contract); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// TODO: Validate request
	// TODO: Get user ID from context

	createdContract, err := h.service.CreateContractFromMandate(r.Context(), mandateID, &contract)
	if err != nil {
		if err.Error() == "mandate not found" {
			h.writeError(w, http.StatusNotFound, "Mandate not found")
			return
		}
		if _, ok := err.(*ValidationError); ok {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to create contract from mandate")
		return
	}

	h.writeSuccess(w, createdContract)
}

// Contract handlers

func (h *Handler) CreateContract(w http.ResponseWriter, r *http.Request) {
	var contract domain.Contract
	if err := json.NewDecoder(r.Body).Decode(&contract); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// TODO: Validate request
	// TODO: Get user ID from context

	doc, err := h.service.CreateContract(r.Context(), &contract)
	if err != nil {
		if _, ok := err.(*ValidationError); ok {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to create contract")
		return
	}

	h.writeSuccess(w, doc)
}

func (h *Handler) SignContract(w http.ResponseWriter, r *http.Request) {
	contractID := chi.URLParam(r, "id")
	if contractID == "" {
		h.writeError(w, http.StatusBadRequest, "Contract ID is required")
		return
	}

	var req domain.SignDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// TODO: Validate request
	// TODO: Get user ID from context

	contract, err := h.service.SignContract(r.Context(), contractID, &req)
	if err != nil {
		if err.Error() == "contract not found" {
			h.writeError(w, http.StatusNotFound, "Contract not found")
			return
		}
		if _, ok := err.(*ValidationError); ok {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to sign contract")
		return
	}

	h.writeSuccess(w, contract)
}

func (h *Handler) ActivateContract(w http.ResponseWriter, r *http.Request) {
	contractID := chi.URLParam(r, "id")
	if contractID == "" {
		h.writeError(w, http.StatusBadRequest, "Contract ID is required")
		return
	}

	// TODO: Get user ID from context

	contract, err := h.service.ActivateContract(r.Context(), contractID)
	if err != nil {
		if err.Error() == "contract not found" {
			h.writeError(w, http.StatusNotFound, "Contract not found")
			return
		}
		if _, ok := err.(*ValidationError); ok {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to activate contract")
		return
	}

	h.writeSuccess(w, contract)
}

func (h *Handler) CreateInvoiceFromContract(w http.ResponseWriter, r *http.Request) {
	contractID := chi.URLParam(r, "id")
	if contractID == "" {
		h.writeError(w, http.StatusBadRequest, "Contract ID is required")
		return
	}

	var invoice domain.Invoice
	if err := json.NewDecoder(r.Body).Decode(&invoice); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// TODO: Validate request
	// TODO: Get user ID from context

	createdInvoice, err := h.service.CreateInvoiceFromContract(r.Context(), contractID, &invoice)
	if err != nil {
		if err.Error() == "contract not found" {
			h.writeError(w, http.StatusNotFound, "Contract not found")
			return
		}
		if _, ok := err.(*ValidationError); ok {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to create invoice from contract")
		return
	}

	h.writeSuccess(w, createdInvoice)
}

// Invoice handlers

func (h *Handler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	var invoice domain.Invoice
	if err := json.NewDecoder(r.Body).Decode(&invoice); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// TODO: Validate request
	// TODO: Get user ID from context

	doc, err := h.service.CreateInvoice(r.Context(), &invoice)
	if err != nil {
		if _, ok := err.(*ValidationError); ok {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to create invoice")
		return
	}

	h.writeSuccess(w, doc)
}

func (h *Handler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	invoiceID := chi.URLParam(r, "id")
	if invoiceID == "" {
		h.writeError(w, http.StatusBadRequest, "Invoice ID is required")
		return
	}

	var req domain.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// TODO: Validate request
	// TODO: Get user ID from context

	invoice, err := h.service.ProcessPayment(r.Context(), invoiceID, &req)
	if err != nil {
		if err.Error() == "invoice not found" {
			h.writeError(w, http.StatusNotFound, "Invoice not found")
			return
		}
		if _, ok := err.(*ValidationError); ok {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to process payment")
		return
	}

	h.writeSuccess(w, invoice)
}

func (h *Handler) VoidInvoice(w http.ResponseWriter, r *http.Request) {
	invoiceID := chi.URLParam(r, "id")
	if invoiceID == "" {
		h.writeError(w, http.StatusBadRequest, "Invoice ID is required")
		return
	}

	// TODO: Get user ID from context

	invoice, err := h.service.VoidInvoice(r.Context(), invoiceID)
	if err != nil {
		if err.Error() == "invoice not found" {
			h.writeError(w, http.StatusNotFound, "Invoice not found")
			return
		}
		if _, ok := err.(*ValidationError); ok {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to void invoice")
		return
	}

	h.writeSuccess(w, invoice)
}

func (h *Handler) GetOverdueInvoices(w http.ResponseWriter, r *http.Request) {
	pagination, err := h.getPaginationFromRequest(r)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	invoices, total, err := h.service.GetOverdueInvoices(r.Context(), pagination)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get overdue invoices")
		return
	}

	response := struct {
		Success bool              `json:"success"`
		Data    []*domain.Invoice `json:"data"`
		Total   int               `json:"total"`
		Page    int               `json:"page"`
		PerPage int               `json:"per_page"`
	}{
		Success: true,
		Data:    invoices,
		Total:   total,
		Page:    pagination.Page,
		PerPage: pagination.PageSize,
	}

	h.writeJSON(w, http.StatusOK, response)
}