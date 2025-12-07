package caseHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/middleware/auth"
	"github.com/hengadev/cluo_api/internal/ports"
)

type Handler interface {
	CreateCase(w http.ResponseWriter, r *http.Request)
	RegisterRoutes(router *http.ServeMux)
}

type handler struct {
	svc    ports.CaseService
	authmw auth.AuthMiddleware
}

func New(svc ports.CaseService, authmw auth.AuthMiddleware) Handler {
	return &handler{svc: svc, authmw: authmw}
}
