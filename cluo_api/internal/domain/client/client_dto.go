package client

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

type CreateClientRequest struct {
	ID   uuid.UUID
	Name string `json:"name"`
	Type string `json:"type"`
}

func (r *CreateClientRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	if r.Name == "" {
		errs.Set("name", "name is required")
	}

	if r.Type == "" {
		errs.Set("type", "type is required")
	}

	if len(r.Name) > 100 {
		errs.Set("name", "name must be less than 100 characters")
	}

	if len(r.Type) > 50 {
		errs.Set("type", "type must be less than 50 characters")
	}

	if r.ID == uuid.Nil {
		errs.Set("id", "valid ID is required")
	}

	return errs.AsError()
}

func NewClient(r *CreateClientRequest) *Client {
	return &Client{
		ID:         uuid.New(),
		Name:       r.Name,
		Type:       ClientType(r.Type),
		ContactIDs: []string{},
		CreatedAt:  time.Now(),
	}
}

type DeleteClientRequest struct {
	ID uuid.UUID `json:"id"`
}

type UpdateClientRequest struct {
	ID   uuid.UUID `json:"id"`
	Name *string   `json:"name" `
	Type *string   `json:"type"`
}

func (r *UpdateClientRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	// Validate name if provided
	if r.Name != nil {
		if *r.Name == "" {
			errs.Set("name", "client name is required")
		}

		if len(*r.Name) > 100 {
			errs.Set("name", "name must be less than 100 characters")
		}
	}

	// Validate type if provided
	if r.Type != nil {
		if *r.Type == "" {
			errs.Set("type", "client type is required")
		}

		if len(*r.Type) > 50 {
			errs.Set("type", "type must be less than 50 characters")
		}
	}

	return errs.AsError()
}

type GetClientByIDRequest struct {
	ID uuid.UUID `json:"id"`
}

type ClientResponse struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	ContactIDs []string `json:"contacts" encx:"encrypt"` // the list of contact IDs
}
