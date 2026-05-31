/**
 * Recording types for the audio transcription and analysis workflow.
 * All types are designed to work with the Golang backend API.
 */

/**
 * Processing status of a recording.
 */
export type RecordingStatus =
	| "uploading"
	| "transcribing"
	| "analyzing"
	| "completed"
	| "failed";

/**
 * Individual processing step with its status.
 */
export interface ProcessingStep {
	title: string;
	status: "pending" | "processing" | "completed" | "failed";
}

/**
 * Main recording metadata and state.
 */
export interface Recording {
	id: string;
	caseId?: string;
	title: string;
	date: string; // Formatted date string (e.g., "30 Jan, 2026")
	startTime: string; // HH:MM format
	duration: number | string; // Duration in seconds (number) or formatted "MM:SS" (string)
	fileSize?: number; // File size in bytes
	status: RecordingStatus;
	processingSteps?: ProcessingStep[];
}

/**
 * Transcript data for a recording.
 */
export interface Transcript {
	recordingId: string;
	text: string;
	confidence?: number; // Overall confidence score (0-1)
	isConfirmed: boolean; // Whether user has confirmed the transcript
	createdAt: string;
	updatedAt: string;
}

/**
 * Suggestion category for analysis results.
 */
export type SuggestionCategory = "observations" | "statements" | "actions" | "unclear";

/**
 * Individual suggestion from AI analysis.
 */
export interface Suggestion {
	id: string;
	category: SuggestionCategory;
	text: string;
	selected: boolean; // User selection state
	timestamp?: string; // Reference point in transcript
}

/**
 * Analysis result from AI processing of transcript.
 * Matches the backend GET /ai/analysis/{id} response shape.
 */
export interface AnalysisResult {
	id: string;
	transcriptionId: string;
	keyFindings: string;
	summary: string;
	sentiment: string;
	topics: string; // JSON-encoded array string from the backend
	suggestedActions: string;
	createdAt: string;
}

/**
 * Response from upload API call.
 */
export interface UploadRecordingResponse {
	id: string;
	status: RecordingStatus;
}

/**
 * Response from recording status API call.
 */
export interface RecordingStatusResponse {
	id: string;
	status: RecordingStatus;
	processingSteps: ProcessingStep[];
	error?: string;
}

/**
 * Response from transcript API call.
 */
export interface TranscriptResponse {
	recordingId: string;
	text: string;
	confidence?: number;
	isConfirmed: boolean;
	createdAt: string;
	updatedAt: string;
}

/**
 * Request body for confirming/updating transcript.
 */
export interface ConfirmTranscriptRequest {
	text: string;
}

/**
 * Response from analysis API call.
 * Alias of AnalysisResult for compatibility with the API layer.
 */
export type AnalysisResponse = AnalysisResult;

/**
 * Response from recordings list API call.
 */
export interface RecordingsListResponse {
	recordings: Recording[];
	totalCount: number;
}

/**
 * Audio blob with metadata for local storage.
 */
export interface AudioBlobWithMetadata {
	blob: Blob;
	duration: number;
	createdAt: string;
	caseId?: string;
	title?: string;
}
