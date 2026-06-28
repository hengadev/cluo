package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/user"
)

// UserRepository defines the interface for user data persistence operations
type UserRepository interface {
	// ExistsByEmailHash checks if a user exists by their email hash
	ExistsByEmailHash(ctx context.Context, emailHash string) (bool, error)

	// GetUserByEmailHash retrieves a user by their email hash
	GetUserByEmailHash(ctx context.Context, emailHash string) (*user.UserEncx, error)

	// GetUserByID retrieves a user by their ID
	GetUserByID(ctx context.Context, userID uuid.UUID) (*user.UserEncx, error)

	// CreateUser creates a new user in the database
	CreateUser(ctx context.Context, user *user.UserEncx) error

	// UpdateUserName updates the encrypted name for the given user
	UpdateUserName(ctx context.Context, userID uuid.UUID, nameEncrypted []byte) error
}
