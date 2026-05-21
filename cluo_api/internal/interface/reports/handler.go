package reportHandler

import "net/http"

// Deprecated: use internal/interface/rapport instead.
// This stub is preserved for historical reference only.

type Handler interface {
	RegisterRoutes(router *http.ServeMux)
}
