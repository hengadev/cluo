-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- Create ai schema if it doesn't exist
CREATE SCHEMA IF NOT EXISTS ai;

-- ============================================================================
-- Transcription Jobs Table
-- ============================================================================
-- Tracks async speech-to-text transcription jobs with status and progress

CREATE TABLE IF NOT EXISTS ai.transcription_jobs (
    id UUID PRIMARY KEY,
    media_file_id UUID NOT NULL,
    audio_path TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    progress INTEGER NOT NULL DEFAULT 0,
    error_message TEXT,
    transcription_id UUID,
    webhook_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    claimed_at TIMESTAMPTZ,
    claimed_by VARCHAR(100),
    created_by UUID NOT NULL,

    -- Foreign key constraints
    CONSTRAINT fk_jobs_media_file
        FOREIGN KEY (media_file_id) REFERENCES media.media_files(id) ON DELETE CASCADE,
    CONSTRAINT fk_jobs_created_by
        FOREIGN KEY (created_by) REFERENCES auth.users(id) ON DELETE CASCADE,

    -- Check constraints
    CONSTRAINT chk_jobs_status
        CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'cancelled')),
    CONSTRAINT chk_jobs_progress_range
        CHECK (progress >= 0 AND progress <= 100),
    CONSTRAINT chk_jobs_key_version_positive
        CHECK (key_version >= 0)
);

-- ============================================================================
-- Transcriptions Table
-- ============================================================================
-- Stores completed speech-to-text transcription results

CREATE TABLE IF NOT EXISTS ai.transcriptions (
    id UUID PRIMARY KEY,
    job_id UUID NOT NULL UNIQUE,
    media_file_id UUID NOT NULL,
    audio_url TEXT NOT NULL,
    transcript_encrypted BYTEA NOT NULL,
    confidence_score REAL NOT NULL,
    language VARCHAR(10) NOT NULL,
    duration BIGINT NOT NULL,
    model_name VARCHAR(100) NOT NULL,
    processing_time_ms BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    audio_deleted_at TIMESTAMPTZ,
    dek_encrypted BYTEA NOT NULL,
    key_version INTEGER NOT NULL,

    -- Foreign key constraints
    CONSTRAINT fk_transcriptions_job
        FOREIGN KEY (job_id) REFERENCES ai.transcription_jobs(id) ON DELETE CASCADE,
    CONSTRAINT fk_transcriptions_media_file
        FOREIGN KEY (media_file_id) REFERENCES media.media_files(id) ON DELETE CASCADE,

    -- Check constraints
    CONSTRAINT chk_transcriptions_confidence_range
        CHECK (confidence_score >= 0.0 AND confidence_score <= 1.0),
    CONSTRAINT chk_transcriptions_duration_positive
        CHECK (duration >= 0),
    CONSTRAINT chk_transcriptions_key_version_positive
        CHECK (key_version >= 0)
);

-- ============================================================================
-- Text Transformations Table
-- ============================================================================
-- Stores text transformation results (reword, summarize, formalize, clarify)

CREATE TABLE IF NOT EXISTS ai.text_transformations (
    id UUID PRIMARY KEY,
    input_text_encrypted BYTEA NOT NULL,
    output_text_encrypted BYTEA NOT NULL,
    transformation_type VARCHAR(20) NOT NULL,
    model_used VARCHAR(100) NOT NULL,
    input_length INTEGER NOT NULL,
    output_length INTEGER NOT NULL,
    processing_time_ms BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Check constraints
    CONSTRAINT chk_transformations_type
        CHECK (transformation_type IN ('reword', 'summarize', 'formalize', 'clarify')),
    CONSTRAINT chk_transformations_lengths_positive
        CHECK (input_length >= 0 AND output_length >= 0),
    CONSTRAINT chk_transformations_processing_time_positive
        CHECK (processing_time_ms >= 0)
);

-- ============================================================================
-- Transcript Analyses Table
-- ============================================================================
-- Stores analysis results for transcriptions

