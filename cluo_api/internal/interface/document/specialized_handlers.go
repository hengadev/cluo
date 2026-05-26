package document

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/hengadev/cluo_api/internal/common/auth/session"
	"github.com/hengadev/cluo_api/internal/domain/document"
)

// ---------------------------------------------------------------------------
// Helper: extract authenticated user UUID from session context
// ---------------------------------------------------------------------------

func (h *handler) getUserID(r *http.Request) (uuid.UUID, error) {
	sessionInfo, ok := session.SessionInfoFromContext(r.Context())
	if !ok || sessionInfo == nil {
		return uuid.Nil, fmt.Errorf("unauthorized")
	}
	return sessionInfo.UserID, nil
}

// ---------------------------------------------------------------------------
// Helper: parse a time field from a raw JSON value, returning a 400-friendly error
// ---------------------------------------------------------------------------

func parseTimeField(raw json.RawMessage, fieldName string) (time.Time, error) {
	if len(raw) == 0 {
		return time.Time{}, nil
	}
	// Trim quotes so we can parse the string directly
	var s string
	if err := json.Unmarshal(raw, &s); err != nil {
		return time.Time{}, fmt.Errorf("invalid %s: must be an RFC3339 date string", fieldName)
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid %s: must be a valid RFC3339 date string", fieldName)
	}
	return t, nil
}

// ---------------------------------------------------------------------------
// Estimate handlers
// ---------------------------------------------------------------------------

// estimateCreateRequest is the wire-format DTO for creating an estimate.
// Time fields are sent as RFC3339 strings and parsed explicitly so that
// malformed values produce a clear 400 error instead of a generic JSON error.
type estimateCreateRequest struct {
	CaseID         uuid.UUID           `json:"case_id"`
	ClientID       uuid.UUID           `json:"client_id"`
	EstimateNumber string              `json:"estimate_number"`
	IssueDate      json.RawMessage     `json:"issue_date"`
	ValidUntil     json.RawMessage     `json:"valid_until,omitempty"`
	LineItems      []document.EstimateItem `json:"line_items"`
	EstimatedTotal float64             `json:"estimated_total"`
	Notes          *string             `json:"notes,omitempty"`
}

