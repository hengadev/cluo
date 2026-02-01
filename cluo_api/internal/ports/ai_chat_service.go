package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/ai"
)

// ChatService defines the interface for chat operations
type ChatService interface {
	// SendMessage sends a user message and returns the assistant's response
	SendMessage(ctx context.Context, req *SendMessageRequest) (*SendMessageResponse, error)

	// GetConversation retrieves a conversation with all messages
	GetConversation(ctx context.Context, conversationID uuid.UUID) (*ai.ChatConversation, []*ai.ChatMessage, error)

	// ListConversations retrieves all conversations for a case
	ListConversations(ctx context.Context, caseID uuid.UUID) ([]*ai.ChatConversation, error)

	// DeleteConversation deletes a conversation and all its messages
	DeleteConversation(ctx context.Context, conversationID uuid.UUID) error

	// StreamMessage sends a message and streams the response via the provided callback
	StreamMessage(ctx context.Context, req *SendMessageRequest, streamCallback StreamCallback) error
}

// SendMessageRequest is the request to send a message
type SendMessageRequest struct {
	ConversationID *uuid.UUID // nil for new conversation
	CaseID         uuid.UUID
	UserID         uuid.UUID
	Message        string
	CaseContext    *ai.ChatContext
}

// SendMessageResponse is the response from sending a message
type SendMessageResponse struct {
	ConversationID      uuid.UUID
	UserMessageID       uuid.UUID
	AssistantMessage    *ai.ChatMessage
	IsNewConversation   bool
}

// StreamCallback is called for each token during streaming
type StreamCallback func(token string, done bool, err error)

// ChatRepository defines the storage interface for chat (works with encrypted types)
type ChatRepository interface {
	// CreateConversation creates a new conversation (encryption handled by service layer)
	CreateConversation(ctx context.Context, conv *ai.ChatConversationEncx) error

	// GetConversation retrieves a conversation by ID (returns encrypted, decrypted by service layer)
	GetConversation(ctx context.Context, id uuid.UUID) (*ai.ChatConversationEncx, error)

	// UpdateConversation updates a conversation (encryption handled by service layer)
	UpdateConversation(ctx context.Context, conv *ai.ChatConversationEncx) error

	// DeleteConversation deletes a conversation
	DeleteConversation(ctx context.Context, id uuid.UUID) error

	// ListConversationsByCase retrieves conversations for a case (returns encrypted, decrypted by service layer)
	ListConversationsByCase(ctx context.Context, caseID uuid.UUID) ([]*ai.ChatConversationEncx, error)

	// CreateMessage creates a new message (encryption handled by service layer)
	CreateMessage(ctx context.Context, msg *ai.ChatMessageEncx) error

	// GetMessagesByConversation retrieves messages for a conversation (returns encrypted, decrypted by service layer)
	GetMessagesByConversation(ctx context.Context, conversationID uuid.UUID) ([]*ai.ChatMessageEncx, error)

	// GetRecentMessagesByConversation retrieves recent N messages (returns encrypted, decrypted by service layer)
	GetRecentMessagesByConversation(ctx context.Context, conversationID uuid.UUID, limit int) ([]*ai.ChatMessageEncx, error)
}
