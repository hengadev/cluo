package searchHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/contracts/identity"
	mw "github.com/hengadev/cluo_api/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	RequireAdministrator := h.authmw.RequireMinimumRole(identity.Administrator)
	router.HandleFunc("GET /search", RequireAdministrator(mw.EnableCORS(h.Search)))
}
