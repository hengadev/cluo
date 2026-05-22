package authHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/contracts/identity"
	mw "github.com/hengadev/cluo_api/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	RequireClient := h.authmw.RequireMinimumRole(identity.Client)
	RequireRefreshToken := h.authmw.RequireRefreshToken

	router.HandleFunc("POST /auth/login", mw.EnableCORS(h.Login))
	router.HandleFunc("POST /auth/register", mw.EnableCORS(h.Register))
	router.HandleFunc("POST /auth/logout", RequireClient(mw.EnableCORS(h.SignOut)))
	router.HandleFunc("POST /auth/refresh", RequireRefreshToken(mw.EnableCORS(h.RefreshSession)))
	router.HandleFunc("GET /auth/me", RequireClient(mw.EnableCORS(h.GetCurrentUser)))
}
