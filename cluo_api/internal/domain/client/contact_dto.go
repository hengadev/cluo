package client

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

type GetContactByIDRequest struct {
	ContactID string `json:"contactID"`
}

func (r *GetContactByIDRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	if r.ContactID == "" {
		errs.Set("contactID", "contact ID is required")
		return errs.AsError()
	}

	if err := uuid.Validate(r.ContactID); err != nil {
		errs.Set("contactID", "invalid contact ID format")
	}

	if len(r.ContactID) > 36 {
		errs.Set("contactID", "contact ID too long")
	}

	return errs.AsError()
}

type GetContactByIDResponse struct {
	ID        string    `json:"id"`
	ClientID  string    `json:"clientID"`
	Lastname  string    `json:"lastname"`
	Firstname string    `json:"firstname"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Position  string    `json:"position"`
	CreatedAt time.Time `json:"createdAt"`
}
