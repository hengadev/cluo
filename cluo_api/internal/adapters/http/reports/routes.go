package reportHandler

import (
	"net/http"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	// router.HandleFunc("GET /{caseID}", h.GetReportByID)
	// router.HandleFunc("GET /case/{caseID}", h.GetReportByCaseID)
	// router.HandleFunc("POST /case/{caseID}", h.CreateReportByCaseID)
}
