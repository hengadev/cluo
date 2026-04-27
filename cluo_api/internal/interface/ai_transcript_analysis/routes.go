package aiTranscriptAnalysisHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/contracts/identity"
	mw "github.com/hengadev/cluo_api/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	RequireAdministrator := h.authmw.RequireMinimumRole(identity.Administrator)

	// Transcript analysis endpoints (require admin role)
	// TODO: Re-enable auth middleware when ready
	// router.HandleFunc("POST /ai/analysis/analyze", RequireAdministrator(mw.EnableCORS(h.AnalyzeTranscript)))
	router.HandleFunc("POST /ai/analysis/analyze", mw.EnableCORS(h.AnalyzeTranscript))
	// router.HandleFunc("GET /ai/analysis/{id}", RequireAdministrator(mw.EnableCORS(h.GetAnalysis)))
	router.HandleFunc("GET /ai/analysis/{id}", mw.EnableCORS(h.GetAnalysis))

	// Health check (public)
	router.HandleFunc("GET /ai/analysis/health", mw.EnableCORS(h.HealthCheck))
}
