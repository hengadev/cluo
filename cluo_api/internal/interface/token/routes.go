package tokenHandler

import (
	"net/http"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /token/{token}", h.ValidateToken)
	router.HandleFunc("GET /token/{token}/report", h.GetReportByToken)
	router.HandleFunc("GET /token/{token}/media", h.GetAllMediaByToken)
}
