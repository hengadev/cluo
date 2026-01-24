package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Handler returns an HTTP handler for the /metrics endpoint.
func Handler() http.Handler {
	return promhttp.Handler()
}

// RegisterRoutes registers the metrics endpoint on the given mux.
func RegisterRoutes(mux *http.ServeMux, path string) {
	if path == "" {
		path = "/metrics"
	}
	mux.Handle("GET "+path, Handler())
}
