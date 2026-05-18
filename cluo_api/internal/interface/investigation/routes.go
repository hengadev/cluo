package investigationHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/contracts/identity"
	mw "github.com/hengadev/cluo_api/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	RequireAdministrator := h.authmw.RequireMinimumRole(identity.Administrator)

	// Individual case operations
	router.HandleFunc("POST /cases", RequireAdministrator(mw.EnableCORS(h.CreateCase)))
	router.HandleFunc("GET /cases/{id}", RequireAdministrator(mw.EnableCORS(h.GetCaseByID)))
	router.HandleFunc("PATCH /cases/{id}", RequireAdministrator(mw.EnableCORS(h.UpdateCase)))
	router.HandleFunc("DELETE /cases/{id}", RequireAdministrator(mw.EnableCORS(h.DeleteCase)))

	// List operations
	router.HandleFunc("GET /cases", RequireAdministrator(mw.EnableCORS(h.ListCases)))
	router.HandleFunc("GET /clients/{clientId}/cases", RequireAdministrator(mw.EnableCORS(h.ListCasesByClient)))

	// Legacy/placeholder routes (commented out)
	// router.HandleFunc("POST /cases/{id}", h.DeliverCase)
}
