package case_type

import "time"

type CreateCaseTypeRequest struct {
	Name string `json:"name"`
}

type UpdateCaseTypeRequest struct {
	Name string `json:"name"`
}

type CaseTypeResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (ct *CaseType) ToResponse() *CaseTypeResponse {
	return &CaseTypeResponse{
		ID:        ct.ID.String(),
		Name:      ct.Name,
		CreatedAt: ct.CreatedAt,
		UpdatedAt: ct.UpdatedAt,
	}
}
