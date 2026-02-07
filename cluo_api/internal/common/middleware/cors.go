package middleware

import (
	"net/http"
	"os"
)

// allowedOrigin returns the configured CORS origin from CLUO_CORS_ORIGIN,
// defaulting to http://localhost:5173 for development.
func allowedOrigin() string {
	if origin := os.Getenv("CLUO_CORS_ORIGIN"); origin != "" {
		return origin
	}
	return "http://localhost:5173"
}

func EnableCORS(next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin())
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CLIENT-IP")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}
