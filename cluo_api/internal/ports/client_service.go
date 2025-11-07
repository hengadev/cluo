package ports

import (
	"context"

	"github.com/hengadev/cluo_api/internal/domain/client"
)

type ClientService interface {
	CreateContact(ctx context.Context, r *client.CreateContactRequest) error
	DeleteContact(ctx context.Context, r *client.DeleteContactRequest) error
	UpdateContact(ctx context.Context, r *client.UpdateContactRequest) error
	GetContactByID(ctx context.Context, r *client.GetContactByIDRequest) (*client.ContactResponse, error)
	GetAllContactsByClientID(ctx context.Context, clientIDStr string) ([]*client.ContactResponse, error)
}
