package ai

import (
	"time"
)

// ChatConversationDTO is the data transfer object for conversations
type ChatConversationDTO struct {
	ID           string    `json:"id"`
	CaseID       string    `json:"caseId"`
	Title        string    `json:"title"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	CreatedBy    string    `json:"createdBy"`
	MessageCount int       `json:"messageCount"`
	TotalTokens  int       `json:"totalTokens"`
}

// ChatMessageDTO is the data transfer object for messages
type ChatMessageDTO struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversationId"`
	Role           string    `json:"role"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"createdAt"`
	TokenCount     int       `json:"tokenCount,omitempty"`
}

// ToDTO converts a ChatConversation to DTO
func (c *ChatConversation) ToDTO() ChatConversationDTO {
	return ChatConversationDTO{
		ID:           c.ID.String(),
		CaseID:       c.CaseID.String(),
		Title:        c.Title,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
		CreatedBy:    c.CreatedBy.String(),
		MessageCount: c.MessageCount,
		TotalTokens:  c.TotalTokens,
	}
}

// ToDTO converts a ChatMessage to DTO
func (m *ChatMessage) ToDTO() ChatMessageDTO {
	return ChatMessageDTO{
		ID:             m.ID.String(),
		ConversationID: m.ConversationID.String(),
		Role:           string(m.Role),
		Content:        m.Content,
		CreatedAt:      m.CreatedAt,
		TokenCount:     m.TokenCount,
	}
}

// ToDTOs converts a slice of ChatConversation to DTOs
func ConversationsToDTOs(convs []*ChatConversation) []ChatConversationDTO {
	dtos := make([]ChatConversationDTO, len(convs))
	for i, conv := range convs {
		dtos[i] = conv.ToDTO()
	}
	return dtos
}

// ToDTOs converts a slice of ChatMessage to DTOs
func MessagesToDTOs(msgs []*ChatMessage) []ChatMessageDTO {
	dtos := make([]ChatMessageDTO, len(msgs))
	for i, msg := range msgs {
		dtos[i] = msg.ToDTO()
	}
	return dtos
}
