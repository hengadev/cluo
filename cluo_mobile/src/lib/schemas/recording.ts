/**
 * ArkType validation schemas for recording-related API responses and requests.
 */

import { type } from "arktype";

/**
 * Recording status enum values.
 */
const RecordingStatusValues = ["uploading", "transcribing", "analyzing", "completed", "failed"] as const;

/**
 * Processing step status values.
 */
const StepStatusValues = ["pending", "processing", "completed", "failed"] as const;

/**
 * Suggestion category values.
 */
const SuggestionCategoryValues = ["observations", "statements", "actions", "unclear"] as const;

/**
 * Schema for a single processing step.
 */
export const ProcessingStepSchema = type({
	title: "string",
	status: type(StepStatusValues),
});

/**
 * Schema for recording status response from API.
 */
export const RecordingStatusResponseSchema = type({
	id: "string",
	status: type(RecordingStatusValues),
	processingSteps: ProcessingStepSchema.array(),
	error: "string?",
});

/**
 * Schema for transcript response from API.
 */
export const TranscriptResponseSchema = type({
	recordingId: "string",
	text: "string.min(1)",
	confidence: "number.between(0, 1)?",
	isConfirmed: "boolean",
	createdAt: "string",
	updatedAt: "string",
});

/**
 * Schema for confirming/updating transcript request.
 */
export const ConfirmTranscriptRequestSchema = type({
	text: "string.min(1).max(100000)", // Max 100k characters
});

/**
 * Schema for transcript text input validation (client-side).
 * More lenient than API version for better UX.
 */
export const TranscriptTextInputSchema = type({
	text: "string.max(100000)",
});

/**
 * Schema for a single suggestion.
 */
export const SuggestionSchema = type({
	id: "string",
	category: type(SuggestionCategoryValues),
	text: "string.min(1)",
	selected: "boolean",
	timestamp: "string?",
});

/**
 * Schema for analysis response from API.
 * Matches the backend GET /ai/analysis/{id} shape.
 */
export const AnalysisResponseSchema = type({
	id: "string",
	transcriptionId: "string",
	keyFindings: "string",
	summary: "string",
	sentiment: "string",
	topics: "string",
	suggestedActions: "string",
	createdAt: "string",
});

/**
 * Schema for upload recording response from API.
 */
export const UploadRecordingResponseSchema = type({
	id: "string",
	status: type(RecordingStatusValues),
});

/**
 * Schema for recordings list response from API.
 */
export const RecordingsListResponseSchema = type({
	recordings: type({
		id: "string",
		"caseId?": "string",
		title: "string",
		date: "string",
		startTime: "string",
		duration: "number",
		"fileSize?": "number",
		status: type(RecordingStatusValues),
		"processingSteps?": ProcessingStepSchema.array(),
	}).array(),
	totalCount: "number",
});

/**
 * Schema for recording metadata (from list or detail).
 */
export const RecordingSchema = type({
	id: "string",
	"caseId?": "string",
	title: "string",
	date: "string",
	startTime: "string",
	duration: "number",
	"fileSize?": "number",
	status: type(RecordingStatusValues),
	"processingSteps?": ProcessingStepSchema.array(),
});

/**
 * Schema for audio file validation before upload.
 */
export const AudioFileSchema = type({
	type: 'string.startsWith("audio/")',
	size: "number.between(1, 500000000)", // Max 500MB
});

/**
 * Schema for recording ID parameter validation.
 */
export const RecordingIdParamSchema = type({
	id: "string.minLength(1)",
});

/**
 * Type exports inferred from schemas.
 */
export type RecordingStatusResponse = typeof RecordingStatusResponseSchema.infer;
export type TranscriptResponse = typeof TranscriptResponseSchema.infer;
export type ConfirmTranscriptRequest = typeof ConfirmTranscriptRequestSchema.infer;
export type AnalysisResponse = typeof AnalysisResponseSchema.infer;
export type UploadRecordingResponse = typeof UploadRecordingResponseSchema.infer;
export type RecordingsListResponse = typeof RecordingsListResponseSchema.infer;
export type Recording = typeof RecordingSchema.infer;
