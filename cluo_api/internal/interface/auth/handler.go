package authHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/ports"
	// "github.com/hengadev/cluo_api/core/middleware/auth"
)

type Handler interface {
	RegisterRoutes(router *http.ServeMux)
}

type handler struct {
	svc ports.AuthService
	// authmw auth.AuthMiddleware
}

// func New(svc ports.AuthService, authmw auth.AuthMiddleware) Handler {
func New(svc ports.AuthService) Handler {
	// return &handler{svc: svc, authmw: authmw}
	return &handler{svc: svc}
}
