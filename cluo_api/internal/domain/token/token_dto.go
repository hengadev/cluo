package token

import "time"

type TokenResponse struct {
	ID        string     `json:"id"`
	CaseID    string     `json:"caseId"`
	ExpiresAt time.Time  `json:"expiresAt"`
	RevokedAt *time.Time `json:"revokedAt,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
}

type CreateTokenResponse struct {
	ID        string    `json:"id"`
	CaseID    string    `json:"caseId"`
	RawToken  string    `json:"rawToken"` // returned once, never stored
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}

func (t *Token) ToResponse() *TokenResponse {
	return &TokenResponse{
		ID:        t.ID.String(),
		CaseID:    t.CaseID.String(),
		ExpiresAt: t.ExpiresAt,
		RevokedAt: t.RevokedAt,
		CreatedAt: t.CreatedAt,
	}
}
