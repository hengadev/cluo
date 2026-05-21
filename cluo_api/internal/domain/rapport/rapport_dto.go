package rapport

import (
	"time"

	"github.com/google/uuid"
)

type CreateRapportRequest struct {
	CaseID  string `json:"caseId"`
	Content []byte `json:"content"` // raw TipTap JSON bytes
}

type UpdateRapportRequest struct {
	CaseID  uuid.UUID
	Content []byte `json:"content"`
}

type RapportResponse struct {
	ID        string    `json:"id"`
	CaseID    string    `json:"caseId"`
	Content   []byte    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (r *Rapport) ToResponse() *RapportResponse {
	return &RapportResponse{
		ID:        r.ID.String(),
		CaseID:    r.CaseID.String(),
		Content:   r.Content,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
