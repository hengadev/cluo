-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Create media schema if it doesn't exist
CREATE SCHEMA IF NOT EXISTS media;

-- Media files table
CREATE TABLE IF NOT EXISTS media.media_files (
    id UUID PRIMARY KEY,
    caseid UUID NOT NULL,
    filesize BIGINT NOT NULL,
    ispublished BOOLEAN NOT NULL DEFAULT FALSE,
    createdat TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    url_encrypted BYTEA NOT NULL,
    type_encrypted BYTEA NOT NULL,
    mimetype_encrypted BYTEA NOT NULL,
    filename_encrypted BYTEA NOT NULL,
    caption_encrypted BYTEA,
    dek_encrypted BYTEA NOT NULL,
    key_version INTEGER NOT NULL,
    metadata JSONB DEFAULT '{}',

    -- Foreign key constraints
    CONSTRAINT fk_media_files_caseid
        FOREIGN KEY (caseid) REFERENCES cases.cases(id) ON DELETE CASCADE
);

-- Media files table indexes
CREATE INDEX IF NOT EXISTS idx_media_files_caseid ON media.media_files(caseid);
CREATE INDEX IF NOT EXISTS idx_media_files_createdat ON media.media_files(createdat);
CREATE INDEX IF NOT EXISTS idx_media_files_key_version ON media.media_files(key_version);
CREATE INDEX IF NOT EXISTS idx_media_files_metadata ON media.media_files USING gin(metadata);

-- Media files table business constraints
ALTER TABLE media.media_files ADD CONSTRAINT chk_media_files_key_version_positive
    CHECK (key_version >= 0);

ALTER TABLE media.media_files ADD CONSTRAINT chk_media_files_filesize_positive
    CHECK (filesize > 0);

ALTER TABLE media.media_files ADD CONSTRAINT chk_media_files_caseid_not_empty
    CHECK (caseid IS NOT NULL);

-- Add comments for documentation
COMMENT ON SCHEMA media IS 'Media file management schema';
COMMENT ON TABLE media.media_files IS 'Stores encrypted media file information';
COMMENT ON COLUMN media.media_files.id IS 'Unique identifier for the media file';
COMMENT ON COLUMN media.media_files.caseid IS 'Foreign key reference to cases.cases.id';
COMMENT ON COLUMN media.media_files.filesize IS 'Size of the file in bytes';
COMMENT ON COLUMN media.media_files.ispublished IS 'Whether the media file is published';
COMMENT ON COLUMN media.media_files.createdat IS 'Media file creation timestamp';
COMMENT ON COLUMN media.media_files.url_encrypted IS 'Encrypted URL to the media file in S3';
COMMENT ON COLUMN media.media_files.type_encrypted IS 'Encrypted media type (image, video, audio)';
COMMENT ON COLUMN media.media_files.mimetype_encrypted IS 'Encrypted MIME type of the file';
COMMENT ON COLUMN media.media_files.filename_encrypted IS 'Encrypted original filename';
COMMENT ON COLUMN media.media_files.caption_encrypted IS 'Encrypted caption for the media file (nullable)';
COMMENT ON COLUMN media.media_files.dek_encrypted IS 'Encrypted Data Encryption Key for field-level encryption';
COMMENT ON COLUMN media.media_files.key_version IS 'Version of the Key Encryption Key used';
COMMENT ON COLUMN media.media_files.metadata IS 'Additional metadata in JSON format';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE IF EXISTS media.media_files;
DROP SCHEMA IF EXISTS media;
