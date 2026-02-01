package aiChatHandler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/common/ctxutil"
	"github.com/hengadev/cluo_api/internal/common/httpx"
	"github.com/hengadev/cluo_api/internal/domain/ai"
	"github.com/hengadev/cluo_api/internal/ports"
	mw "github.com/hengadev/cluo_api/internal/common/middleware/auth"
)

type Handler interface {
	SendMessage(w http.ResponseWriter, r *http.Request)
	GetConversation(w http.ResponseWriter, r *http.Request)
	ListConversations(w http.ResponseWriter, r *http.Request)
	DeleteConversation(w http.ResponseWriter, r *http.Request)
	RegisterRoutes(router *http.ServeMux)
}

type handler struct {
	svc    ports.ChatService
	authmw mw.AuthMiddleware
}

func New(svc ports.ChatService, authmw mw.AuthMiddleware) Handler {
	return &handler{svc: svc, authmw: authmw}
}

// Request/Response DTOs
type SendMessageRequest struct {
	ConversationID *uuid.UUID `json:"conversationId,omitempty"`
	Message        string    `json:"message"`
}

type SendMessageResponse struct {
	ConversationID     string                 `json:"conversationId"`
	UserMessageID      string                 `json:"userMessageId"`
	AssistantMessage   ai.ChatMessageDTO      `json:"assistantMessage"`
	Conversation       *ai.ChatConversationDTO `json:"conversation,omitempty"`
}

type GetConversationResponse struct {
	Conversation ai.ChatConversationDTO `json:"conversation"`
	Messages     []ai.ChatMessageDTO    `json:"messages"`
}

type ListConversationsResponse struct {
	Conversations []ai.ChatConversationDTO `json:"conversations"`
}

func (h *handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Get user ID from context
	userID, err := ctxutil.GetUserIDFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusUnauthorized)
		return
	}

	// Get case ID from query param
	caseIDStr := r.URL.Query().Get("case_id")
	if caseIDStr == "" {
		httpx.RespondWithError(w, fmt.Errorf("case_id is required"), http.StatusBadRequest)
		return
	}
	caseID, err := uuid.Parse(caseIDStr)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	// Parse request body
	var payload SendMessageRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&payload); err != nil {
		logger.WarnContext(ctx, "Handler: Invalid JSON request body",
			"error", err,
			"operation", "send_message",
			"method", r.Method,
			"path", r.URL.Path)
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing send message request",
		"operation", "send_message",
		"method", r.Method,
		"path", r.URL.Path,
		"case_id", caseID,
		"message_length", len(payload.Message))

	// Build case context (would fetch from case service)
	caseContext := h.buildCaseContext(ctx, caseID)

	// Call service
	req := &ports.SendMessageRequest{
		ConversationID: payload.ConversationID,
		CaseID:         caseID,
		UserID:         userID,
		Message:        payload.Message,
		CaseContext:    caseContext,
	}

	resp, err := h.svc.SendMessage(ctx, req)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "send message")
		return
	}

	response := &SendMessageResponse{
		ConversationID:   resp.ConversationID.String(),
		UserMessageID:    resp.UserMessageID.String(),
		AssistantMessage: resp.AssistantMessage.ToDTO(),
	}

	// Include conversation details if new
	if resp.IsNewConversation {
		conv, _, _ := h.svc.GetConversation(ctx, resp.ConversationID)
		dto := conv.ToDTO()
		response.Conversation = &dto
	}

	logger.InfoContext(ctx, "Handler: Send message request completed successfully",
		"operation", "send_message",
		"method", r.Method,
		"path", r.URL.Path,
		"status_code", http.StatusOK)

	httpx.RespondWithJSON(w, response, http.StatusOK)
}

func (h *handler) GetConversation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Get conversation ID from URL path
	conversationIDStr := r.PathValue("id")
	if conversationIDStr == "" {
		httpx.RespondWithError(w, fmt.Errorf("conversation id is required"), http.StatusBadRequest)
		return
	}
	conversationID, err := uuid.Parse(conversationIDStr)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing get conversation request",
		"operation", "get_conversation",
		"method", r.Method,
		"path", r.URL.Path,
		"conversation_id", conversationID)

	conv, messages, err := h.svc.GetConversation(ctx, conversationID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "get conversation")
		return
	}

	response := &GetConversationResponse{
		Conversation: conv.ToDTO(),
		Messages:     ai.MessagesToDTOs(messages),
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}

func (h *handler) ListConversations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Get case ID from query param
	caseIDStr := r.URL.Query().Get("case_id")
	if caseIDStr == "" {
		httpx.RespondWithError(w, fmt.Errorf("case_id is required"), http.StatusBadRequest)
		return
	}
	caseID, err := uuid.Parse(caseIDStr)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing list conversations request",
		"operation", "list_conversations",
		"method", r.Method,
		"path", r.URL.Path,
		"case_id", caseID)

	conversations, err := h.svc.ListConversations(ctx, caseID)
	if err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "list conversations")
		return
	}

	response := &ListConversationsResponse{
		Conversations: ai.ConversationsToDTOs(conversations),
	}

	httpx.RespondWithJSON(w, response, http.StatusOK)
}

func (h *handler) DeleteConversation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger, err := ctxutil.GetLoggerFromContext(ctx)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusInternalServerError)
		return
	}

	// Get conversation ID from URL path
	conversationIDStr := r.PathValue("id")
	if conversationIDStr == "" {
		httpx.RespondWithError(w, fmt.Errorf("conversation id is required"), http.StatusBadRequest)
		return
	}
	conversationID, err := uuid.Parse(conversationIDStr)
	if err != nil {
		httpx.RespondWithError(w, err, http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "Handler: Processing delete conversation request",
		"operation", "delete_conversation",
		"method", r.Method,
		"path", r.URL.Path,
		"conversation_id", conversationID)

	if err := h.svc.DeleteConversation(ctx, conversationID); err != nil {
		httpx.RespondWithServiceError(w, logger, ctx, err, "delete conversation")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) buildCaseContext(ctx context.Context, caseID uuid.UUID) *ai.ChatContext {
	// TODO: Fetch case details from case service
	// For now, return minimal context
	return &ai.ChatContext{
		CaseID: caseID,
	}
}

func (h *handler) RegisterRoutes(router *http.ServeMux) {
	// Chat endpoints (require authenticated user)
	router.HandleFunc("POST /api/ai/chat/message", h.authmw.RequireAccessToken(h.SendMessage))
	router.HandleFunc("GET /api/ai/chat/conversations", h.authmw.RequireAccessToken(h.ListConversations))
	router.HandleFunc("GET /api/ai/chat/conversations/{id}", h.authmw.RequireAccessToken(h.GetConversation))
	router.HandleFunc("DELETE /api/ai/chat/conversations/{id}", h.authmw.RequireAccessToken(h.DeleteConversation))
}
