package client

import (
	"context"

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
	return errs.AsError()
}
