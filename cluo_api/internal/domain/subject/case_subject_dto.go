package subject

import (
	"time"

	"github.com/google/uuid"
)

type CreateCaseSubjectRequest struct {
	Lastname   string `json:"lastname"`
	Firstname  string `json:"firstname"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	City       string `json:"city"`
	PostalCode string `json:"postalCode"`
	Address1   string `json:"address1"`
	Address2   string `json:"address2"`
	Occupation string `json:"occupation"`
	Notes      string `json:"notes"`
}

type UpdateCaseSubjectRequest struct {
	ID         uuid.UUID
	Lastname   *string `json:"lastname"`
	Firstname  *string `json:"firstname"`
	Email      *string `json:"email"`
	Phone      *string `json:"phone"`
	City       *string `json:"city"`
	PostalCode *string `json:"postalCode"`
	Address1   *string `json:"address1"`
	Address2   *string `json:"address2"`
	Occupation *string `json:"occupation"`
	Notes      *string `json:"notes"`
}

type CaseSubjectResponse struct {
	ID         string    `json:"id"`
	Lastname   string    `json:"lastname"`
	Firstname  string    `json:"firstname"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	City       string    `json:"city"`
	PostalCode string    `json:"postalCode"`
	Address1   string    `json:"address1"`
	Address2   string    `json:"address2"`
	Occupation string    `json:"occupation"`
	Notes      string    `json:"notes"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (s *Subject) ToResponse() *CaseSubjectResponse {
	return &CaseSubjectResponse{
		ID:         s.ID.String(),
		Lastname:   s.Lastname,
		Firstname:  s.Firstname,
		Email:      s.Email,
		Phone:      s.Phone,
		City:       s.City,
		PostalCode: s.PostalCode,
		Address1:   s.Address1,
		Address2:   s.Address2,
		Occupation: s.Occupation,
		Notes:      s.Notes,
		CreatedAt:  s.CreatedAt,
	}
}
