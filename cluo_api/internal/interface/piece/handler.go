package pieceHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/middleware/auth"
	"github.com/hengadev/cluo_api/internal/ports"
)

// Handler is the public interface for the piece HTTP handler.
type Handler interface {
	RegisterRoutes(mux *http.ServeMux)
}

type handler struct {
	svc    ports.PieceService
	authmw auth.AuthMiddleware
}

// New creates a new piece Handler.
func New(svc ports.PieceService, authmw auth.AuthMiddleware) Handler {
	return &handler{svc: svc, authmw: authmw}
}
