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
