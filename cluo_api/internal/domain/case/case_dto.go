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
		Status:            string(c.Status),
		CreatedAt:         c.CreatedAt,
		UpdatedAt:         c.UpdatedAt,
	}

}

type CaseResponse struct {
	ID                string    `json:"id"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	ClientID          string    `json:"clientId"`
	AssignedContactID *string   `json:"assignedContactID"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}
}
