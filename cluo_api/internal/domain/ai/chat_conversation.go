package ai

import (
	"time"

	"github.com/google/uuid"
)

// ChatConversation represents a chat session associated with a case
type ChatConversation struct {
	ID           uuid.UUID `db:"id"`
	CaseID       uuid.UUID `db:"case_id"`
	Title        string    `encx:"encrypt" db:"title_encrypted"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	CreatedBy    uuid.UUID `db:"created_by"`
	MessageCount int       `db:"message_count"`
	TotalTokens  int       `db:"total_tokens"`
}

// NewChatConversation creates a new ChatConversation entity
func NewChatConversation(caseID uuid.UUID, title string, createdBy uuid.UUID) *ChatConversation {
	now := time.Now()
	return &ChatConversation{
		ID:           uuid.New(),
		CaseID:       caseID,
		Title:        title,
		CreatedAt:    now,
		UpdatedAt:    now,
		CreatedBy:    createdBy,
		MessageCount: 0,
		TotalTokens:  0,
	}
}
