-- +goose Up
ALTER TABLE media.media_files
    ADD COLUMN purpose TEXT NOT NULL DEFAULT 'general';

ALTER TABLE media.media_files
    ADD CONSTRAINT chk_media_files_purpose
        CHECK (purpose IN ('general', 'witness_interview'));

COMMENT ON COLUMN media.media_files.purpose IS 'Recording purpose: general (for the report) or witness_interview';

-- +goose Down
ALTER TABLE media.media_files
    DROP CONSTRAINT IF EXISTS chk_media_files_purpose;

ALTER TABLE media.media_files
    DROP COLUMN IF EXISTS purpose;
