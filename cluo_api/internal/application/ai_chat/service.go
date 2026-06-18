package aiChat

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/ai"
	"github.com/hengadev/cluo_api/internal/ports"
	"github.com/hengadev/encx"
)

const (
	maxConversationHistory = 20 // Maximum messages to include in context
	systemPromptTemplate   = `You are an AI assistant helping with an investigation case.

CASE CONTEXT:
- Case ID: %s
- Case Title: %s
- Case Type: %s
- Status: %s
- Client: %s
- Location: %s
- Subjects:
%s

IMPORTANT: You are discussing this specific case. Be helpful but concise. If you don't know something, say so. Keep responses focused on the case at hand.

CONVERSATION HISTORY:
%s`
)

type Service struct {
	chatRepo ports.ChatRepository
	llm      ports.LLMClient
	crypto   encx.CryptoService
	logger   *slog.Logger
}

func New(chatRepo ports.ChatRepository, llm ports.LLMClient, crypto encx.CryptoService, logger *slog.Logger) *Service {
	return &Service{
		chatRepo: chatRepo,
		llm:      llm,
		crypto:   crypto,
		logger:   logger.With("component", "ai_chat"),
	}
}

func (s *Service) SendMessage(ctx context.Context, req *ports.SendMessageRequest) (*ports.SendMessageResponse, error) {
	// Determine if new or existing conversation
	var conversation *ai.ChatConversation
	isNewConversation := req.ConversationID == nil

	if isNewConversation {
		// Create new conversation with title from first message
		title := s.generateConversationTitle(req.Message)
		conversation = ai.NewChatConversation(req.CaseID, title, req.UserID)

		// Encrypt and create conversation
		convEncx, err := ai.ProcessChatConversationEncx(ctx, s.crypto, conversation)
		if err != nil {
			return nil, fmt.Errorf("encrypt conversation: %w", err)
		}
		if err := s.chatRepo.CreateConversation(ctx, convEncx); err != nil {
			return nil, fmt.Errorf("create conversation: %w", err)
		}
	} else {
		// Get existing conversation
		convEncx, err := s.chatRepo.GetConversation(ctx, *req.ConversationID)
		if err != nil {
			return nil, fmt.Errorf("get conversation: %w", err)
		}
		conversation, err = ai.DecryptChatConversationEncx(ctx, s.crypto, convEncx)
		if err != nil {
			return nil, fmt.Errorf("decrypt conversation: %w", err)
		}
	}

	// Create and encrypt user message
	userMsg := ai.NewChatMessage(conversation.ID, ai.ChatRoleUser, req.Message)
	userMsgEncx, err := ai.ProcessChatMessageEncx(ctx, s.crypto, userMsg)
	if err != nil {
		return nil, fmt.Errorf("encrypt user message: %w", err)
	}
	if err := s.chatRepo.CreateMessage(ctx, userMsgEncx); err != nil {
		return nil, fmt.Errorf("create user message: %w", err)
	}

	// Build conversation history for context
	var history []*ai.ChatMessage
	historyEncx, err := s.chatRepo.GetRecentMessagesByConversation(ctx, conversation.ID, maxConversationHistory)
	if err != nil {
		s.logger.WarnContext(ctx, "Failed to get conversation history", "error", err)
		history = []*ai.ChatMessage{userMsg}
	} else {
		// Decrypt messages for context building
		history = make([]*ai.ChatMessage, len(historyEncx))
		for i, msgEncx := range historyEncx {
			history[i], err = ai.DecryptChatMessageEncx(ctx, s.crypto, msgEncx)
			if err != nil {
				s.logger.WarnContext(ctx, "Failed to decrypt message", "error", err)
				history[i] = ai.NewChatMessage(msgEncx.ConversationID, msgEncx.Role, "[decryption error]")
			}
		}
		history = append(history, userMsg)
	}

	// Build system prompt with case context
	systemPrompt := s.buildSystemPrompt(req.CaseContext, conversation.CaseID)

	// Build user prompt with conversation history
	userPrompt := s.buildConversationPrompt(history)

	// Call LLM
	startTime := time.Now()
	response, err := s.llm.Generate(ctx, userPrompt, systemPrompt)
	processingTime := time.Since(startTime).Milliseconds()

	if err != nil {
		s.logger.ErrorContext(ctx, "LLM generation failed", "error", err)
		return nil, fmt.Errorf("llm generation: %w", err)
	}

	// Create and encrypt assistant message
	assistantMsg := ai.NewChatMessage(conversation.ID, ai.ChatRoleAssistant, response)
	assistantMsgEncx, err := ai.ProcessChatMessageEncx(ctx, s.crypto, assistantMsg)
	if err != nil {
		return nil, fmt.Errorf("encrypt assistant message: %w", err)
	}

	if err := s.chatRepo.CreateMessage(ctx, assistantMsgEncx); err != nil {
		return nil, fmt.Errorf("create assistant message: %w", err)
	}

	// Update conversation stats
	conversation.MessageCount += 2
	conversation.TotalTokens += userMsg.TokenCount + assistantMsg.TokenCount
	conversation.UpdatedAt = time.Now()

	// Encrypt and update conversation
	convEncx, err := ai.ProcessChatConversationEncx(ctx, s.crypto, conversation)
	if err != nil {
		s.logger.WarnContext(ctx, "Failed to encrypt conversation for update", "error", err)
	} else {
		if err := s.chatRepo.UpdateConversation(ctx, convEncx); err != nil {
			s.logger.WarnContext(ctx, "Failed to update conversation stats", "error", err)
		}
	}

	s.logger.InfoContext(ctx, "Chat message processed",
		"conversation_id", conversation.ID,
		"processing_time_ms", processingTime,
		"message_count", conversation.MessageCount)

	return &ports.SendMessageResponse{
		ConversationID:    conversation.ID,
		UserMessageID:     userMsg.ID,
		AssistantMessage:  assistantMsg,
		IsNewConversation: isNewConversation,
	}, nil
}

