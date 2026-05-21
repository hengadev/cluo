package caseTypeHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/contracts/identity"
	mw "github.com/hengadev/cluo_api/internal/common/middleware"
)

func (h *handler) RegisterRoutes(mux *http.ServeMux) {
	RequireAdministrator := h.authmw.RequireMinimumRole(identity.Administrator)

	mux.HandleFunc("POST /case-types", RequireAdministrator(mw.EnableCORS(h.CreateCaseType)))
	mux.HandleFunc("GET /case-types", RequireAdministrator(mw.EnableCORS(h.ListCaseTypes)))
	mux.HandleFunc("GET /case-types/{id}", RequireAdministrator(mw.EnableCORS(h.GetCaseType)))
	mux.HandleFunc("PATCH /case-types/{id}", RequireAdministrator(mw.EnableCORS(h.UpdateCaseType)))
	mux.HandleFunc("DELETE /case-types/{id}", RequireAdministrator(mw.EnableCORS(h.DeleteCaseType)))
}
