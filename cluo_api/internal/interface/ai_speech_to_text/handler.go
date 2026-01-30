package aiSpeechToTextHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/middleware/auth"
	"github.com/hengadev/cluo_api/internal/ports"
)

type Handler interface {
	// Job management
	SubmitJob(w http.ResponseWriter, r *http.Request)
	GetJobStatus(w http.ResponseWriter, r *http.Request)
	CancelJob(w http.ResponseWriter, r *http.Request)
	ListJobs(w http.ResponseWriter, r *http.Request)

	// Transcription access
	GetTranscription(w http.ResponseWriter, r *http.Request)
	ListTranscriptions(w http.ResponseWriter, r *http.Request)
	DeleteTranscription(w http.ResponseWriter, r *http.Request)

	// Health check
	HealthCheck(w http.ResponseWriter, r *http.Request)

	RegisterRoutes(router *http.ServeMux)
}

type handler struct {
	svc    ports.SpeechToTextService
	authmw auth.AuthMiddleware
}

func New(svc ports.SpeechToTextService, authmw auth.AuthMiddleware) Handler {
	return &handler{svc: svc, authmw: authmw}
}
