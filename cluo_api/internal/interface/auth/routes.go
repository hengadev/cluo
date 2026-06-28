package authHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/contracts/identity"
	mw "github.com/hengadev/cluo_api/internal/common/middleware"
)

func (h *AuthHandler) RegisterRoutes(router *http.ServeMux) {
	RequireClient := h.authmw.RequireMinimumRole(identity.Client)
	RequireRefreshToken := h.authmw.RequireRefreshToken

	loginHandler := mw.EnableCORS(h.Login)
	if h.loginRateLimiter != nil {
		loginHandler = handlerToFunc(h.loginRateLimiter(funcToHandler(loginHandler)))
	}

	router.HandleFunc("POST /auth/login", loginHandler)
	router.HandleFunc("POST /auth/register", mw.EnableCORS(h.Register))
	router.HandleFunc("POST /auth/logout", RequireClient(mw.EnableCORS(h.SignOut)))
	router.HandleFunc("POST /auth/refresh", RequireRefreshToken(mw.EnableCORS(h.RefreshSession)))
	router.HandleFunc("GET /auth/me", RequireClient(mw.EnableCORS(h.GetCurrentUser)))
	router.HandleFunc("PATCH /auth/me", RequireClient(mw.EnableCORS(h.UpdateCurrentUser)))
}
