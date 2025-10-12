package tokenHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/ports"
	// "github.com/hengadev/cluo_api/core/middleware/auth"
)

type Handler interface {
	RegisterRoutes(router *http.ServeMux)
}

type handler struct {
	svc ports.TokenService
	// authmw auth.AuthMiddleware
}

// func New(svc ports.TokenService, authmw auth.AuthMiddleware) Handler {
func New(svc ports.TokenService) Handler {
	// return &handler{svc: svc, authmw: authmw}
	return &handler{svc: svc}
}
