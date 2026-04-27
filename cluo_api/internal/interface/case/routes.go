package caseHandler

import (
	"net/http"

	// "github.com/hengadev/cluo_api/internal/common/contracts/identity"
	mw "github.com/hengadev/cluo_api/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	// RequireAdministrator := h.authmw.RequireMinimumRole(identity.Administrator)

	// Individual case operations
	// TODO: Re-enable auth middleware when ready
	// router.HandleFunc("POST /cases", RequireAdministrator(mw.EnableCORS(h.CreateCase)))
	router.HandleFunc("POST /cases", mw.EnableCORS(h.CreateCase))
	// router.HandleFunc("GET /cases/{id}", RequireAdministrator(mw.EnableCORS(h.GetCaseByID)))
	router.HandleFunc("GET /cases/{id}", mw.EnableCORS(h.GetCaseByID))
	// router.HandleFunc("PATCH /cases/{id}", RequireAdministrator(mw.EnableCORS(h.UpdateCase)))
	router.HandleFunc("PATCH /cases/{id}", mw.EnableCORS(h.UpdateCase))
	// router.HandleFunc("DELETE /cases/{id}", RequireAdministrator(mw.EnableCORS(h.DeleteCase)))
	router.HandleFunc("DELETE /cases/{id}", mw.EnableCORS(h.DeleteCase))

	// List operations
	// router.HandleFunc("GET /cases", RequireAdministrator(mw.EnableCORS(h.ListCases)))
	router.HandleFunc("GET /cases", mw.EnableCORS(h.ListCases))
	// router.HandleFunc("GET /clients/{clientId}/cases", RequireAdministrator(mw.EnableCORS(h.ListCasesByClient)))
	router.HandleFunc("GET /clients/{clientId}/cases", mw.EnableCORS(h.ListCasesByClient))

	// Legacy/placeholder routes (commented out)
	// router.HandleFunc("POST /cases/{id}", h.DeliverCase)
}
