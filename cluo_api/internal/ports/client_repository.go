package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/client"
)

type ClientRepository interface {
	CreateContact(ctx context.Context, contactEncx *client.ContactEncx) error
	ExistsByClientID(ctx context.Context, clientIDHashBasic string) (bool, error)
	DeleteContact(ctx context.Context, contactID uuid.UUID) error
	UpdateContact(ctx context.Context, request *client.ContactEncx) error
	GetContactByID(ctx context.Context, contactID uuid.UUID) (*client.ContactEncx, error)
	GetAllContactsByClientID(ctx context.Context, clientIDHashBasic string) ([]*client.ContactEncx, error)
}
