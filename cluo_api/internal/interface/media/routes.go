package mediaHandler

import (
	"net/http"

	// "github.com/hengadev/cluo_api/internal/common/contracts/identity"
	mw "github.com/hengadev/cluo_api/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	// RequireAuthenticated := h.authmw.RequireMinimumRole(identity.Client)

	// Upload media (authenticated users)
	// TODO: Re-enable auth middleware when ready
	// router.HandleFunc("POST "+UploadMediaEndpoint, RequireAuthenticated(mw.EnableCORS(h.UploadMedia)))
	router.HandleFunc("POST "+UploadMediaEndpoint, mw.EnableCORS(h.UploadMedia))

	// Get media by ID (authenticated users)
	// router.HandleFunc("GET "+GetMediaByIDEndpoint, RequireAuthenticated(mw.EnableCORS(h.GetMediaByID)))
	router.HandleFunc("GET "+GetMediaByIDEndpoint, mw.EnableCORS(h.GetMediaByID))

	// Update media (authenticated users)
	// router.HandleFunc("PATCH "+UpdateMediaEndpoint, RequireAuthenticated(mw.EnableCORS(h.UpdateMedia)))
	router.HandleFunc("PATCH "+UpdateMediaEndpoint, mw.EnableCORS(h.UpdateMedia))

	// Delete media (authenticated users)
	// router.HandleFunc("DELETE "+DeleteMediaEndpoint, RequireAuthenticated(mw.EnableCORS(h.DeleteMedia)))
	router.HandleFunc("DELETE "+DeleteMediaEndpoint, mw.EnableCORS(h.DeleteMedia))

	// List media by case (authenticated users)
	// router.HandleFunc("GET "+ListMediaByCaseEndpoint, RequireAuthenticated(mw.EnableCORS(h.ListMediaByCaseID)))
	router.HandleFunc("GET "+ListMediaByCaseEndpoint, mw.EnableCORS(h.ListMediaByCaseID))
}