CREATE TABLE IF NOT EXISTS ai.transcript_analyses (
    id UUID PRIMARY KEY,
    transcription_id UUID NOT NULL UNIQUE,
    key_findings_encrypted BYTEA NOT NULL,
    summary_encrypted BYTEA NOT NULL,
    sentiment VARCHAR(20) NOT NULL,
    topics_encrypted BYTEA NOT NULL,
    suggested_actions_encrypted BYTEA NOT NULL,
    model_used VARCHAR(100) NOT NULL,
    processing_time_ms BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Foreign key constraints
    CONSTRAINT fk_analyses_transcription
        FOREIGN KEY (transcription_id) REFERENCES ai.transcriptions(id) ON DELETE CASCADE,

    -- Check constraints
    CONSTRAINT chk_analyses_sentiment
        CHECK (sentiment IN ('positive', 'neutral', 'negative', 'mixed')),
    CONSTRAINT chk_analyses_processing_time_positive
        CHECK (processing_time_ms >= 0)
);

-- ============================================================================
-- Report Suggestions Table
-- ============================================================================
-- Stores report suggestions based on transcript analysis

CREATE TABLE IF NOT EXISTS ai.report_suggestions (
    id UUID PRIMARY KEY,
    analysis_id UUID NOT NULL,
    report_type VARCHAR(50) NOT NULL,
    reasoning_encrypted BYTEA NOT NULL,
    confidence REAL NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Foreign key constraints
    CONSTRAINT fk_suggestions_analysis
        FOREIGN KEY (analysis_id) REFERENCES ai.transcript_analyses(id) ON DELETE CASCADE,

    -- Check constraints
    CONSTRAINT chk_suggestions_confidence_range
        CHECK (confidence >= 0.0 AND confidence <= 1.0)
);

-- ============================================================================
-- Indexes: Transcription Jobs
-- ============================================================================

-- Index for polling pending jobs (used by workers)
CREATE INDEX IF NOT EXISTS idx_transcription_jobs_pending
    ON ai.transcription_jobs (status, created_at)
    WHERE status = 'pending';

-- Index for user's jobs (most recent first)
CREATE INDEX IF NOT EXISTS idx_transcription_jobs_user_created
    ON ai.transcription_jobs (created_by, created_at DESC);

-- Index for media file lookups
CREATE INDEX IF NOT EXISTS idx_transcription_jobs_media_file
    ON ai.transcription_jobs (media_file_id);

-- Index for webhook processing
CREATE INDEX IF NOT EXISTS idx_transcription_jobs_webhook
    ON ai.transcription_jobs (webhook_url)
    WHERE webhook_url IS NOT NULL;

-- Index for claimed jobs (worker recovery)
CREATE INDEX IF NOT EXISTS idx_transcription_jobs_claimed
    ON ai.transcription_jobs (status, claimed_at)
    WHERE status = 'processing';

-- ============================================================================
-- Indexes: Transcriptions
-- ============================================================================

-- Index for job lookups
CREATE INDEX IF NOT EXISTS idx_transcriptions_job
    ON ai.transcriptions (job_id);

-- Index for media file lookups
CREATE INDEX IF NOT EXISTS idx_transcriptions_media_file
    ON ai.transcriptions (media_file_id);

-- Index for language filtering
CREATE INDEX IF NOT EXISTS idx_transcriptions_language
    ON ai.transcriptions (language);

-- Index for created date filtering
CREATE INDEX IF NOT EXISTS idx_transcriptions_created_at
    ON ai.transcriptions (created_at DESC);

-- ============================================================================
-- Indexes: Text Transformations
-- ============================================================================

-- Index for transformation type filtering
CREATE INDEX IF NOT EXISTS idx_text_transformations_type
    ON ai.text_transformations (transformation_type);

-- Index for created date filtering
CREATE INDEX IF NOT EXISTS idx_text_transformations_created_at
    ON ai.text_transformations (created_at DESC);

-- ============================================================================
-- Indexes: Transcript Analyses
-- ============================================================================

-- Index for transcription lookups
CREATE INDEX IF NOT EXISTS idx_transcript_analyses_transcription
    ON ai.transcript_analyses (transcription_id);

-- Index for sentiment filtering
CREATE INDEX IF NOT EXISTS idx_transcript_analyses_sentiment
    ON ai.transcript_analyses (sentiment);

-- Index for created date filtering
CREATE INDEX IF NOT EXISTS idx_transcript_analyses_created_at
    ON ai.transcript_analyses (created_at DESC);

-- ============================================================================
-- Indexes: Report Suggestions
-- ============================================================================