func (s *Service) StreamMessage(ctx context.Context, req *ports.SendMessageRequest, streamCallback ports.StreamCallback) error {
	// Note: This is a simplified streaming implementation.
	// For full streaming support, the LLM client interface would need to be extended
	// to support streaming, and the Ollama client would need to use the streaming API.

	// For now, we'll use the non-streaming SendMessage and send the result as a single "chunk"
	streamCallback("", false, nil)

	resp, err := s.SendMessage(ctx, req)
	if err != nil {
		streamCallback("", false, err)
		return err
	}

	// Send the complete response
	streamCallback(resp.AssistantMessage.Content, true, nil)

	return nil
}

func (s *Service) GetConversation(ctx context.Context, conversationID uuid.UUID) (*ai.ChatConversation, []*ai.ChatMessage, error) {
	convEncx, err := s.chatRepo.GetConversation(ctx, conversationID)
	if err != nil {
		return nil, nil, err
	}

	conv, err := ai.DecryptChatConversationEncx(ctx, s.crypto, convEncx)
	if err != nil {
		return nil, nil, fmt.Errorf("decrypt conversation: %w", err)
	}

	messagesEncx, err := s.chatRepo.GetMessagesByConversation(ctx, conversationID)
	if err != nil {
		return nil, nil, err
	}

	messages := make([]*ai.ChatMessage, len(messagesEncx))
	for i, msgEncx := range messagesEncx {
		messages[i], err = ai.DecryptChatMessageEncx(ctx, s.crypto, msgEncx)
		if err != nil {
			return nil, nil, fmt.Errorf("decrypt message: %w", err)
		}
	}

	return conv, messages, nil
}

func (s *Service) ListConversations(ctx context.Context, caseID uuid.UUID) ([]*ai.ChatConversation, error) {
	convsEncx, err := s.chatRepo.ListConversationsByCase(ctx, caseID)
	if err != nil {
		return nil, err
	}

	conversations := make([]*ai.ChatConversation, len(convsEncx))
	for i, convEncx := range convsEncx {
		conversations[i], err = ai.DecryptChatConversationEncx(ctx, s.crypto, convEncx)
		if err != nil {
			return nil, fmt.Errorf("decrypt conversation: %w", err)
		}
	}

	return conversations, nil
}

func (s *Service) DeleteConversation(ctx context.Context, conversationID uuid.UUID) error {
	return s.chatRepo.DeleteConversation(ctx, conversationID)
}

func (s *Service) buildSystemPrompt(caseCtx *ai.ChatContext, caseID uuid.UUID) string {
	if caseCtx == nil {
		return "You are an AI assistant helping with an investigation case."
	}

	subjects := "(none recorded)"
	if len(caseCtx.Subjects) > 0 {
		var b strings.Builder
		for _, sub := range caseCtx.Subjects {
			b.WriteString(fmt.Sprintf("- %s: %s\n", sub.Role, sub.Name))
		}
		subjects = strings.TrimRight(b.String(), "\n")
	}

	notes := "(none)"
	if len(caseCtx.RecentNotes) > 0 {
		notes = strings.Join(caseCtx.RecentNotes, "\n")
	}

	return fmt.Sprintf(systemPromptTemplate,
		caseID,
		caseCtx.CaseTitle,
		caseCtx.CaseType,
		caseCtx.Status,
		caseCtx.ClientName,
		caseCtx.Location,
		subjects,
		notes,
	)
}

func (s *Service) buildConversationPrompt(messages []*ai.ChatMessage) string {
	var builder strings.Builder

	for _, msg := range messages {
		switch msg.Role {
		case ai.ChatRoleUser:
			builder.WriteString(fmt.Sprintf("User: %s\n", msg.Content))
		case ai.ChatRoleAssistant:
			builder.WriteString(fmt.Sprintf("Assistant: %s\n", msg.Content))
		}
	}

	builder.WriteString("Assistant:")
	return builder.String()
}

func (s *Service) generateConversationTitle(firstMessage string) string {
	// Truncate to first 50 chars, adding ellipsis if needed
	if len(firstMessage) <= 50 {
		return firstMessage
	}
	return firstMessage[:47] + "..."
}
