package mediaHandler

import (
	"net/http"

	// "github.com/hengadev/cluo_api/internal/common/contracts/identity"
	mw "github.com/hengadev/cluo_api/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	// RequireAuthenticated := h.authmw.RequireMinimumRole(identity.Client)

	// Upload media (authenticated users)
	router.HandleFunc("POST "+UploadMediaEndpoint, RequireAuthenticated(mw.EnableCORS(h.UploadMedia)))

	// Get media by ID (authenticated users)
	router.HandleFunc("GET "+GetMediaByIDEndpoint, RequireAuthenticated(mw.EnableCORS(h.GetMediaByID)))

	// Update media (authenticated users)
	router.HandleFunc("PATCH "+UpdateMediaEndpoint, RequireAuthenticated(mw.EnableCORS(h.UpdateMedia)))

	// Delete media (authenticated users)
	router.HandleFunc("DELETE "+DeleteMediaEndpoint, RequireAuthenticated(mw.EnableCORS(h.DeleteMedia)))

	// List media by case (authenticated users)
	router.HandleFunc("GET "+ListMediaByCaseEndpoint, RequireAuthenticated(mw.EnableCORS(h.ListMediaByCaseID)))
}
