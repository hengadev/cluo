package client

import (
	"time"

	"github.com/google/uuid"
)

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
