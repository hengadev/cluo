package tokenHandler

import (
	"net/http"

	"github.com/hengadev/encx"
	"github.com/hengadev/cluo_api/internal/common/archive"
	"github.com/hengadev/cluo_api/internal/common/middleware/auth"
	"github.com/hengadev/cluo_api/internal/ports"
)

type Handler interface {
	RegisterRoutes(router *http.ServeMux)
}

type handler struct {
	svc          ports.TokenService
	rapportSvc   ports.RapportService
	documentRepo ports.DocumentRepository // may be nil while document packages are disabled
	crypto       encx.CryptoService      // may be nil while document packages are disabled
	authmw       auth.AuthMiddleware
	archiveDeps  archive.Dependencies // may be nil while storage is unavailable
	caseSvc      ports.CaseService   // used to resolve the case reference for archive filenames
}

func New(svc ports.TokenService, rapportSvc ports.RapportService, documentRepo ports.DocumentRepository, crypto encx.CryptoService, authmw auth.AuthMiddleware, caseSvc ports.CaseService) Handler {
	return &handler{svc: svc, rapportSvc: rapportSvc, documentRepo: documentRepo, crypto: crypto, authmw: authmw, caseSvc: caseSvc}
}

// NewWithArchive creates a handler with archive download support.
func NewWithArchive(svc ports.TokenService, rapportSvc ports.RapportService, documentRepo ports.DocumentRepository, crypto encx.CryptoService, authmw auth.AuthMiddleware, archiveDeps archive.Dependencies, caseSvc ports.CaseService) Handler {
	return &handler{svc: svc, rapportSvc: rapportSvc, documentRepo: documentRepo, crypto: crypto, authmw: authmw, archiveDeps: archiveDeps, caseSvc: caseSvc}
}
