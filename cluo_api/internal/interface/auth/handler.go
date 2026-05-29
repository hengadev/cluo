package authHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/middleware/auth"
	mw "github.com/hengadev/cluo_api/internal/common/middleware"
	"github.com/hengadev/cluo_api/internal/ports"
)

type Handler interface {
	RegisterRoutes(router *http.ServeMux)
}

type AuthHandler struct {
	svc              ports.AuthService
	authmw           auth.AuthMiddleware
	loginRateLimiter func(http.Handler) http.Handler // optional, may be nil
}

func New(svc ports.AuthService, authmw auth.AuthMiddleware) Handler {
	return &AuthHandler{svc: svc, authmw: authmw}
}

// WithLoginRateLimiter returns a copy of the handler with the given rate limiter
// applied to the login route.
func (h *AuthHandler) WithLoginRateLimiter(rl func(http.Handler) http.Handler) Handler {
	return &AuthHandler{svc: h.svc, authmw: h.authmw, loginRateLimiter: rl}
}

// handlerToFunc converts an http.Handler to a Handler func.
func handlerToFunc(h http.Handler) mw.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
}

// funcToHandler converts a Handler func to an http.Handler.
func funcToHandler(fn mw.Handler) http.Handler {
	return http.HandlerFunc(fn)
}
