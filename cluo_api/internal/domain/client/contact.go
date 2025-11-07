package client

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

type Contact struct {
	ID        uuid.UUID `json:"-"`
	ClientID  uuid.UUID `json:"-" encx:"encrypt,hash_basic"`
	Lastname  string    `json:"lastname" encx:"encrypt"`
	Firstname string    `json:"firstname" encx:"encrypt"`
	Email     string    `json:"email" encx:"encrypt,hash_basic"`
	Phone     string    `json:"phone" encx:"encrypt"`
	Position  string    `json:"position" encx:"encrypt"`
	CreatedAt time.Time
}

func NewContact(r *CreateContactRequest) *Contact {
	return &Contact{
		ID:        uuid.New(),
		Lastname:  r.Lastname,
		Firstname: r.Firstname,
		Email:     r.Email,
		Phone:     r.Phone,
		Position:  r.Position,
		CreatedAt: time.Now(),
	}
}

type CreateContactRequest struct {
	// TODO: is this how the json should be formatted ?
	ClientID  string `json:"clientID"`
	Lastname  string `json:"lastname"`
	Firstname string `json:"firstname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Position  string `json:"position"`
}

func (r *CreateContactRequest) Valid(ctx context.Context) error {
	var errs errsx.Map
	if err := uuid.Validate(r.ClientID); err != nil {
		errs.Set("clientID", err)
	}
	return errs.AsError()
}

type DeleteContactRequest struct {
	ContactID string `json:"contactID"`
}

func (r *DeleteContactRequest) Valid(ctx context.Context) error {
	var errs errsx.Map
	if err := uuid.Validate(r.ContactID); err != nil {
		errs.Set("contact ID", err)
	}
	return errs.AsError()
}

type UpdateContactRequest struct {
	ID        string  `json:"contactID"`
	ClientID  string  `json:"clientID" encx:"encrypt,hash_basic"`
	Lastname  *string `json:"lastname" encx:"encrypt"`
	Firstname *string `json:"firstname" encx:"encrypt"`
	Email     *string `json:"email" encx:"encrypt,hash_basic"`
	Phone     *string `json:"phone" encx:"encrypt"`
	Position  *string `json:"position" encx:"encrypt"`
}

func (r *UpdateContactRequest) Valid(ctx context.Context) error {
	var errs errsx.Map
	// TODO: complete that validation with other rules that make sense for the different fields
	if err := uuid.Validate(r.ID); err != nil {
		errs.Set("contact ID", err)
	}
	if err := uuid.Validate(r.ClientID); err != nil {
		errs.Set("client ID", err)
	}
	return errs.AsError()
}

type ContactResponse struct {
	ID        string    `json:"id"`
	ClientID  string    `json:"clientID"`
	Lastname  string    `json:"lastname"`
	Firstname string    `json:"firstname"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Position  string    `json:"position"`
	CreatedAt time.Time `json:"createdAt"`
}
