package tokenHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/contracts/identity"
	mw "github.com/hengadev/cluo_api/internal/common/middleware"
)

func (h *TokenHandler) RegisterRoutes(mux *http.ServeMux) {
	RequireAdministrator := h.authmw.RequireMinimumRole(identity.Administrator)

	// PI routes (authenticated)
	mux.HandleFunc("POST /cases/{id}/tokens", RequireAdministrator(mw.EnableCORS(h.CreateToken)))
	mux.HandleFunc("GET /cases/{id}/tokens", RequireAdministrator(mw.EnableCORS(h.ListTokens)))
	mux.HandleFunc("DELETE /cases/{id}/tokens/{tokenId}", RequireAdministrator(mw.EnableCORS(h.RevokeToken)))

	// Portal routes (public, token in URL path or query param)
	// Apply rate limiter to both public token routes.
	validateHandler := h.ValidateToken
	validateQueryHandler := mw.Handler(h.ValidateTokenQuery)
	if h.tokenRateLimiter != nil {
		validateHandler = handlerToFunc(h.tokenRateLimiter(funcToHandler(h.ValidateToken)))
		validateQueryHandler = handlerToFunc(h.tokenRateLimiter(funcToHandler(h.ValidateTokenQuery)))
	}

	mux.HandleFunc("GET /tokens/validate", validateQueryHandler)
	mux.HandleFunc("GET /token/{token}", validateHandler)

	// Other portal routes – not rate limited per issue spec
	mux.HandleFunc("GET /token/{token}/media", h.GetAllMediaByToken)
	mux.HandleFunc("GET /token/{token}/media/{mediaId}", h.GetMediaByIDByToken)
	mux.HandleFunc("GET /token/{token}/report", h.GetReportByToken)
	mux.HandleFunc("GET /token/{token}/report/html", h.GetReportHTMLByToken)
	mux.HandleFunc("GET /token/{token}/report/pdf", h.GetReportPDFByToken)
	mux.HandleFunc("GET /token/{token}/documents", h.GetDocumentsByToken)
	mux.HandleFunc("GET /token/{token}/documents/{type}/pdf", h.GetDocumentPDFByToken)
	mux.HandleFunc("GET /token/{token}/archive", h.GetFullArchiveByToken)
	mux.HandleFunc("GET /token/{token}/media/archive", h.GetMediaArchiveByToken)
}
