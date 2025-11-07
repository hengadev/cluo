package authHandler

import (
	"net/http"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /auth/login", h.Login)
	router.HandleFunc("POST /auth/register", h.Register)
}