-- Index for analysis lookups
CREATE INDEX IF NOT EXISTS idx_report_suggestions_analysis
    ON ai.report_suggestions (analysis_id);

-- Index for report type filtering
CREATE INDEX IF NOT EXISTS idx_report_suggestions_type
    ON ai.report_suggestions (report_type);

-- ============================================================================
-- Comments for Documentation
-- ============================================================================

COMMENT ON SCHEMA ai IS 'AI services schema - text transformation, speech-to-text, and transcript analysis';

COMMENT ON TABLE ai.transcription_jobs IS 'Tracks async speech-to-text transcription jobs with status and progress';
COMMENT ON TABLE ai.transcriptions IS 'Stores completed speech-to-text transcription results';
COMMENT ON TABLE ai.text_transformations IS 'Stores text transformation results (reword, summarize, formalize, clarify)';
COMMENT ON TABLE ai.transcript_analyses IS 'Stores analysis results for transcriptions';
COMMENT ON TABLE ai.report_suggestions IS 'Stores report suggestions based on transcript analysis';

-- Transcription Jobs columns
COMMENT ON COLUMN ai.transcription_jobs.id IS 'Unique identifier for the job';
COMMENT ON COLUMN ai.transcription_jobs.media_file_id IS 'Reference to the media file being transcribed';
COMMENT ON COLUMN ai.transcription_jobs.audio_path IS 'Temporary path to the audio file';
COMMENT ON COLUMN ai.transcription_jobs.status IS 'Job status: pending, processing, completed, failed, cancelled';
COMMENT ON COLUMN ai.transcription_jobs.progress IS 'Progress percentage (0-100)';
COMMENT ON COLUMN ai.transcription_jobs.error_message IS 'Error message if the job failed';
COMMENT ON COLUMN ai.transcription_jobs.transcription_id IS 'Reference to the completed transcription';
COMMENT ON COLUMN ai.transcription_jobs.webhook_url IS 'Optional callback URL for job completion';
COMMENT ON COLUMN ai.transcription_jobs.created_at IS 'Job creation timestamp';
COMMENT ON COLUMN ai.transcription_jobs.started_at IS 'Job processing start timestamp';
COMMENT ON COLUMN ai.transcription_jobs.completed_at IS 'Job completion timestamp';
COMMENT ON COLUMN ai.transcription_jobs.claimed_at IS 'When the job was claimed by a worker';
COMMENT ON COLUMN ai.transcription_jobs.claimed_by IS 'Worker ID that claimed the job';
COMMENT ON COLUMN ai.transcription_jobs.created_by IS 'User who submitted the job';

-- Transcriptions columns
COMMENT ON COLUMN ai.transcriptions.id IS 'Unique identifier for the transcription';
COMMENT ON COLUMN ai.transcriptions.job_id IS 'Reference to the transcription job';
COMMENT ON COLUMN ai.transcriptions.media_file_id IS 'Reference to the original media file';
COMMENT ON COLUMN ai.transcriptions.audio_url IS 'URL to the audio file';
COMMENT ON COLUMN ai.transcriptions.transcript_encrypted IS 'Encrypted transcript text';
COMMENT ON COLUMN ai.transcriptions.confidence_score IS 'Confidence score (0.0-1.0)';
COMMENT ON COLUMN ai.transcriptions.language IS 'Detected language code (e.g., en, es, fr)';
COMMENT ON COLUMN ai.transcriptions.duration IS 'Audio duration in milliseconds';
COMMENT ON COLUMN ai.transcriptions.model_name IS 'Whisper model used (e.g., base, small, medium)';
COMMENT ON COLUMN ai.transcriptions.processing_time_ms IS 'Processing time in milliseconds';
COMMENT ON COLUMN ai.transcriptions.created_at IS 'Transcription creation timestamp';
COMMENT ON COLUMN ai.transcriptions.audio_deleted_at IS 'When the audio file was deleted';
COMMENT ON COLUMN ai.transcriptions.dek_encrypted IS 'Encrypted Data Encryption Key for field-level encryption';
COMMENT ON COLUMN ai.transcriptions.key_version IS 'Version of the Key Encryption Key used';

