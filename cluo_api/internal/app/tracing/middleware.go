package tracing

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Middleware returns an HTTP middleware that adds tracing to requests.
func Middleware(serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return otelhttp.NewHandler(next, serviceName,
			otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
				return r.Method + " " + r.URL.Path
			}),
		)
	}
}

// WrapHandler wraps an http.Handler with tracing.
func WrapHandler(handler http.Handler, operation string) http.Handler {
	return otelhttp.NewHandler(handler, operation)
}

// WrapRoundTripper wraps an http.RoundTripper with tracing for outgoing requests.
func WrapRoundTripper(rt http.RoundTripper) http.RoundTripper {
	return otelhttp.NewTransport(rt)
}
