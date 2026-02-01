-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- ============================================================================
-- Chat Conversations Table
-- ============================================================================
-- Stores chat sessions associated with cases for AI discussions

CREATE TABLE IF NOT EXISTS ai.chat_conversations (
    id UUID PRIMARY KEY,
    case_id UUID NOT NULL,
    title_encrypted BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID NOT NULL,
    message_count INTEGER NOT NULL DEFAULT 0,
    total_tokens INTEGER NOT NULL DEFAULT 0,
    dek_encrypted BYTEA NOT NULL,
    key_version INTEGER NOT NULL,

    -- Foreign key constraints
    CONSTRAINT fk_chat_conversations_case
        FOREIGN KEY (case_id) REFERENCES cases.cases(id) ON DELETE CASCADE,
    CONSTRAINT fk_chat_conversations_created_by
        FOREIGN KEY (created_by) REFERENCES auth.users(id) ON DELETE CASCADE,

    -- Check constraints
    CONSTRAINT chk_chat_conversations_key_version
        CHECK (key_version >= 0),
    CONSTRAINT chk_chat_conversations_message_count
        CHECK (message_count >= 0),
    CONSTRAINT chk_chat_conversations_total_tokens
        CHECK (total_tokens >= 0)
);

-- ============================================================================
-- Chat Messages Table
-- ============================================================================
-- Stores individual messages within chat conversations

CREATE TABLE IF NOT EXISTS ai.chat_messages (
    id UUID PRIMARY KEY,
    conversation_id UUID NOT NULL,
    role VARCHAR(20) NOT NULL,
    content_encrypted BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    token_count INTEGER,

    -- Foreign key constraints
    CONSTRAINT fk_chat_messages_conversation
        FOREIGN KEY (conversation_id)
        REFERENCES ai.chat_conversations(id) ON DELETE CASCADE,

    -- Check constraints
    CONSTRAINT chk_chat_messages_role
        CHECK (role IN ('user', 'assistant', 'system')),
    CONSTRAINT chk_chat_messages_token_count
        CHECK (token_count IS NULL OR token_count >= 0)
);

-- ============================================================================
-- Indexes: Chat Conversations
-- ============================================================================

-- Index for case conversations (most recent first)
CREATE INDEX IF NOT EXISTS idx_chat_conversations_case
    ON ai.chat_conversations (case_id, updated_at DESC);

-- Index for user's conversations (most recent first)
CREATE INDEX IF NOT EXISTS idx_chat_conversations_created_by
    ON ai.chat_conversations (created_by, updated_at DESC);

-- ============================================================================
-- Indexes: Chat Messages
-- ============================================================================

-- Index for conversation messages (chronological order)
CREATE INDEX IF NOT EXISTS idx_chat_messages_conversation
    ON ai.chat_messages (conversation_id, created_at ASC);

-- ============================================================================
-- Comments for Documentation
-- ============================================================================

COMMENT ON TABLE ai.chat_conversations IS 'Stores chat sessions for AI-assisted case discussions';
COMMENT ON TABLE ai.chat_messages IS 'Stores individual messages within chat conversations';

-- Chat Conversations columns
COMMENT ON COLUMN ai.chat_conversations.id IS 'Unique identifier for the conversation';
COMMENT ON COLUMN ai.chat_conversations.case_id IS 'Reference to the associated case';
COMMENT ON COLUMN ai.chat_conversations.title_encrypted IS 'Encrypted conversation title (typically first message)';
COMMENT ON COLUMN ai.chat_conversations.created_at IS 'Conversation creation timestamp';
COMMENT ON COLUMN ai.chat_conversations.updated_at IS 'Last update timestamp';
COMMENT ON COLUMN ai.chat_conversations.created_by IS 'User who created the conversation';
COMMENT ON COLUMN ai.chat_conversations.message_count IS 'Total number of messages in the conversation';
COMMENT ON COLUMN ai.chat_conversations.total_tokens IS 'Estimated total tokens used';
COMMENT ON COLUMN ai.chat_conversations.dek_encrypted IS 'Encrypted Data Encryption Key for field-level encryption';
COMMENT ON COLUMN ai.chat_conversations.key_version IS 'Version of the Key Encryption Key used';

-- Chat Messages columns
COMMENT ON COLUMN ai.chat_messages.id IS 'Unique identifier for the message';
COMMENT ON COLUMN ai.chat_messages.conversation_id IS 'Reference to the parent conversation';
COMMENT ON COLUMN ai.chat_messages.role IS 'Message sender role: user, assistant, or system';
COMMENT ON COLUMN ai.chat_messages.content_encrypted IS 'Encrypted message content';
COMMENT ON COLUMN ai.chat_messages.created_at IS 'Message creation timestamp';
COMMENT ON COLUMN ai.chat_messages.token_count IS 'Estimated token count for the message';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

-- Drop indexes
DROP INDEX IF EXISTS ai.idx_chat_messages_conversation;
DROP INDEX IF EXISTS ai.idx_chat_conversations_created_by;
DROP INDEX IF EXISTS ai.idx_chat_conversations_case;

-- Drop tables
DROP TABLE IF EXISTS ai.chat_messages;
DROP TABLE IF EXISTS ai.chat_conversations;
