package rapportHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/middleware/auth"
	"github.com/hengadev/cluo_api/internal/ports"
)

type Handler interface {
	RegisterRoutes(mux *http.ServeMux)
}

type handler struct {
	svc    ports.RapportService
	authmw auth.AuthMiddleware
}

func New(svc ports.RapportService, authmw auth.AuthMiddleware) Handler {
	return &handler{svc: svc, authmw: authmw}
}
