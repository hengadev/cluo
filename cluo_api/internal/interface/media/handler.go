package mediaHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/middleware/auth"
	"github.com/hengadev/cluo_api/internal/ports"
)

type Handler interface {
	RegisterRoutes(router *http.ServeMux)
}

type handler struct {
	svc    ports.MediaService
	authmw auth.AuthMiddleware
}

func New(svc ports.MediaService, authmw auth.AuthMiddleware) Handler {
	return &handler{
		svc:    svc,
		authmw: authmw,
	}
}
