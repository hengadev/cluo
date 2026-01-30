package aiSpeechToTextHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/contracts/identity"
	mw "github.com/hengadev/cluo_api/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	RequireAdministrator := h.authmw.RequireMinimumRole(identity.Administrator)

	// Job management endpoints (require admin role)
	router.HandleFunc("POST /ai/speech/jobs", RequireAdministrator(mw.EnableCORS(h.SubmitJob)))
	router.HandleFunc("GET /ai/speech/jobs/{id}", RequireAdministrator(mw.EnableCORS(h.GetJobStatus)))
	router.HandleFunc("DELETE /ai/speech/jobs/{id}", RequireAdministrator(mw.EnableCORS(h.CancelJob)))
	router.HandleFunc("GET /ai/speech/jobs", RequireAdministrator(mw.EnableCORS(h.ListJobs)))

	// Transcription access endpoints (require admin role)
	router.HandleFunc("GET /ai/speech/transcriptions/{id}", RequireAdministrator(mw.EnableCORS(h.GetTranscription)))
	router.HandleFunc("DELETE /ai/speech/transcriptions/{id}", RequireAdministrator(mw.EnableCORS(h.DeleteTranscription)))
	router.HandleFunc("GET /ai/speech/transcriptions", RequireAdministrator(mw.EnableCORS(h.ListTranscriptions)))

	// Health check (public)
	router.HandleFunc("GET /ai/speech/health", mw.EnableCORS(h.HealthCheck))
}
