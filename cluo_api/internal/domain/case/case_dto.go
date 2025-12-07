package caseDomain

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

func (c *Case) ToResponse() *CaseResponse {
	var assignedContactIDStr *string
	if c.AssignedContactID != nil {
		contactIDStr := c.AssignedContactID.String()
		assignedContactIDStr = &contactIDStr
	}

	return &CaseResponse{
		ID:                c.ID.String(),
		Title:             c.Title,
		Description:       c.Description,
		ClientID:          c.ClientID.String(),
		AssignedContactID: assignedContactIDStr,
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

type CreateCaseRequest struct {
	Title             string  `json:"title"`
	Description       string  `json:"description"`
	ClientID          string  `json:"clientId"`
	AssignedContactID *string `json:"assignedContactID"`
	Status            string  `json:"status"`
}

func (r *CreateCaseRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	// Validate Title
	if strings.TrimSpace(r.Title) == "" {
		errs.Set("title", "title is required")
	} else if len(r.Title) > 200 {
		errs.Set("title", "title must be less than 200 characters")
	}

	// Validate Description (optional but if provided, check length)
	if len(r.Description) > 2000 {
		errs.Set("description", "description must be less than 2000 characters")
	}

	// Validate ClientID
	if strings.TrimSpace(r.ClientID) == "" {
		errs.Set("clientId", "clientId is required")
	} else {
		if _, err := uuid.Parse(r.ClientID); err != nil {
			errs.Set("clientId", "clientId must be a valid UUID")
		}
	}

	// Validate AssignedContactID (optional)
	if r.AssignedContactID != nil && strings.TrimSpace(*r.AssignedContactID) != "" {
		if _, err := uuid.Parse(*r.AssignedContactID); err != nil {
			errs.Set("assignedContactId", "assignedContactId must be a valid UUID")
		}
	}

	// Validate Status
	if strings.TrimSpace(r.Status) == "" {
		errs.Set("status", "status is required")
	} else {
		status := CaseStatus(strings.ToLower(strings.TrimSpace(r.Status)))
		if !status.IsValid() {
			errs.Set("status", "status must be one of: draft, in_progress, ready, released")
		}
	}

	return errs.AsError()
}

// NewCase creates a new Case domain object from a CreateCaseRequest
func NewCase(r *CreateCaseRequest) *Case {
	now := time.Now()

	// Parse status, defaulting to Draft if invalid
	status := CaseStatus(strings.ToLower(strings.TrimSpace(r.Status)))
	if !status.IsValid() {
		status = CaseStatusDraft
	}

	// Parse client ID
	clientID, err := uuid.Parse(r.ClientID)
	if err != nil {
		// This should be handled in validation, but for safety, use nil UUID
		clientID = uuid.Nil
	}

	// Parse assigned contact ID if provided
	var assignedContactID *uuid.UUID
	if r.AssignedContactID != nil {
		contactID, err := uuid.Parse(*r.AssignedContactID)
		if err == nil {
			assignedContactID = &contactID
		}
	}

	return &Case{
		ID:                uuid.New(),
		Title:             r.Title,
		Description:       r.Description,
		ClientID:          clientID,
		AssignedContactID: assignedContactID,
		Status:            status,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

type GetCaseByIDRequest struct {
	ID uuid.UUID `json:"id"`
}
