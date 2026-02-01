package aiChat

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hengadev/cluo_api/internal/domain/ai"
	"github.com/jackc/pgx/v5/pgxpool"
)

const Schema = "ai"

// Repository implements the ChatRepository interface.
type Repository struct {
	pool   *pgxpool.Pool
	schema string
}

// New creates a new chat repository.
func New(ctx context.Context, pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool:   pool,
		schema: Schema,
	}
}

// CreateConversation creates a new conversation.
func (r *Repository) CreateConversation(ctx context.Context, conv *ai.ChatConversationEncx) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.chat_conversations (
			id, case_id, title_encrypted, created_at, updated_at,
			created_by, message_count, total_tokens, dek_encrypted, key_version
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, r.schema)

	_, err := r.pool.Exec(ctx, query,
		conv.ID, conv.CaseID, conv.TitleEncrypted, conv.CreatedAt, conv.UpdatedAt,
		conv.CreatedBy, conv.MessageCount, conv.TotalTokens,
		conv.DEKEncrypted, conv.KeyVersion)

	if err != nil {
		return fmt.Errorf("create conversation: %w", err)
	}

	return nil
}

// GetConversation retrieves a conversation by ID.
func (r *Repository) GetConversation(ctx context.Context, id uuid.UUID) (*ai.ChatConversationEncx, error) {
	query := fmt.Sprintf(`
		SELECT id, case_id, title_encrypted, created_at, updated_at,
		       created_by, message_count, total_tokens, dek_encrypted, key_version
		FROM %s.chat_conversations
		WHERE id = $1
	`, r.schema)

	row := r.pool.QueryRow(ctx, query, id)

	var conv ai.ChatConversationEncx
	err := row.Scan(
		&conv.ID, &conv.CaseID, &conv.TitleEncrypted, &conv.CreatedAt, &conv.UpdatedAt,
		&conv.CreatedBy, &conv.MessageCount, &conv.TotalTokens,
		&conv.DEKEncrypted, &conv.KeyVersion,
	)

	if err != nil {
		return nil, fmt.Errorf("get conversation: %w", err)
	}

	return &conv, nil
}

// UpdateConversation updates a conversation.
func (r *Repository) UpdateConversation(ctx context.Context, conv *ai.ChatConversationEncx) error {
	query := fmt.Sprintf(`
		UPDATE %s.chat_conversations
		SET message_count = $2, total_tokens = $3, updated_at = $4
		WHERE id = $1
	`, r.schema)

	_, err := r.pool.Exec(ctx, query, conv.ID, conv.MessageCount, conv.TotalTokens, conv.UpdatedAt)

	if err != nil {
		return fmt.Errorf("update conversation: %w", err)
	}

	return nil
}

// DeleteConversation deletes a conversation by ID.
func (r *Repository) DeleteConversation(ctx context.Context, id uuid.UUID) error {
	query := fmt.Sprintf("DELETE FROM %s.chat_conversations WHERE id = $1", r.schema)

	_, err := r.pool.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("delete conversation: %w", err)
	}

	return nil
}

// ListConversationsByCase retrieves conversations for a case.
func (r *Repository) ListConversationsByCase(
	ctx context.Context,
	caseID uuid.UUID,
) ([]*ai.ChatConversationEncx, error) {
	query := fmt.Sprintf(`
		SELECT id, case_id, title_encrypted, created_at, updated_at,
		       created_by, message_count, total_tokens, dek_encrypted, key_version
		FROM %s.chat_conversations
		WHERE case_id = $1
		ORDER BY updated_at DESC
	`, r.schema)

	rows, err := r.pool.Query(ctx, query, caseID)
	if err != nil {
		return nil, fmt.Errorf("query conversations by case: %w", err)
	}
	defer rows.Close()

	var conversations []*ai.ChatConversationEncx
	for rows.Next() {
		var conv ai.ChatConversationEncx
		err := rows.Scan(
			&conv.ID, &conv.CaseID, &conv.TitleEncrypted, &conv.CreatedAt, &conv.UpdatedAt,
			&conv.CreatedBy, &conv.MessageCount, &conv.TotalTokens,
			&conv.DEKEncrypted, &conv.KeyVersion,
		)
		if err != nil {
			return nil, fmt.Errorf("scan conversation: %w", err)
		}

		conversations = append(conversations, &conv)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("iterate conversations: %w", rows.Err())
	}

	return conversations, nil
}

// CreateMessage creates a new message.
func (r *Repository) CreateMessage(ctx context.Context, msg *ai.ChatMessageEncx) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.chat_messages (
			id, conversation_id, role, content_encrypted,
			created_at, token_count, dek_encrypted, key_version
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, r.schema)

	_, err := r.pool.Exec(ctx, query,
		msg.ID, msg.ConversationID, msg.Role, msg.ContentEncrypted,
		msg.CreatedAt, msg.TokenCount, msg.DEKEncrypted, msg.KeyVersion)

	if err != nil {
		return fmt.Errorf("create message: %w", err)
	}

	return nil
}

// GetMessagesByConversation retrieves all messages for a conversation.
func (r *Repository) GetMessagesByConversation(
	ctx context.Context,
	conversationID uuid.UUID,
) ([]*ai.ChatMessageEncx, error) {
	query := fmt.Sprintf(`
		SELECT id, conversation_id, role, content_encrypted,
		       created_at, token_count, dek_encrypted, key_version
		FROM %s.chat_messages
		WHERE conversation_id = $1
		ORDER BY created_at ASC
	`, r.schema)

	rows, err := r.pool.Query(ctx, query, conversationID)
	if err != nil {
		return nil, fmt.Errorf("query messages by conversation: %w", err)
	}
	defer rows.Close()

	var messages []*ai.ChatMessageEncx
	for rows.Next() {
		var msg ai.ChatMessageEncx
		err := rows.Scan(
			&msg.ID, &msg.ConversationID, &msg.Role, &msg.ContentEncrypted,
			&msg.CreatedAt, &msg.TokenCount, &msg.DEKEncrypted, &msg.KeyVersion,
		)
		if err != nil {
			return nil, fmt.Errorf("scan message: %w", err)
		}

		messages = append(messages, &msg)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("iterate messages: %w", rows.Err())
	}

	return messages, nil
}

// GetRecentMessagesByConversation retrieves recent N messages for a conversation.
func (r *Repository) GetRecentMessagesByConversation(
	ctx context.Context,
	conversationID uuid.UUID,
	limit int,
) ([]*ai.ChatMessageEncx, error) {
	query := fmt.Sprintf(`
		SELECT id, conversation_id, role, content_encrypted,
		       created_at, token_count, dek_encrypted, key_version
		FROM %s.chat_messages
		WHERE conversation_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`, r.schema)

	rows, err := r.pool.Query(ctx, query, conversationID, limit)
	if err != nil {
		return nil, fmt.Errorf("query recent messages: %w", err)
	}
	defer rows.Close()

	var messages []*ai.ChatMessageEncx
	for rows.Next() {
		var msg ai.ChatMessageEncx
		err := rows.Scan(
			&msg.ID, &msg.ConversationID, &msg.Role, &msg.ContentEncrypted,
			&msg.CreatedAt, &msg.TokenCount, &msg.DEKEncrypted, &msg.KeyVersion,
		)
		if err != nil {
			return nil, fmt.Errorf("scan message: %w", err)
		}

		messages = append(messages, &msg)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("iterate messages: %w", rows.Err())
	}

	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}
