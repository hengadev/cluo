package tokenHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/middleware/auth"
	"github.com/hengadev/cluo_api/internal/ports"
)

type Handler interface {
	RegisterRoutes(router *http.ServeMux)
}

type handler struct {
	svc        ports.TokenService
	rapportSvc ports.RapportService
	authmw     auth.AuthMiddleware
}

func New(svc ports.TokenService, rapportSvc ports.RapportService, authmw auth.AuthMiddleware) Handler {
	return &handler{svc: svc, rapportSvc: rapportSvc, authmw: authmw}
}
