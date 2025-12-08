package caseDomain

import (
	"context"
	"fmt"
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

type DeleteCaseByIDRequest struct {
	ID uuid.UUID `json:"id"`
}

type UpdateCaseRequest struct {
	ID                uuid.UUID  `json:"id"`
	Title             *string    `json:"title"`
	Description       *string    `json:"description"`
	ClientID          *uuid.UUID `json:"clientId"`
	AssignedContactID *uuid.UUID `json:"assignedContactId"`
	Status            *string    `json:"status"`
}

func (r *UpdateCaseRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	// Validate ID
	if r.ID == uuid.Nil {
		errs.Set("id", "id is required")
	}

	// Validate Title (optional but if provided, check length)
	if r.Title != nil {
		if strings.TrimSpace(*r.Title) == "" {
			errs.Set("title", "title cannot be empty if provided")
		} else if len(*r.Title) > 200 {
			errs.Set("title", "title must be less than 200 characters")
		}
	}

	// Validate Description (optional but if provided, check length)
	if r.Description != nil && len(*r.Description) > 2000 {
		errs.Set("description", "description must be less than 2000 characters")
	}

	// Validate AssignedContactID (optional but if provided, must be valid UUID)
	if r.AssignedContactID != nil && *r.AssignedContactID == uuid.Nil {
		errs.Set("assignedContactId", "assignedContactId cannot be nil UUID if provided")
	}

	// Validate Status (optional but if provided, check valid values)
	if r.Status != nil {
		status := strings.TrimSpace(*r.Status)
		if status == "" {
			errs.Set("status", "status cannot be empty if provided")
		} else {
			caseStatus := CaseStatus(strings.ToLower(status))
			if !caseStatus.IsValid() {
				errs.Set("status", "status must be one of: draft, in_progress, ready, released")
			}
		}
	}

	return errs.AsError()
}

// ListCasesRequest represents a request to list cases with filtering and pagination
type ListCasesRequest struct {
	ClientID          *string `json:"clientId,omitempty"`
	Status            *string `json:"status,omitempty"`
	AssignedContactID *string `json:"assignedContactId,omitempty"`
	DateCreatedFrom   *string `json:"dateCreatedFrom,omitempty"`
	DateCreatedTo     *string `json:"dateCreatedTo,omitempty"`
	DateUpdatedFrom   *string `json:"dateUpdatedFrom,omitempty"`
	DateUpdatedTo     *string `json:"dateUpdatedTo,omitempty"`
	Search            *string `json:"search,omitempty"`
	Page              int     `json:"page" validate:"min=1"`
	PageSize          int     `json:"pageSize" validate:"min=1,max=100"`
}

func (r *ListCasesRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	// Validate Page
	if r.Page < 1 {
		errs.Set("page", "page must be at least 1")
	}

	// Validate PageSize
	if r.PageSize < 1 {
		errs.Set("pageSize", "pageSize must be at least 1")
	} else if r.PageSize > 100 {
		errs.Set("pageSize", "pageSize cannot exceed 100")
	}

	// Validate ClientID (optional)
	if r.ClientID != nil && strings.TrimSpace(*r.ClientID) != "" {
		if _, err := uuid.Parse(*r.ClientID); err != nil {
			errs.Set("clientId", "clientId must be a valid UUID")
		}
	}

	// Validate AssignedContactID (optional)
	if r.AssignedContactID != nil && strings.TrimSpace(*r.AssignedContactID) != "" {
		if _, err := uuid.Parse(*r.AssignedContactID); err != nil {
			errs.Set("assignedContactId", "assignedContactId must be a valid UUID")
		}
	}

	// Validate Status (optional)
	if r.Status != nil && strings.TrimSpace(*r.Status) != "" {
		status := CaseStatus(strings.ToLower(strings.TrimSpace(*r.Status)))
		if !status.IsValid() {
			errs.Set("status", "status must be one of: draft, in_progress, ready, released")
		}
	}

	// Validate date formats (optional)
	if r.DateCreatedFrom != nil && strings.TrimSpace(*r.DateCreatedFrom) != "" {
		if _, err := time.Parse(time.RFC3339, *r.DateCreatedFrom); err != nil {
			errs.Set("dateCreatedFrom", "dateCreatedFrom must be a valid RFC3339 datetime")
		}
	}

	if r.DateCreatedTo != nil && strings.TrimSpace(*r.DateCreatedTo) != "" {
		if _, err := time.Parse(time.RFC3339, *r.DateCreatedTo); err != nil {
			errs.Set("dateCreatedTo", "dateCreatedTo must be a valid RFC3339 datetime")
		}
	}

	if r.DateUpdatedFrom != nil && strings.TrimSpace(*r.DateUpdatedFrom) != "" {
		if _, err := time.Parse(time.RFC3339, *r.DateUpdatedFrom); err != nil {
			errs.Set("dateUpdatedFrom", "dateUpdatedFrom must be a valid RFC3339 datetime")
		}
	}

	if r.DateUpdatedTo != nil && strings.TrimSpace(*r.DateUpdatedTo) != "" {
		if _, err := time.Parse(time.RFC3339, *r.DateUpdatedTo); err != nil {
			errs.Set("dateUpdatedTo", "dateUpdatedTo must be a valid RFC3339 datetime")
		}
	}

	return errs.AsError()
}

// ToCaseFilter converts ListCasesRequest to CaseFilter for repository layer
func (r *ListCasesRequest) ToCaseFilter() (CaseFilter, error) {
	filter := CaseFilter{}

	// Parse ClientID
	if r.ClientID != nil && strings.TrimSpace(*r.ClientID) != "" {
		clientID, err := uuid.Parse(*r.ClientID)
		if err != nil {
			return filter, fmt.Errorf("invalid clientId: %w", err)
		}
		filter.ClientID = &clientID
	}

	// Parse Status
	if r.Status != nil && strings.TrimSpace(*r.Status) != "" {
		status := CaseStatus(strings.ToLower(strings.TrimSpace(*r.Status)))
		if !status.IsValid() {
			return filter, fmt.Errorf("invalid status: %s", *r.Status)
		}
		filter.Status = &status
	}

	// Parse AssignedContactID
	if r.AssignedContactID != nil && strings.TrimSpace(*r.AssignedContactID) != "" {
		contactID, err := uuid.Parse(*r.AssignedContactID)
		if err != nil {
			return filter, fmt.Errorf("invalid assignedContactId: %w", err)
		}
		filter.AssignedContactID = &contactID
	}

	// Parse date fields
	if r.DateCreatedFrom != nil && strings.TrimSpace(*r.DateCreatedFrom) != "" {
		dateFrom, err := time.Parse(time.RFC3339, *r.DateCreatedFrom)
		if err != nil {
			return filter, fmt.Errorf("invalid dateCreatedFrom: %w", err)
		}
		filter.DateCreatedFrom = &dateFrom
	}

	if r.DateCreatedTo != nil && strings.TrimSpace(*r.DateCreatedTo) != "" {
		dateTo, err := time.Parse(time.RFC3339, *r.DateCreatedTo)
		if err != nil {
			return filter, fmt.Errorf("invalid dateCreatedTo: %w", err)
		}
		filter.DateCreatedTo = &dateTo
	}

	if r.DateUpdatedFrom != nil && strings.TrimSpace(*r.DateUpdatedFrom) != "" {
		dateFrom, err := time.Parse(time.RFC3339, *r.DateUpdatedFrom)
		if err != nil {
			return filter, fmt.Errorf("invalid dateUpdatedFrom: %w", err)
		}
		filter.DateUpdatedFrom = &dateFrom
	}

	if r.DateUpdatedTo != nil && strings.TrimSpace(*r.DateUpdatedTo) != "" {
		dateTo, err := time.Parse(time.RFC3339, *r.DateUpdatedTo)
		if err != nil {
			return filter, fmt.Errorf("invalid dateUpdatedTo: %w", err)
		}
		filter.DateUpdatedTo = &dateTo
	}

	// Search field (string, no parsing needed)
	if r.Search != nil && strings.TrimSpace(*r.Search) != "" {
		search := strings.TrimSpace(*r.Search)
		filter.Search = &search
	}

	return filter, nil
}

// ListByClientRequest represents a request to list cases for a specific client
type ListByClientRequest struct {
	ClientID string `json:"clientId"`
	Page     int    `json:"page" validate:"min=1"`
	PageSize int    `json:"pageSize" validate:"min=1,max=100"`
}

func (r *ListByClientRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	// Validate ClientID
	if strings.TrimSpace(r.ClientID) == "" {
		errs.Set("clientId", "clientId is required")
	} else if _, err := uuid.Parse(r.ClientID); err != nil {
		errs.Set("clientId", "clientId must be a valid UUID")
	}

	// Validate Page
	if r.Page < 1 {
		errs.Set("page", "page must be at least 1")
	}

	// Validate PageSize
	if r.PageSize < 1 {
		errs.Set("pageSize", "pageSize must be at least 1")
	} else if r.PageSize > 100 {
		errs.Set("pageSize", "pageSize cannot exceed 100")
	}

	return errs.AsError()
}

// ListCasesResponse represents the response for list operations
type ListCasesResponse struct {
	Cases      []*CaseResponse `json:"cases"`
	Pagination PaginationInfo `json:"pagination"`
}

// PaginationInfo represents pagination metadata
type PaginationInfo struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

// NewPaginationInfo creates pagination metadata
func NewPaginationInfo(page, pageSize, totalItems int) PaginationInfo {
	totalPages := (totalItems + pageSize - 1) / pageSize
	if totalPages < 1 {
		totalPages = 1
	}

	return PaginationInfo{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}
