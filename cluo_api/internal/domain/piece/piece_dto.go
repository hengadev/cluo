package piece

import (
	"time"

	"github.com/google/uuid"
)

// CreatePieceRequest holds the data needed to create a new piece record after upload.
type CreatePieceRequest struct {
	CaseID     uuid.UUID
	Filename   string
	StorageKey string
	MimeType   string
	SizeBytes  int64
	Notes      string
}

// PieceResponse is the public-facing representation of a piece.
type PieceResponse struct {
	ID         string    `json:"id"`
	CaseID     string    `json:"caseId"`
	Filename   string    `json:"filename"`
	StorageKey string    `json:"storageKey"`
	MimeType   string    `json:"mimeType"`
	SizeBytes  int64     `json:"sizeBytes"`
	Notes      string    `json:"notes,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// Pagination holds page/pageSize for list queries.
type Pagination struct {
	Page     int
	PageSize int
}

// Offset returns the SQL OFFSET for the given page.
func (p Pagination) Offset() int { return (p.Page - 1) * p.PageSize }

// ListPiecesResponse wraps a page of results with pagination metadata.
type ListPiecesResponse struct {
	Pieces     []*PieceResponse `json:"pieces"`
	Pagination PaginationInfo   `json:"pagination"`
}

// PaginationInfo contains pagination metadata returned to callers.
type PaginationInfo struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

// NewPaginationInfo builds a PaginationInfo from raw counts.
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

// ToResponse converts a Piece to its JSON-friendly representation.
func (p *Piece) ToResponse() *PieceResponse {
	return &PieceResponse{
		ID:         p.ID.String(),
		CaseID:     p.CaseID.String(),
		Filename:   p.Filename,
		StorageKey: p.StorageKey,
		MimeType:   p.MimeType,
		SizeBytes:  p.SizeBytes,
		Notes:      p.Notes,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}
}