func (h *handler) CreateEstimate(w http.ResponseWriter, r *http.Request) {
	if _, err := h.getUserID(r); err != nil {
		h.writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	var req estimateCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Parse date fields
	issueDate, err := parseTimeField(req.IssueDate, "issue_date")
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var validUntil *time.Time
	if len(req.ValidUntil) > 0 && string(req.ValidUntil) != "null" {
		vu, err := parseTimeField(req.ValidUntil, "valid_until")
		if err != nil {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		validUntil = &vu
	}

	estimate := &document.Estimate{
		DocumentBase:   document.NewDocumentBase(req.CaseID, req.ClientID),
		EstimateNumber: req.EstimateNumber,
		IssueDate:      issueDate,
		ValidUntil:     validUntil,
		LineItems:      req.LineItems,
		EstimatedTotal: req.EstimatedTotal,
		Notes:          req.Notes,
	}

	// Validation enforcement
	if err := estimate.Validate(); err != nil {
		h.writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	doc, err := h.service.CreateEstimate(r.Context(), estimate)
	if err != nil {
		if _, ok := err.(*ValidationError); ok {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to create estimate")
		return
	}

	h.writeCreated(w, doc)
}

func (h *handler) AcceptEstimate(w http.ResponseWriter, r *http.Request) {
	estimateID := r.PathValue("id")
	if estimateID == "" {
		h.writeError(w, http.StatusBadRequest, "Estimate ID is required")
		return
	}

	userID, err := h.getUserID(r)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	mandate, err := h.service.AcceptEstimate(r.Context(), estimateID, userID.String())
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

func (h *handler) UpdateEstimate(w http.ResponseWriter, r *http.Request) {
	estimateID := r.PathValue("id")
	if estimateID == "" {
		h.writeError(w, http.StatusBadRequest, "Estimate ID is required")
		return
	}

	if _, err := h.getUserID(r); err != nil {
		h.writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	var req struct {
		LineItems []document.EstimateItem `json:"line_items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(req.LineItems) == 0 {
		h.writeError(w, http.StatusBadRequest, "At least one line item is required")
		return
	}

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

// ---------------------------------------------------------------------------
// Mandate handlers
// ---------------------------------------------------------------------------

// mandateCreateRequest is the wire-format DTO for creating a mandate.
type mandateCreateRequest struct {
	CaseID           uuid.UUID       `json:"case_id"`
	ClientID         uuid.UUID       `json:"client_id"`
	MandateNumber    string          `json:"mandate_number"`
	IssueDate        json.RawMessage `json:"issue_date"`
	ScopeOfWork      string          `json:"scope_of_work"`
	ValidFrom        json.RawMessage `json:"valid_from"`
	ValidUntil       json.RawMessage `json:"valid_until,omitempty"`
	TermsConditions  string          `json:"terms_conditions"`
	SpecialInstructions *string      `json:"special_instructions,omitempty"`
	Jurisdiction     *string         `json:"jurisdiction,omitempty"`
}

func (h *handler) CreateMandate(w http.ResponseWriter, r *http.Request) {
	if _, err := h.getUserID(r); err != nil {
		h.writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	var req mandateCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Parse date fields
	issueDate, err := parseTimeField(req.IssueDate, "issue_date")
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	validFrom, err := parseTimeField(req.ValidFrom, "valid_from")
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var validUntil *time.Time
	if len(req.ValidUntil) > 0 && string(req.ValidUntil) != "null" {
		vu, err := parseTimeField(req.ValidUntil, "valid_until")
		if err != nil {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		validUntil = &vu
	}

	mandate := &document.Mandate{
		DocumentBase:       document.NewDocumentBase(req.CaseID, req.ClientID),
		MandateNumber:      req.MandateNumber,
		IssueDate:          issueDate,
		ScopeOfWork:        req.ScopeOfWork,
		ValidFrom:          validFrom,
		ValidUntil:         validUntil,
		TermsConditions:    req.TermsConditions,
		SpecialInstructions: req.SpecialInstructions,
		Jurisdiction:       req.Jurisdiction,
	}

	// Validation enforcement
	if err := mandate.Validate(); err != nil {
		h.writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	doc, err := h.service.CreateMandate(r.Context(), mandate)
	if err != nil {
		if _, ok := err.(*ValidationError); ok {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to create mandate")
		return
	}

	h.writeCreated(w, doc)
}

func (h *handler) SignMandate(w http.ResponseWriter, r *http.Request) {
	mandateID := r.PathValue("id")
	if mandateID == "" {
		h.writeError(w, http.StatusBadRequest, "Mandate ID is required")
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

func (h *handler) ActivateMandate(w http.ResponseWriter, r *http.Request) {
	mandateID := r.PathValue("id")
	if mandateID == "" {
		h.writeError(w, http.StatusBadRequest, "Mandate ID is required")
		return
	}

	if _, err := h.getUserID(r); err != nil {
		h.writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

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

func (h *handler) CreateContractFromMandate(w http.ResponseWriter, r *http.Request) {
	mandateID := r.PathValue("id")
	if mandateID == "" {
		h.writeError(w, http.StatusBadRequest, "Mandate ID is required")
		return
	}

	if _, err := h.getUserID(r); err != nil {
		h.writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	var contract document.Contract
	if err := json.NewDecoder(r.Body).Decode(&contract); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := contract.Validate(); err != nil {
		h.writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

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

// ---------------------------------------------------------------------------
// Contract handlers
// ---------------------------------------------------------------------------

// contractCreateRequest is the wire-format DTO for creating a contract.
type contractCreateRequest struct {
	CaseID           uuid.UUID       `json:"case_id"`
	ClientID         uuid.UUID       `json:"client_id"`
	ContractNumber   string          `json:"contract_number"`
	StartDate        json.RawMessage `json:"start_date"`
	EndDate          json.RawMessage `json:"end_date,omitempty"`
	ScopeOfServices  string          `json:"scope_of_services"`
	PaymentTerms     string          `json:"payment_terms"`
	Confidentiality  string          `json:"confidentiality"`
	TerminationClause string         `json:"termination_clause"`
	ContractValue    *float64        `json:"contract_value,omitempty"`
	Currency         *string         `json:"currency,omitempty"`
	RenewalTerms     *string         `json:"renewal_terms,omitempty"`
	GoverningLaw     *string         `json:"governing_law,omitempty"`
}

func (h *handler) CreateContract(w http.ResponseWriter, r *http.Request) {
	if _, err := h.getUserID(r); err != nil {
		h.writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	var req contractCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Parse date fields
	startDate, err := parseTimeField(req.StartDate, "start_date")
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var endDate *time.Time
	if len(req.EndDate) > 0 && string(req.EndDate) != "null" {
		ed, err := parseTimeField(req.EndDate, "end_date")
		if err != nil {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		endDate = &ed
	}

	contract := &document.Contract{
		DocumentBase:      document.NewDocumentBase(req.CaseID, req.ClientID),
		ContractNumber:    req.ContractNumber,
		StartDate:         startDate,
		EndDate:           endDate,
		ScopeOfServices:   req.ScopeOfServices,
		PaymentTerms:      req.PaymentTerms,
		Confidentiality:   req.Confidentiality,
		TerminationClause: req.TerminationClause,
		ContractValue:     req.ContractValue,
		Currency:          req.Currency,
		RenewalTerms:      req.RenewalTerms,
		GoverningLaw:      req.GoverningLaw,
	}

	// Validation enforcement
	if err := contract.Validate(); err != nil {
		h.writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	doc, err := h.service.CreateContract(r.Context(), contract)
	if err != nil {
		if _, ok := err.(*ValidationError); ok {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to create contract")
		return
	}

	h.writeCreated(w, doc)
}

func (h *handler) SignContract(w http.ResponseWriter, r *http.Request) {
	contractID := r.PathValue("id")
	if contractID == "" {
		h.writeError(w, http.StatusBadRequest, "Contract ID is required")
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

func (h *handler) ActivateContract(w http.ResponseWriter, r *http.Request) {
	contractID := r.PathValue("id")
	if contractID == "" {
		h.writeError(w, http.StatusBadRequest, "Contract ID is required")
		return
	}

	if _, err := h.getUserID(r); err != nil {
		h.writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

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

func (h *handler) CreateInvoiceFromContract(w http.ResponseWriter, r *http.Request) {
	contractID := r.PathValue("id")
	if contractID == "" {
		h.writeError(w, http.StatusBadRequest, "Contract ID is required")
		return
	}

	if _, err := h.getUserID(r); err != nil {
		h.writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	var invoice document.Invoice
	if err := json.NewDecoder(r.Body).Decode(&invoice); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := invoice.Validate(); err != nil {
		h.writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

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

// ---------------------------------------------------------------------------
// Invoice handlers
// ---------------------------------------------------------------------------

// invoiceCreateRequest is the wire-format DTO for creating an invoice.
type invoiceCreateRequest struct {
	CaseID         uuid.UUID          `json:"case_id"`
	ClientID       uuid.UUID          `json:"client_id"`
	InvoiceNumber  string             `json:"invoice_number"`
	IssueDate      json.RawMessage    `json:"issue_date"`
	DueDate        json.RawMessage    `json:"due_date"`
	LineItems      []document.InvoiceItem `json:"line_items"`
	TotalAmount    float64            `json:"total_amount"`
	TaxRate        float64            `json:"tax_rate"`
	TaxAmount      float64            `json:"tax_amount"`
	Notes          *string            `json:"notes,omitempty"`
}

func (h *handler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	if _, err := h.getUserID(r); err != nil {
		h.writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	var req invoiceCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// Parse date fields
	issueDate, err := parseTimeField(req.IssueDate, "issue_date")
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	dueDate, err := parseTimeField(req.DueDate, "due_date")
	if err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	invoice := &document.Invoice{
		DocumentBase:  document.NewDocumentBase(req.CaseID, req.ClientID),
		InvoiceNumber: req.InvoiceNumber,
		IssueDate:     issueDate,
		DueDate:       dueDate,
		LineItems:     req.LineItems,
		TotalAmount:   req.TotalAmount,
		TaxRate:       req.TaxRate,
		TaxAmount:     req.TaxAmount,
		Notes:         req.Notes,
		PaymentStatus: document.PaymentStatusUnpaid,
	}

	// Validation enforcement
	if err := invoice.Validate(); err != nil {
		h.writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	doc, err := h.service.CreateInvoice(r.Context(), invoice)
	if err != nil {
		if _, ok := err.(*ValidationError); ok {
			h.writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to create invoice")
		return
	}

	h.writeCreated(w, doc)
}

func (h *handler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	invoiceID := r.PathValue("id")
	if invoiceID == "" {
		h.writeError(w, http.StatusBadRequest, "Invoice ID is required")
		return
	}

	if _, err := h.getUserID(r); err != nil {
		h.writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	var req document.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

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

func (h *handler) VoidInvoice(w http.ResponseWriter, r *http.Request) {
	invoiceID := r.PathValue("id")
	if invoiceID == "" {
		h.writeError(w, http.StatusBadRequest, "Invoice ID is required")
		return
	}

	if _, err := h.getUserID(r); err != nil {
		h.writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

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

func (h *handler) GetOverdueInvoices(w http.ResponseWriter, r *http.Request) {
	if _, err := h.getUserID(r); err != nil {
		h.writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

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
		Data    []*document.Invoice `json:"data"`
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