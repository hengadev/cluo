package caseSubjectHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/contracts/identity"
	mw "github.com/hengadev/cluo_api/internal/common/middleware"
)

func (h *handler) RegisterRoutes(mux *http.ServeMux) {
	RequireAdministrator := h.authmw.RequireMinimumRole(identity.Administrator)

	mux.HandleFunc("POST /subjects", RequireAdministrator(mw.EnableCORS(h.CreateCaseSubject)))
	mux.HandleFunc("GET /subjects", RequireAdministrator(mw.EnableCORS(h.ListCaseSubjects)))
	mux.HandleFunc("GET /subjects/{id}", RequireAdministrator(mw.EnableCORS(h.GetCaseSubject)))
	mux.HandleFunc("PATCH /subjects/{id}", RequireAdministrator(mw.EnableCORS(h.UpdateCaseSubject)))
	mux.HandleFunc("DELETE /subjects/{id}", RequireAdministrator(mw.EnableCORS(h.DeleteCaseSubject)))
}
