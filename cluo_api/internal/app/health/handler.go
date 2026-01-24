package health

import (
	"encoding/json"
	"net/http"
)

// Handler handles health check HTTP requests.
type Handler struct {
	checker *Checker
}

// NewHandler creates a new health handler.
func NewHandler(checker *Checker) *Handler {
	return &Handler{checker: checker}
}

// RegisterRoutes registers health check routes on the given mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", h.Liveness)
	mux.HandleFunc("GET /ready", h.Readiness)
}

// Liveness handles the /health endpoint (liveness probe).
// Returns 200 OK if the server is running.
func (h *Handler) Liveness(w http.ResponseWriter, r *http.Request) {
	report := h.checker.CheckLiveness()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(report)
}

// Readiness handles the /ready endpoint (readiness probe).
// Returns 200 OK if all dependencies are healthy, 503 otherwise.
func (h *Handler) Readiness(w http.ResponseWriter, r *http.Request) {
	report := h.checker.CheckReadiness(r.Context())

	w.Header().Set("Content-Type", "application/json")

	switch report.Status {
	case StatusHealthy:
		w.WriteHeader(http.StatusOK)
	case StatusDegraded:
		w.WriteHeader(http.StatusOK) // Still accepting traffic but degraded
	case StatusUnhealthy:
		w.WriteHeader(http.StatusServiceUnavailable)
	default:
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	json.NewEncoder(w).Encode(report)
}
