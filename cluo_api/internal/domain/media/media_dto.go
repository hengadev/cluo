package domain

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/errsx"
)

// UploadMediaRequest represents the request to upload a new media file with file content
type UploadMediaRequest struct {
	CaseID      string    `json:"caseId"`
	File        io.Reader `json:"-"` // The file content
	FileName    string    `json:"fileName"`
	MimeType    string    `json:"mimeType"`
	FileSize    int64     `json:"fileSize"`
	Caption     *string   `json:"caption"`
	IsPublished *bool     `json:"isPublished"`
}

func (r *UploadMediaRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	// Validate CaseID
	if strings.TrimSpace(r.CaseID) == "" {
		errs.Set("caseId", "caseId is required")
	} else {
		if _, err := uuid.Parse(r.CaseID); err != nil {
			errs.Set("caseId", "caseId must be a valid UUID")
		}
	}

	// Validate File
	if r.File == nil {
		errs.Set("file", "file is required")
	}

	// Validate MimeType
	if strings.TrimSpace(r.MimeType) == "" {
		errs.Set("mimeType", "mimeType is required")
	} else if len(r.MimeType) > 100 {
		errs.Set("mimeType", "mimeType must be less than 100 characters")
	}

	// Validate FileName
	if strings.TrimSpace(r.FileName) == "" {
		errs.Set("fileName", "fileName is required")
	} else if len(r.FileName) > 255 {
		errs.Set("fileName", "fileName must be less than 255 characters")
	}

	// Validate FileSize
	if r.FileSize <= 0 {
		errs.Set("fileSize", "fileSize must be greater than 0")
	}

	// Validate Caption (optional)
	if r.Caption != nil && len(*r.Caption) > 500 {
		errs.Set("caption", "caption must be less than 500 characters")
	}

	return errs.AsError()
}

// CreateMediaRequest represents the request to create a new media file (internal use after upload)
type CreateMediaRequest struct {
	CaseID      string  `json:"caseId"`
	URL         string  `json:"url"`
	Type        string  `json:"type"`
	MimeType    string  `json:"mimeType"`
	FileName    string  `json:"fileName"`
	FileSize    int64   `json:"fileSize"`
	Caption     *string `json:"caption"`
	IsPublished *bool   `json:"isPublished"`
}

func (r *CreateMediaRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	// Validate CaseID
	if strings.TrimSpace(r.CaseID) == "" {
		errs.Set("caseId", "caseId is required")
	} else {
		if _, err := uuid.Parse(r.CaseID); err != nil {
			errs.Set("caseId", "caseId must be a valid UUID")
		}
	}

	// Validate URL
	if strings.TrimSpace(r.URL) == "" {
		errs.Set("url", "url is required")
	} else if len(r.URL) > 2048 {
		errs.Set("url", "url must be less than 2048 characters")
	}

	// Validate Type
	if strings.TrimSpace(r.Type) == "" {
		errs.Set("type", "type is required")
	} else {
		mediaType := MediaType(strings.ToLower(strings.TrimSpace(r.Type)))
		if !mediaType.IsValid() {
			errs.Set("type", "type must be one of: image, video, audio")
		}
	}

	// Validate MimeType
	if strings.TrimSpace(r.MimeType) == "" {
		errs.Set("mimeType", "mimeType is required")
	} else if len(r.MimeType) > 100 {
		errs.Set("mimeType", "mimeType must be less than 100 characters")
	}

	// Validate FileName
	if strings.TrimSpace(r.FileName) == "" {
		errs.Set("fileName", "fileName is required")
	} else if len(r.FileName) > 255 {
		errs.Set("fileName", "fileName must be less than 255 characters")
	}

	// Validate FileSize
	if r.FileSize <= 0 {
		errs.Set("fileSize", "fileSize must be greater than 0")
	}

	// Validate Caption (optional)
	if r.Caption != nil && len(*r.Caption) > 500 {
		errs.Set("caption", "caption must be less than 500 characters")
	}

	return errs.AsError()
}

