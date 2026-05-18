package investigationHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/middleware/auth"
	"github.com/hengadev/cluo_api/internal/ports"
)

type Handler interface {
	CreateCase(w http.ResponseWriter, r *http.Request)
	GetCaseByID(w http.ResponseWriter, r *http.Request)
	UpdateCase(w http.ResponseWriter, r *http.Request)
	DeleteCase(w http.ResponseWriter, r *http.Request)
	ListCases(w http.ResponseWriter, r *http.Request)
	ListCasesByClient(w http.ResponseWriter, r *http.Request)
	RegisterRoutes(router *http.ServeMux)
}

type handler struct {
	svc    ports.CaseService
	authmw auth.AuthMiddleware
}

func New(svc ports.CaseService, authmw auth.AuthMiddleware) Handler {
	return &handler{svc: svc, authmw: authmw}
}
