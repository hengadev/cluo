package tokenHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/contracts/identity"
	mw "github.com/hengadev/cluo_api/internal/common/middleware"
)

func (h *handler) RegisterRoutes(mux *http.ServeMux) {
	RequireAdministrator := h.authmw.RequireMinimumRole(identity.Administrator)

	// PI routes (authenticated)
	mux.HandleFunc("POST /cases/{id}/tokens", RequireAdministrator(mw.EnableCORS(h.CreateToken)))
	mux.HandleFunc("GET /cases/{id}/tokens", RequireAdministrator(mw.EnableCORS(h.ListTokens)))
	mux.HandleFunc("DELETE /cases/{id}/tokens/{tokenId}", RequireAdministrator(mw.EnableCORS(h.RevokeToken)))

	// Portal routes (public, token in URL path)
	mux.HandleFunc("GET /token/{token}", h.ValidateToken)
	mux.HandleFunc("GET /token/{token}/media", h.GetAllMediaByToken)
	mux.HandleFunc("GET /token/{token}/media/{mediaId}", h.GetMediaByIDByToken)
	mux.HandleFunc("GET /token/{token}/report", h.GetReportByToken)
	mux.HandleFunc("GET /token/{token}/report/html", h.GetReportHTMLByToken)
	mux.HandleFunc("GET /token/{token}/documents", h.GetDocumentsByToken)
}
