package clientHandler

import (
	"net/http"

	"github.com/hengadev/cluo_api/internal/common/contracts/identity"
	mw "github.com/hengadev/cluo_api/internal/common/middleware"
)

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	RequireClient := h.authmw.RequireMinimumRole(identity.Client)
	RequireAdministrator := h.authmw.RequireMinimumRole(identity.Administrator)

	// === Client CRUD Endpoints ===

	// Creates a new client
	router.HandleFunc("POST "+CreateClientEndpoint, RequireAdministrator(mw.EnableCORS(h.CreateClient)))

	// Gets all clients
	router.HandleFunc("GET "+GetAllClientsEndpoint, RequireAdministrator(mw.EnableCORS(h.GetAllClients)))

	// Gets a client by ID
	router.HandleFunc("GET "+GetClientByIDEndpoint, RequireAdministrator(mw.EnableCORS(h.GetClientByID)))

	// Updates a client
	router.HandleFunc("PATCH "+UpdateClientEndpoint, RequireAdministrator(mw.EnableCORS(h.UpdateClient)))

	// Deletes a client
	router.HandleFunc("DELETE "+DeleteClientEndpoint, RequireAdministrator(mw.EnableCORS(h.DeleteClient)))

	// === Contact CRUD Endpoints ===

	// Creates a contact for a client
	router.HandleFunc("POST "+CreateContactEndpoint, RequireClient(mw.EnableCORS(h.CreateContact)))

	// Gets all contacts for a specific client
	router.HandleFunc("GET "+GetAllContactsByClientIDEndpoint, RequireClient(mw.EnableCORS(h.GetAllContactsByClientID)))

	// Gets all contact IDs for a specific client
	router.HandleFunc("GET "+GetContactIDsForClientEndpoint, RequireAdministrator(mw.EnableCORS(h.GetContactIDsForClient)))

	// Gets a contact by ID
	router.HandleFunc("GET "+GetContactByIDEndpoint, RequireClient(mw.EnableCORS(h.GetContactByID)))

	// Updates a contact
	router.HandleFunc("PATCH "+UpdateContactEndpoint, RequireClient(mw.EnableCORS(h.UpdateContact)))

	// Deletes a contact
	router.HandleFunc("DELETE "+DeleteContactEndpoint, RequireClient(mw.EnableCORS(h.DeleteContact)))
}
