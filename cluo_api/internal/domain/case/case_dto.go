package caseDomain

import (
	"time"
)

func (c *Case) ToResponse() *CaseResponse {
	return &CaseResponse{
		ID:                c.ID.String(),
		Title:             c.Title,
		Description:       c.Description,
		ClientID:          c.ClientID,
		AssignedContactID: c.AssignedContactID,
		Status:            c.Status,
		CreatedAt:         c.CreatedAt,
		UpdatedAt:         c.UpdatedAt,
	}

}

type CaseResponse struct {
	ID                string
	Title             string `encx:"encrypt"`
	Description       string `encx:"encrypt"`
	ClientID          string
	AssignedContactID *string
	Status            CaseStatus `encx:"encrypt"`
	CreatedAt         time.Time
	UpdatedAt         time.Time `encx:"encrypt"`
}
