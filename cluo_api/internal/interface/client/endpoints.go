package clientHandler

const (
	// Base path
	ClientBasePath  = "/client"
	ContactBasePath = "/contact"

	// === Client Endpoints ===

	// Client CRUD operations
	CreateClientEndpoint  = ClientBasePath
	GetAllClientsEndpoint = ClientBasePath
	GetClientByIDEndpoint = ClientBasePath + "/{id}"
	UpdateClientEndpoint  = ClientBasePath + "/{id}"
	DeleteClientEndpoint  = ClientBasePath + "/{id}"

	// === Contact Endpoints ===

	// Contact CRUD operations
	CreateContactEndpoint            = ClientBasePath + "/{id}" + ContactBasePath
	GetAllContactsByClientIDEndpoint = ClientBasePath + "/{id}" + ContactBasePath
	GetContactIDsForClientEndpoint   = ClientBasePath + "/{id}/contact-ids"
	GetContactByIDEndpoint           = ContactBasePath + "/{id}"
	UpdateContactEndpoint            = ContactBasePath + "/{id}"
	DeleteContactEndpoint            = ContactBasePath + "/{id}"
)
