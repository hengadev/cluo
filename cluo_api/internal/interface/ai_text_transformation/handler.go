package aiTextTransformationHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/middleware/auth"
	"github.com/hengadev/cluo_api/internal/ports"
)

type Handler interface {
	TransformText(w http.ResponseWriter, r *http.Request)
	HealthCheck(w http.ResponseWriter, r *http.Request)
	RegisterRoutes(router *http.ServeMux)
}

type handler struct {
	svc    ports.TextTransformationService
	authmw auth.AuthMiddleware
}

func New(svc ports.TextTransformationService, authmw auth.AuthMiddleware) Handler {
	return &handler{svc: svc, authmw: authmw}
}
