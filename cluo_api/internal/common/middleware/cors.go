package middleware

import (
	"net/http"
	"os"
	"strings"
)

// allowedOrigins returns the set of permitted CORS origins from CLUO_CORS_ORIGIN
// (comma-separated). Defaults to http://localhost:5173 for local development.
// For Wails desktop builds add "wails://wails"; for Wails dev mode add
// "http://localhost:34115".
func allowedOrigins() map[string]struct{} {
	raw := os.Getenv("CLUO_CORS_ORIGIN")
	if raw == "" {
		raw = "http://localhost:5173"
	}
	origins := make(map[string]struct{})
	for _, o := range strings.Split(raw, ",") {
		o = strings.TrimSpace(o)
		if o != "" {
			origins[o] = struct{}{}
		}
	}
	return origins
}

func EnableCORS(next Handler) Handler {
	origins := allowedOrigins()
	return func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w, r, origins)
		next(w, r)
	}
}

// PreflightMiddleware handles OPTIONS preflight requests at the top level,
// before Go 1.22 method-specific routing consumes them and returns 404.
// Must be applied as the outermost wrapper around the ServeMux.
func PreflightMiddleware(next http.Handler) http.Handler {
	origins := allowedOrigins()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			setCORSHeaders(w, r, origins)
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// GlobalCORSMiddleware adds CORS headers to every response, including 404s
// from unregistered routes. Apply it between PreflightMiddleware and the mux
// so that routes which aren't registered (because an optional service is
// disabled) still return a readable response instead of being CORS-blocked.
func GlobalCORSMiddleware(next http.Handler) http.Handler {
	origins := allowedOrigins()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w, r, origins)
		next.ServeHTTP(w, r)
	})
}

func setCORSHeaders(w http.ResponseWriter, r *http.Request, origins map[string]struct{}) {
	origin := r.Header.Get("Origin")
	if _, ok := origins[origin]; ok {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CLIENT-IP")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Vary", "Origin")
}
