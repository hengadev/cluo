package pieceHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/contracts/identity"
	mw "github.com/hengadev/cluo_api/internal/common/middleware"
)

func (h *handler) RegisterRoutes(mux *http.ServeMux) {
	RequireAdministrator := h.authmw.RequireMinimumRole(identity.Administrator)

	mux.HandleFunc("POST /cases/{id}/pieces", RequireAdministrator(mw.EnableCORS(h.UploadPiece)))
	mux.HandleFunc("GET /cases/{id}/pieces", RequireAdministrator(mw.EnableCORS(h.ListPieces)))
	mux.HandleFunc("GET /cases/{id}/pieces/{pieceId}", RequireAdministrator(mw.EnableCORS(h.GetPiece)))
	mux.HandleFunc("DELETE /cases/{id}/pieces/{pieceId}", RequireAdministrator(mw.EnableCORS(h.DeletePiece)))
}
