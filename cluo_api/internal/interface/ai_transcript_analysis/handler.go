package aiTranscriptAnalysisHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/middleware/auth"
	"github.com/hengadev/cluo_api/internal/ports"
)

type Handler interface {
	AnalyzeTranscript(w http.ResponseWriter, r *http.Request)
	GetAnalysis(w http.ResponseWriter, r *http.Request)
	HealthCheck(w http.ResponseWriter, r *http.Request)
	RegisterRoutes(router *http.ServeMux)
}

type handler struct {
	svc    ports.TranscriptAnalysisService
	authmw auth.AuthMiddleware
}

func New(svc ports.TranscriptAnalysisService, authmw auth.AuthMiddleware) Handler {
	return &handler{svc: svc, authmw: authmw}
}
