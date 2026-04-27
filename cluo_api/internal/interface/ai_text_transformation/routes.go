package aiTextTransformationHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/contracts/identity"
	mw "github.com/hengadev/cluo_api/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	RequireAdministrator := h.authmw.RequireMinimumRole(identity.Administrator)

	// Text transformation endpoints (require admin role)
	// TODO: Re-enable auth middleware when ready
	// router.HandleFunc("POST /ai/text/transform", RequireAdministrator(mw.EnableCORS(h.TransformText)))
	router.HandleFunc("POST /ai/text/transform", mw.EnableCORS(h.TransformText))
	router.HandleFunc("GET /ai/text/health", mw.EnableCORS(h.HealthCheck))
}