// NewMedia creates a new MediaFile domain object from a CreateMediaRequest
func NewMedia(r *CreateMediaRequest) *MediaFile {
	caseID, _ := uuid.Parse(r.CaseID)
	mediaType := MediaType(strings.ToLower(strings.TrimSpace(r.Type)))

	caption := ""
	if r.Caption != nil {
		caption = *r.Caption
	}

	isPublished := false
	if r.IsPublished != nil {
		isPublished = *r.IsPublished
	}

	return &MediaFile{
		ID:          uuid.New(),
		CaseID:      caseID,
		URL:         r.URL,
		Type:        mediaType,
		MimeType:    r.MimeType,
		FileName:    r.FileName,
		FileSize:    r.FileSize,
		Caption:     caption,
		IsPublished: isPublished,
		CreatedAt:   time.Now(),
	}
}

// UpdateMediaRequest represents the request to update an existing media file
type UpdateMediaRequest struct {
	ID          uuid.UUID `json:"id"`
	Caption     *string   `json:"caption"`
	IsPublished *bool     `json:"isPublished"`
}

func (r *UpdateMediaRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	if r.ID == uuid.Nil {
		errs.Set("id", "id is required")
	}

	// Validate Caption if provided
	if r.Caption != nil && len(*r.Caption) > 500 {
		errs.Set("caption", "caption must be less than 500 characters")
	}

	return errs.AsError()
}

// GetMediaByIDRequest represents the request to get a media file by ID
type GetMediaByIDRequest struct {
	ID uuid.UUID `json:"id"`
}

// DeleteMediaRequest represents the request to delete a media file
type DeleteMediaRequest struct {
	ID uuid.UUID `json:"id"`
}

// ListMediaByCaseIDRequest represents the request to list media files for a case
type ListMediaByCaseIDRequest struct {
	CaseID   string  `json:"caseId"`
	Type     *string `json:"type,omitempty"`
	Page     int     `json:"page"`
	PageSize int     `json:"pageSize"`
}

func (r *ListMediaByCaseIDRequest) Valid(ctx context.Context) error {
	var errs errsx.Map

	// Validate CaseID
	if strings.TrimSpace(r.CaseID) == "" {
		errs.Set("caseId", "caseId is required")
	} else {
		if _, err := uuid.Parse(r.CaseID); err != nil {
			errs.Set("caseId", "caseId must be a valid UUID")
		}
	}

	// Validate Type if provided
	if r.Type != nil && strings.TrimSpace(*r.Type) != "" {
		mediaType := MediaType(strings.ToLower(strings.TrimSpace(*r.Type)))
		if !mediaType.IsValid() {
			errs.Set("type", "type must be one of: image, video, audio")
		}
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

// MediaResponse represents the response for a media file
type MediaResponse struct {
	ID          string    `json:"id"`
	CaseID      string    `json:"caseId"`
	URL         string    `json:"url"`
	Type        string    `json:"type"`
	MimeType    string    `json:"mimeType"`
	FileName    string    `json:"fileName"`
	FileSize    int64     `json:"fileSize"`
	Caption     string    `json:"caption"`
	IsPublished bool      `json:"isPublished"`
	CreatedAt   time.Time `json:"createdAt"`
}

// ToResponse converts a MediaFile to a MediaResponse
func (m *MediaFile) ToResponse() *MediaResponse {
	return &MediaResponse{
		ID:          m.ID.String(),
		CaseID:      m.CaseID.String(),
		URL:         m.URL,
		Type:        string(m.Type),
		MimeType:    m.MimeType,
		FileName:    m.FileName,
		FileSize:    m.FileSize,
		Caption:     m.Caption,
		IsPublished: m.IsPublished,
		CreatedAt:   m.CreatedAt,
	}
}

// ListMediaResponse represents the response for listing media files
type ListMediaResponse struct {
	Media      []*MediaResponse `json:"media"`
	Pagination PaginationInfo   `json:"pagination"`
}

// PaginationInfo contains pagination metadata
type PaginationInfo struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

// NewPaginationInfo creates a new PaginationInfo
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
