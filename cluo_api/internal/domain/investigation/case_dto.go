package investigation

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

type ReleaseResponse struct {
	CaseID    string    `json:"caseId"`
	TokenID   string    `json:"tokenId"`
	RawToken  string    `json:"rawToken"`
	PortalURL string    `json:"portalUrl"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type PortalCaseResponse struct {
	*CaseResponse
	TokenExpiresAt time.Time `json:"tokenExpiresAt"`
}

func (c *Investigation) ToResponse() *CaseResponse {
	var assignedContactIDStr *string
	if c.AssignedContactID != nil {
		contactIDStr := c.AssignedContactID.String()
		assignedContactIDStr = &contactIDStr
	}

	var caseSubjectIDStr *string
	if c.CaseSubjectID != nil {
		subjectIDStr := c.CaseSubjectID.String()
		caseSubjectIDStr = &subjectIDStr
	}

	var caseTypeIDStr *string
	if c.CaseTypeID != nil {
		typeIDStr := c.CaseTypeID.String()
		caseTypeIDStr = &typeIDStr
	}

	// Helper function to convert non-empty strings to pointers
	stringToPtr := func(s string) *string {
		if s == "" {
			return nil
		}
		return &s
	}

	return &CaseResponse{
		ID:                c.ID.String(),
		Title:             c.Title,
		Description:       c.Description,
		ClientID:          c.ClientID.String(),
		AssignedContactID: assignedContactIDStr,
		CaseSubjectID:     caseSubjectIDStr,
		ExternalReference: c.ExternalReference,
		CaseTypeID:        caseTypeIDStr,
		Status:            string(c.Status),
		Placename:         stringToPtr(c.Placename),
		Address1:          stringToPtr(c.Address1),
		Address2:          stringToPtr(c.Address2),
		City:              stringToPtr(c.City),
		PostalCode:        stringToPtr(c.PostalCode),
		Country:           stringToPtr(c.Country),
		Latitude:          c.Latitude,
		Longitude:         c.Longitude,
		LocationType:      stringToPtr(c.LocationType),
		LocationNotes:     stringToPtr(c.LocationNotes),
		CreatedAt:         c.CreatedAt,
		UpdatedAt:         c.UpdatedAt,
	}

}

type CaseResponse struct {
	ID                string    `json:"id"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	ClientID          string    `json:"clientId"`
	AssignedContactID *string   `json:"assignedContactID,omitempty"`
	CaseSubjectID     *string   `json:"caseSubjectId,omitempty"`
	ExternalReference *string   `json:"externalReference,omitempty"`
	CaseTypeID        *string   `json:"caseTypeId,omitempty"`
	Status            string    `json:"status"`
	Placename         *string   `json:"placename,omitempty"`
	Address1          *string   `json:"address1,omitempty"`
	Address2          *string   `json:"address2,omitempty"`
	City              *string   `json:"city,omitempty"`
	PostalCode        *string   `json:"postalCode,omitempty"`
	Country           *string   `json:"country,omitempty"`
	Latitude          *string   `json:"latitude,omitempty"`
	Longitude         *string   `json:"longitude,omitempty"`
	LocationType      *string   `json:"locationType,omitempty"`
	LocationNotes     *string   `json:"locationNotes,omitempty"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type CreateCaseRequest struct {
	Title             string  `json:"title"`
	Description       string  `json:"description"`
	ClientID          string  `json:"clientId"`
	AssignedContactID *string `json:"assignedContactID,omitempty"`
	CaseSubjectID     *string `json:"caseSubjectId,omitempty"`
	ExternalReference *string `json:"externalReference,omitempty"`
	CaseTypeID        *string `json:"caseTypeId,omitempty"`
	Status            string  `json:"status"`
	Placename         *string `json:"placename,omitempty"`
	Address1          *string `json:"address1,omitempty"`
	Address2          *string `json:"address2,omitempty"`
	City              *string `json:"city,omitempty"`
	PostalCode        *string `json:"postalCode,omitempty"`
	Country           *string `json:"country,omitempty"`
	Latitude          *string `json:"latitude,omitempty"`
	Longitude         *string `json:"longitude,omitempty"`
	LocationType      *string `json:"locationType,omitempty"`
	LocationNotes     *string `json:"locationNotes,omitempty"`
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
		status := Status(strings.ToLower(strings.TrimSpace(r.Status)))
		if !status.IsValid() {
			errs.Set("status", "status must be one of: in_progress, ready, released")
		}
	}

	// Validate CaseSubjectID (optional)
	if r.CaseSubjectID != nil && strings.TrimSpace(*r.CaseSubjectID) != "" {
		if _, err := uuid.Parse(*r.CaseSubjectID); err != nil {
			errs.Set("caseSubjectId", "caseSubjectId must be a valid UUID")
		}
	}

	// Validate location fields (all optional)
	if r.Placename != nil && len(*r.Placename) > 200 {
		errs.Set("placename", "placename must be less than 200 characters")
	}

	if r.Address1 != nil && len(*r.Address1) > 200 {
		errs.Set("address1", "address1 must be less than 200 characters")
	}

	if r.Address2 != nil && len(*r.Address2) > 200 {
		errs.Set("address2", "address2 must be less than 200 characters")
	}

	if r.City != nil && len(*r.City) > 100 {
		errs.Set("city", "city must be less than 100 characters")
	}

	if r.PostalCode != nil && len(*r.PostalCode) > 20 {
		errs.Set("postalCode", "postalCode must be less than 20 characters")
	}

	if r.Country != nil && len(*r.Country) > 100 {
		errs.Set("country", "country must be less than 100 characters")
	}

	if r.LocationType != nil && len(*r.LocationType) > 100 {
		errs.Set("locationType", "locationType must be less than 100 characters")
	}

	if r.LocationNotes != nil && len(*r.LocationNotes) > 1000 {
		errs.Set("locationNotes", "locationNotes must be less than 1000 characters")
	}

	// Validate latitude and longitude format (optional, but must be valid decimal if provided)
	if r.Latitude != nil && strings.TrimSpace(*r.Latitude) != "" {
		// Simple validation: should be parseable as float and within valid range
		var lat float64
		if _, err := fmt.Sscanf(*r.Latitude, "%f", &lat); err != nil {
			errs.Set("latitude", "latitude must be a valid decimal number")
		} else if lat < -90.0 || lat > 90.0 {
			errs.Set("latitude", "latitude must be between -90 and 90")
		}
	}

	if r.Longitude != nil && strings.TrimSpace(*r.Longitude) != "" {
		// Simple validation: should be parseable as float and within valid range
		var lon float64
		if _, err := fmt.Sscanf(*r.Longitude, "%f", &lon); err != nil {
			errs.Set("longitude", "longitude must be a valid decimal number")
		} else if lon < -180.0 || lon > 180.0 {
			errs.Set("longitude", "longitude must be between -180 and 180")
		}
	}

	return errs.AsError()
}

// New creates a new Investigation domain object from a CreateCaseRequest
func New(r *CreateCaseRequest) *Investigation {
	now := time.Now()

	// Parse status, defaulting to InProgress if invalid
	status := Status(strings.ToLower(strings.TrimSpace(r.Status)))
	if !status.IsValid() {
		status = StatusInProgress
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

	// Parse case subject ID if provided
	var caseSubjectID *uuid.UUID
	if r.CaseSubjectID != nil {
		subjectID, err := uuid.Parse(*r.CaseSubjectID)
		if err == nil {
			caseSubjectID = &subjectID
		}
	}

	// Parse case type ID if provided
	var parsedCaseTypeID *uuid.UUID
	if r.CaseTypeID != nil && strings.TrimSpace(*r.CaseTypeID) != "" {
		typeID, err := uuid.Parse(*r.CaseTypeID)
		if err == nil {
			parsedCaseTypeID = &typeID
		}
	}

	// Helper to dereference string pointers, returning empty string if nil
	ptrToString := func(s *string) string {
		if s == nil {
			return ""
		}
		return *s
	}

	return &Investigation{
		ID:                uuid.New(),
		Title:             r.Title,
		Description:       r.Description,
		ClientID:          clientID,
		AssignedContactID: assignedContactID,
		CaseSubjectID:     caseSubjectID,
		ExternalReference: r.ExternalReference,
		CaseTypeID:        parsedCaseTypeID,
		Status:            status,
		Placename:         ptrToString(r.Placename),
		Address1:          ptrToString(r.Address1),
		Address2:          ptrToString(r.Address2),
		City:              ptrToString(r.City),
		PostalCode:        ptrToString(r.PostalCode),
		Country:           ptrToString(r.Country),
		Latitude:          r.Latitude,
		Longitude:         r.Longitude,
		LocationType:      ptrToString(r.LocationType),
		LocationNotes:     ptrToString(r.LocationNotes),
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
	CaseSubjectID     *uuid.UUID `json:"caseSubjectId"`
	ExternalReference *string    `json:"externalReference"`
	CaseTypeID        *uuid.UUID `json:"caseTypeId,omitempty"`
	Status            *string    `json:"status"`
	Placename         *string    `json:"placename"`
	Address1          *string    `json:"address1"`
	Address2          *string    `json:"address2"`
	City              *string    `json:"city"`
	PostalCode        *string    `json:"postalCode"`
	Country           *string    `json:"country"`
	Latitude          *string    `json:"latitude"`
	Longitude         *string    `json:"longitude"`
	LocationType      *string    `json:"locationType"`
	LocationNotes     *string    `json:"locationNotes"`
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
			caseStatus := Status(strings.ToLower(status))
			if !caseStatus.IsValid() {
				errs.Set("status", "status must be one of: in_progress, ready, released")
			}
		}
	}

	// Validate CaseSubjectID (optional but if provided, cannot be nil UUID)
	if r.CaseSubjectID != nil && *r.CaseSubjectID == uuid.Nil {
		errs.Set("caseSubjectId", "caseSubjectId cannot be nil UUID if provided")
	}

	// Validate location fields (all optional, check length if provided)
	if r.Placename != nil && len(*r.Placename) > 200 {
		errs.Set("placename", "placename must be less than 200 characters")
	}

	if r.Address1 != nil && len(*r.Address1) > 200 {
		errs.Set("address1", "address1 must be less than 200 characters")
	}

	if r.Address2 != nil && len(*r.Address2) > 200 {
		errs.Set("address2", "address2 must be less than 200 characters")
	}

	if r.City != nil && len(*r.City) > 100 {
		errs.Set("city", "city must be less than 100 characters")
	}

	if r.PostalCode != nil && len(*r.PostalCode) > 20 {
		errs.Set("postalCode", "postalCode must be less than 20 characters")
	}

	if r.Country != nil && len(*r.Country) > 100 {
		errs.Set("country", "country must be less than 100 characters")
	}

	if r.LocationType != nil && len(*r.LocationType) > 100 {
		errs.Set("locationType", "locationType must be less than 100 characters")
	}

	if r.LocationNotes != nil && len(*r.LocationNotes) > 1000 {
		errs.Set("locationNotes", "locationNotes must be less than 1000 characters")
	}

	// Validate latitude and longitude format (optional, but must be valid decimal if provided)
	if r.Latitude != nil && strings.TrimSpace(*r.Latitude) != "" {
		var lat float64
		if _, err := fmt.Sscanf(*r.Latitude, "%f", &lat); err != nil {
			errs.Set("latitude", "latitude must be a valid decimal number")
		} else if lat < -90.0 || lat > 90.0 {
			errs.Set("latitude", "latitude must be between -90 and 90")
		}
	}

	if r.Longitude != nil && strings.TrimSpace(*r.Longitude) != "" {
		var lon float64
		if _, err := fmt.Sscanf(*r.Longitude, "%f", &lon); err != nil {
			errs.Set("longitude", "longitude must be a valid decimal number")
		} else if lon < -180.0 || lon > 180.0 {
			errs.Set("longitude", "longitude must be between -180 and 180")
		}
	}

	return errs.AsError()
}

// ListCasesRequest represents a request to list cases with filtering and pagination
type ListCasesRequest struct {
	ClientID          *string `json:"clientId,omitempty"`
	Status            *string `json:"status,omitempty"`
	AssignedContactID *string `json:"assignedContactId,omitempty"`
	CaseSubjectID     *string `json:"caseSubjectId,omitempty"`
	CaseTypeID        *string `json:"caseTypeId,omitempty"`
	City              *string `json:"city,omitempty"`
	PostalCode        *string `json:"postalCode,omitempty"`
	Country           *string `json:"country,omitempty"`
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

	// Validate CaseSubjectID (optional)
	if r.CaseSubjectID != nil && strings.TrimSpace(*r.CaseSubjectID) != "" {
		if _, err := uuid.Parse(*r.CaseSubjectID); err != nil {
			errs.Set("caseSubjectId", "caseSubjectId must be a valid UUID")
		}
	}

	// Validate Status (optional)
	if r.Status != nil && strings.TrimSpace(*r.Status) != "" {
		status := Status(strings.ToLower(strings.TrimSpace(*r.Status)))
		if !status.IsValid() {
			errs.Set("status", "status must be one of: in_progress, ready, released")
		}
	}

	// Validate CaseTypeID (optional)
	if r.CaseTypeID != nil && strings.TrimSpace(*r.CaseTypeID) != "" {
		if _, err := uuid.Parse(*r.CaseTypeID); err != nil {
			errs.Set("caseTypeId", "caseTypeId must be a valid UUID")
		}
	}

	// Validate location filters (optional)
	if r.City != nil && len(*r.City) > 100 {
		errs.Set("city", "city must be less than 100 characters")
	}

	if r.PostalCode != nil && len(*r.PostalCode) > 20 {
		errs.Set("postalCode", "postalCode must be less than 20 characters")
	}

	if r.Country != nil && len(*r.Country) > 100 {
		errs.Set("country", "country must be less than 100 characters")
	}

	// Validate Search (optional)
	if r.Search != nil && len(*r.Search) > 1000 {
		errs.Set("search", "search must be less than 1000 characters")
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

// ToCaseFilter converts ListCasesRequest to Filter for repository layer
func (r *ListCasesRequest) ToFilter() (Filter, error) {
	filter := Filter{}

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
		status := Status(strings.ToLower(strings.TrimSpace(*r.Status)))
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

	// Parse CaseSubjectID
	if r.CaseSubjectID != nil && strings.TrimSpace(*r.CaseSubjectID) != "" {
		subjectID, err := uuid.Parse(*r.CaseSubjectID)
		if err != nil {
			return filter, fmt.Errorf("invalid caseSubjectId: %w", err)
		}
		filter.CaseSubjectID = &subjectID
	}

	// Parse CaseTypeID
	if r.CaseTypeID != nil && strings.TrimSpace(*r.CaseTypeID) != "" {
		caseTypeID, err := uuid.Parse(*r.CaseTypeID)
		if err != nil {
			return filter, fmt.Errorf("invalid caseTypeId: %w", err)
		}
		filter.CaseTypeID = &caseTypeID
	}

	// Parse location filters
	if r.City != nil && strings.TrimSpace(*r.City) != "" {
		city := strings.TrimSpace(*r.City)
		filter.City = &city
	}

	if r.PostalCode != nil && strings.TrimSpace(*r.PostalCode) != "" {
		postalCode := strings.TrimSpace(*r.PostalCode)
		filter.PostalCode = &postalCode
	}

	if r.Country != nil && strings.TrimSpace(*r.Country) != "" {
		country := strings.TrimSpace(*r.Country)
		filter.Country = &country
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
	Pagination PaginationInfo  `json:"pagination"`
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
