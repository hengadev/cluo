package tokenHandler

import (
	"errors"
	"net/http"

	"github.com/hengadev/encx"
	"github.com/hengadev/cluo_api/internal/common/archive"
	mw "github.com/hengadev/cluo_api/internal/common/middleware"
	"github.com/hengadev/cluo_api/internal/common/middleware/auth"
	"github.com/hengadev/cluo_api/internal/ports"
)

var errMissingTokenParam = errors.New("missing required query parameter: token")

type validateTokenResponse struct {
	CaseID string `json:"caseId"`
}

type Handler interface {
	RegisterRoutes(router *http.ServeMux)
}

type TokenHandler struct {
	svc              ports.TokenService
	rapportSvc       ports.RapportService
	documentRepo     ports.DocumentRepository // may be nil while document packages are disabled
	crypto           encx.CryptoService      // may be nil while document packages are disabled
	authmw           auth.AuthMiddleware
	archiveDeps      archive.Dependencies // may be nil while storage is unavailable
	caseSvc          ports.CaseService   // used to resolve the case reference for archive filenames
	tokenRateLimiter func(http.Handler) http.Handler // optional, may be nil
}

func New(svc ports.TokenService, rapportSvc ports.RapportService, documentRepo ports.DocumentRepository, crypto encx.CryptoService, authmw auth.AuthMiddleware, caseSvc ports.CaseService) Handler {
	return &TokenHandler{svc: svc, rapportSvc: rapportSvc, documentRepo: documentRepo, crypto: crypto, authmw: authmw, caseSvc: caseSvc}
}

// NewWithArchive creates a handler with archive download support.
func NewWithArchive(svc ports.TokenService, rapportSvc ports.RapportService, documentRepo ports.DocumentRepository, crypto encx.CryptoService, authmw auth.AuthMiddleware, archiveDeps archive.Dependencies, caseSvc ports.CaseService) Handler {
	return &TokenHandler{svc: svc, rapportSvc: rapportSvc, documentRepo: documentRepo, crypto: crypto, authmw: authmw, archiveDeps: archiveDeps, caseSvc: caseSvc}
}

// WithTokenRateLimiter returns a copy of the handler with the given rate limiter
// applied to the public portal token routes.
func (h *TokenHandler) WithTokenRateLimiter(rl func(http.Handler) http.Handler) Handler {
	return &TokenHandler{svc: h.svc, rapportSvc: h.rapportSvc, documentRepo: h.documentRepo, crypto: h.crypto, authmw: h.authmw, archiveDeps: h.archiveDeps, caseSvc: h.caseSvc, tokenRateLimiter: rl}
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
