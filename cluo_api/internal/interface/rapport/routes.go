package rapportHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/contracts/identity"
	mw "github.com/hengadev/cluo_api/internal/common/middleware"
)

func (h *handler) RegisterRoutes(mux *http.ServeMux) {
	RequireAdministrator := h.authmw.RequireMinimumRole(identity.Administrator)

	mux.HandleFunc("POST /cases/{id}/rapport", RequireAdministrator(mw.EnableCORS(h.CreateRapport)))
	mux.HandleFunc("GET /cases/{id}/rapport", RequireAdministrator(mw.EnableCORS(h.GetRapport)))
	mux.HandleFunc("PUT /cases/{id}/rapport", RequireAdministrator(mw.EnableCORS(h.UpdateRapport)))
	mux.HandleFunc("DELETE /cases/{id}/rapport", RequireAdministrator(mw.EnableCORS(h.DeleteRapport)))
}
