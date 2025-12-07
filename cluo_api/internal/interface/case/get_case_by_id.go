package caseHandler

import (
	"net/http"

	"github.com/google/uuid"
	caseDomain "github.com/hengadev/cluo_api/internal/domain/case"
)

func (h *handler) GetCaseByID(w http.ResponseWriter, r *http.Request) {
	// Extract case ID from URL path
	caseIDStr := r.PathValue("id")
	if caseIDStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Parse case ID as UUID
	caseID, err := uuid.Parse(caseIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Create request object
	request := caseDomain.GetCaseByIDRequest{
		ID: caseID,
	}

	// Call service layer
	response, err := h.svc.GetCaseByID(r.Context(), &request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return response
	w.WriteHeader(http.StatusOK)
	response.ServeHTTP(w, r)
}