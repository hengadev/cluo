package mediaHandler

import (
	"net/http"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	// router.HandleFunc("POST /case/{caseID}/upload", h.UploadMedia)
	// router.HandleFunc("GET /case/{caseID}", h.ListMediaByCaseID)
	// router.HandleFunc("delete /{caseID}", h.DeleteMediaByID)
}
