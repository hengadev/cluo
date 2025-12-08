package caseHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/contracts/identity"
	mw "github.com/hengadev/cluo_api/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	RequireAdministrator := h.authmw.RequireMinimumRole(identity.Administrator)

	router.HandleFunc("POST /cases", RequireAdministrator(mw.EnableCORS(h.CreateCase)))
	router.HandleFunc("GET /cases/{id}", RequireAdministrator(mw.EnableCORS(h.GetCaseByID)))
	router.HandleFunc("PATCH /cases/{id}", RequireAdministrator(mw.EnableCORS(h.UpdateCase)))
	router.HandleFunc("DELETE /cases/{id}", RequireAdministrator(mw.EnableCORS(h.DeleteCase)))
	// router.HandleFunc("GET /cases", h.GetAllCases)
	// router.HandleFunc("POST /cases/{id}", h.DeliverCase)
}
