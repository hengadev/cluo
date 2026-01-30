package ctxutil

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// GetUserIDFromContext retrieves the user ID from the context.
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return uuid.Nil, errors.New("user ID not found in context")
	}
	return userID, nil
}
