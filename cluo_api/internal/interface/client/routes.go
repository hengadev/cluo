package clientHandler

import (
	"net/http"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	// Client routes
	router.HandleFunc("POST /client", h.CreateClient)
	router.HandleFunc("DELETE /client/{id}", h.DeleteClient)
	router.HandleFunc("PATCH /client/{id}", h.UpdateClient)

	// Contact routes
	router.HandleFunc("POST /client/{id}/contact", h.CreateContact)
	router.HandleFunc("GET /contact/{id}", h.GetContactByID)
	router.HandleFunc("DELETE /contact/{id}", h.DeleteContact)
	router.HandleFunc("PATCH /contact/{id}", h.UpdateContact)
	router.HandleFunc("GET /client/{id}/contact", h.GetAllContactsByClientID)
}
