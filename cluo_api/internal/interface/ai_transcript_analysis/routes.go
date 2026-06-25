package aiTranscriptAnalysisHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/contracts/identity"
	mw "github.com/hengadev/cluo_api/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	RequireAdministrator := h.authmw.RequireMinimumRole(identity.Administrator)

	// Transcript analysis endpoints (require admin role)
	router.HandleFunc("POST /ai/analysis/analyze", RequireAdministrator(mw.EnableCORS(h.AnalyzeTranscript)))
	router.HandleFunc("GET /ai/analysis/{id}", RequireAdministrator(mw.EnableCORS(h.GetAnalysis)))
	router.HandleFunc("GET /ai/analysis/by-transcription/{transcriptionId}", RequireAdministrator(mw.EnableCORS(h.GetAnalysisByTranscriptionID)))

	// Health check (public)
	router.HandleFunc("GET /ai/analysis/health", mw.EnableCORS(h.HealthCheck))
}
