package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/client"
)

type ClientRepository interface {
	//contact
	CreateContact(ctx context.Context, contactEncx *client.ContactEncx) error
	ExistsClient(ctx context.Context, clientID uuid.UUID) (bool, error)
	DeleteContact(ctx context.Context, contactID uuid.UUID) error
	UpdateContact(ctx context.Context, request *client.ContactEncx) error
	GetContactByID(ctx context.Context, contactID uuid.UUID) (*client.ContactEncx, error)
	GetAllContactsByClientID(ctx context.Context, clientID uuid.UUID) ([]*client.ContactEncx, error)
	GetContactIDsForClient(ctx context.Context, clientID uuid.UUID) ([]uuid.UUID, error)
	//client
	CreateClient(ctx context.Context, clientEncx *client.ClientEncx) error
	DeleteClient(ctx context.Context, clientID uuid.UUID) error
	UpdateClient(ctx context.Context, request *client.ClientEncx) error
	GetClientByID(ctx context.Context, clientID uuid.UUID) (*client.ClientEncx, error)
	GetAllClients(ctx context.Context) ([]*client.ClientEncx, error)
}
