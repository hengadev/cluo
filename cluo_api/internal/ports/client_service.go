package ports

import (
	"context"

	"github.com/hengadev/cluo_api/internal/domain/client"

	"github.com/google/uuid"
)

type ClientService interface {
	// client
	CreateClient(ctx context.Context, r *client.CreateClientRequest) error
	DeleteClient(ctx context.Context, r *client.DeleteClientRequest) error
	UpdateClient(ctx context.Context, r *client.UpdateClientRequest) error
	// contact
	CreateContact(ctx context.Context, r *client.CreateContactRequest) error
	DeleteContact(ctx context.Context, r *client.DeleteContactRequest) error
	UpdateContact(ctx context.Context, r *client.UpdateContactRequest) error
	GetContactByID(ctx context.Context, r *client.GetContactByIDRequest) (*client.ContactResponse, error)
	GetAllContactsByClientID(ctx context.Context, clientID uuid.UUID) ([]*client.ContactResponse, error)
}
