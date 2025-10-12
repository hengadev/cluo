package mediaHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/ports"
	// "github.com/hengadev/cluo_api/core/middleware/auth"
)

type Handler interface {
	RegisterRoutes(router *http.ServeMux)
}

type handler struct {
	svc ports.MediaService
	// authmw auth.MediaMiddleware
}

// func New(svc ports.MediaService, authmw auth.MediaMiddleware) Handler {
func New(svc ports.MediaService) Handler {
	// return &handler{svc: svc, authmw: authmw}
	return &handler{svc: svc}
}
