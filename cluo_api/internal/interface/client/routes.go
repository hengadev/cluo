package clientHandler

import (
	"net/http"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /client/{id}/contact", h.CreateContact)
	router.HandleFunc("GET /contact/{id}", h.GetContactByID)
	router.HandleFunc("DELETE /contact/{id}", h.DeleteContact)
	router.HandleFunc("PATCH /contact/{id}", h.UpdateContact)
	router.HandleFunc("GET /client/{id}/contact", h.GetAllContactsByClientID)
}