-- Text Transformations columns
COMMENT ON COLUMN ai.text_transformations.id IS 'Unique identifier for the transformation';
COMMENT ON COLUMN ai.text_transformations.input_text_encrypted IS 'Encrypted input text';
COMMENT ON COLUMN ai.text_transformations.output_text_encrypted IS 'Encrypted output text';
COMMENT ON COLUMN ai.text_transformations.transformation_type IS 'Type: reword, summarize, formalize, clarify';
COMMENT ON COLUMN ai.text_transformations.model_used IS 'LLM model used (e.g., llama3.2)';
COMMENT ON COLUMN ai.text_transformations.input_length IS 'Input text character count';
COMMENT ON COLUMN ai.text_transformations.output_length IS 'Output text character count';
COMMENT ON COLUMN ai.text_transformations.processing_time_ms IS 'Processing time in milliseconds';
COMMENT ON COLUMN ai.text_transformations.created_at IS 'Transformation creation timestamp';

-- Transcript Analyses columns
COMMENT ON COLUMN ai.transcript_analyses.id IS 'Unique identifier for the analysis';
COMMENT ON COLUMN ai.transcript_analyses.transcription_id IS 'Reference to the transcription';
COMMENT ON COLUMN ai.transcript_analyses.key_findings_encrypted IS 'Encrypted key findings from the transcript';
COMMENT ON COLUMN ai.transcript_analyses.summary_encrypted IS 'Encrypted summary of the transcript';
COMMENT ON COLUMN ai.transcript_analyses.sentiment IS 'Detected sentiment: positive, neutral, negative, mixed';
COMMENT ON COLUMN ai.transcript_analyses.topics_encrypted IS 'Encrypted topics (JSON array)';
COMMENT ON COLUMN ai.transcript_analyses.suggested_actions_encrypted IS 'Encrypted suggested actions';
COMMENT ON COLUMN ai.transcript_analyses.model_used IS 'LLM model used for analysis';
COMMENT ON COLUMN ai.transcript_analyses.processing_time_ms IS 'Processing time in milliseconds';
COMMENT ON COLUMN ai.transcript_analyses.created_at IS 'Analysis creation timestamp';

-- Report Suggestions columns
COMMENT ON COLUMN ai.report_suggestions.id IS 'Unique identifier for the suggestion';
COMMENT ON COLUMN ai.report_suggestions.analysis_id IS 'Reference to the transcript analysis';
COMMENT ON COLUMN ai.report_suggestions.report_type IS 'Suggested report type';
COMMENT ON COLUMN ai.report_suggestions.reasoning_encrypted IS 'Encrypted reasoning for the suggestion';
COMMENT ON COLUMN ai.report_suggestions.confidence IS 'Confidence score (0.0-1.0)';
COMMENT ON COLUMN ai.report_suggestions.created_at IS 'Suggestion creation timestamp';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

-- Drop indexes
DROP INDEX IF EXISTS ai.idx_report_suggestions_type;
DROP INDEX IF EXISTS ai.idx_report_suggestions_analysis;
DROP INDEX IF EXISTS ai.idx_transcript_analyses_created_at;
DROP INDEX IF EXISTS ai.idx_transcript_analyses_sentiment;
DROP INDEX IF EXISTS ai.idx_transcript_analyses_transcription;
DROP INDEX IF EXISTS ai.idx_text_transformations_created_at;
DROP INDEX IF EXISTS ai.idx_text_transformations_type;
DROP INDEX IF EXISTS ai.idx_transcriptions_created_at;
DROP INDEX IF EXISTS ai.idx_transcriptions_language;
DROP INDEX IF EXISTS ai.idx_transcriptions_media_file;
DROP INDEX IF EXISTS ai.idx_transcriptions_job;
DROP INDEX IF EXISTS ai.idx_transcription_jobs_claimed;
DROP INDEX IF EXISTS ai.idx_transcription_jobs_webhook;
DROP INDEX IF EXISTS ai.idx_transcription_jobs_media_file;
DROP INDEX IF EXISTS ai.idx_transcription_jobs_user_created;
DROP INDEX IF EXISTS ai.idx_transcription_jobs_pending;

-- Drop tables
DROP TABLE IF EXISTS ai.report_suggestions;
DROP TABLE IF EXISTS ai.transcript_analyses;
DROP TABLE IF EXISTS ai.text_transformations;
DROP TABLE IF EXISTS ai.transcriptions;
DROP TABLE IF EXISTS ai.transcription_jobs;

-- Drop schema
DROP SCHEMA IF EXISTS ai;
