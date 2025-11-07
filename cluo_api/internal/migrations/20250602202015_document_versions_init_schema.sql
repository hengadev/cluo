-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Document versions table (generic for all document types)
CREATE TABLE IF NOT EXISTS document_versions (
    id BIGSERIAL PRIMARY KEY,
    document_id UUID NOT NULL,
    doc_type VARCHAR(20) NOT NULL CHECK (doc_type IN ('estimate', 'mandate', 'contract', 'invoice')),
    version INT NOT NULL,
    author_id UUID,
    data JSONB NOT NULL,
    reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(document_id, doc_type, version)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_document_versions_document_id ON document_versions(document_id);
CREATE INDEX IF NOT EXISTS idx_document_versions_doc_type ON document_versions(doc_type);
CREATE INDEX IF NOT EXISTS idx_document_versions_created_at ON document_versions(created_at);
CREATE INDEX IF NOT EXISTS idx_document_versions_author_id ON document_versions(author_id) WHERE author_id IS NOT NULL;

-- Add comments for documentation
COMMENT ON TABLE document_versions IS 'Stores version history for all document types';
COMMENT ON COLUMN document_versions.data IS 'Serialized document data in JSON format';
COMMENT ON COLUMN document_versions.reason IS 'Optional note explaining why this version was created';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE IF EXISTS document_versions;
