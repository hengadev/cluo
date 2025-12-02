package client

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hengadev/cluo_api/internal/common/validation"

	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

type CreateContactRequest struct {
	ClientID  uuid.UUID `json:"clientID"`
	Lastname  string    `json:"lastname"`
	Firstname string    `json:"firstname"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Position  string    `json:"position"`
}

func (r *CreateContactRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	// UUID validation for ClientID
	if err := uuid.Validate(r.ClientID.String()); err != nil {
		errs.Set("clientID", err)
	}

	// Required string field validations
	if strings.TrimSpace(r.Lastname) == "" {
		errs.Set("lastname", fmt.Errorf("lastname is required"))
	}

	if strings.TrimSpace(r.Firstname) == "" {
		errs.Set("firstname", fmt.Errorf("firstname is required"))
	}

	// Use validation utilities for complex fields
	if err := validation.ValidateEmail(r.Email); err != nil {
		errs.Set("email", err)
	}

	if err := validation.ValidatePhone(r.Phone); err != nil {
		errs.Set("phone", err)
	}

	return errs.AsError()
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

type DeleteContactRequest struct {
	ContactID uuid.UUID `json:"contactID"`
}

type UpdateContactRequest struct {
	ID        uuid.UUID `json:"contactID"`
	ClientID  string    `json:"clientID" `
	Lastname  *string   `json:"lastname" `
	Firstname *string   `json:"firstname"`
	Email     *string   `json:"email" `
	Phone     *string   `json:"phone" `
	Position  *string   `json:"position"`
}

func (r *UpdateContactRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	// UUID validation for ID
	if err := uuid.Validate(r.ID.String()); err != nil {
		errs.Set("contact ID", err)
	}

	// UUID validation for ClientID
	if err := uuid.Validate(r.ClientID); err != nil {
		errs.Set("client ID", err)
	}

	// Validate email if provided
	if r.Email != nil {
		if err := validation.ValidateEmail(*r.Email); err != nil {
			errs.Set("email", err)
		}
	}

	// Validate phone if provided
	if r.Phone != nil {
		if err := validation.ValidatePhone(*r.Phone); err != nil {
			errs.Set("phone", err)
		}
	}

	return errs.AsError()
}

type GetContactByIDRequest struct {
	ContactID uuid.UUID `json:"contactID"`
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
