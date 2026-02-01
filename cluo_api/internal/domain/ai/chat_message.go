package ai

import (
	"time"

	"github.com/google/uuid"
)

// ChatRole represents the role of a message sender
type ChatRole string

const (
	ChatRoleUser      ChatRole = "user"
	ChatRoleAssistant ChatRole = "assistant"
	ChatRoleSystem    ChatRole = "system"
)

// IsValid checks if the role is valid
func (r ChatRole) IsValid() bool {
	switch r {
	case ChatRoleUser, ChatRoleAssistant, ChatRoleSystem:
		return true
	}
	return false
}

// ChatMessage represents a single message in a conversation
type ChatMessage struct {
	ID             uuid.UUID  `db:"id"`
	ConversationID uuid.UUID  `db:"conversation_id"`
	Role           ChatRole   `db:"role"`
	Content        string     `encx:"encrypt" db:"content_encrypted"`
	CreatedAt      time.Time  `db:"created_at"`
	TokenCount     int        `db:"token_count,omitempty"`
}

// NewChatMessage creates a new ChatMessage entity
func NewChatMessage(conversationID uuid.UUID, role ChatRole, content string) *ChatMessage {
	return &ChatMessage{
		ID:             uuid.New(),
		ConversationID: conversationID,
		Role:           role,
		Content:        content,
		CreatedAt:      time.Now(),
		TokenCount:     estimateTokenCount(content),
	}
}
