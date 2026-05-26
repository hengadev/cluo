/**
 * API utilities for audio recording, transcription, and analysis.
 * All functions communicate with the Golang backend.
 */

import type {
	UploadRecordingResponse,
	RecordingStatusResponse,
	TranscriptResponse,
	AnalysisResponse,
	RecordingsListResponse,
} from "../schemas/recording";
import type {
	ConfirmTranscriptRequest,
	ProcessingStep,
	RecordingStatusResponse as RecordingStatusInterface,
} from "../types/recording";

/**
 * Base URL for the backend API.
 * Uses environment variable or defaults to local development.
 */
const API_URL = import.meta.env.VITE_API_URL ?? "";

/**
 * Helper function to make API requests with proper error handling.
 */
async function apiFetch<T>(
	path: string,
	options?: RequestInit,
): Promise<T> {
	const url = `${API_URL}${path}`;

	const response = await fetch(url, {
		credentials: "include",
		headers: {
			"Content-Type": "application/json",
			...options?.headers,
		},
		...options,
	});

	if (!response.ok) {
		const errorText = await response.text().catch(() => "Unknown error");
		throw new Error(`API error (${response.status}): ${errorText}`);
	}

	return response.json() as Promise<T>;
}

/**
 * Upload an audio recording to the backend.
 *
 * @param blob - The audio blob to upload
 * @param metadata - Optional metadata (caseId, title, etc.)
 * @returns Response with recording ID and initial status
 */
export async function uploadRecording(
	blob: Blob,
	metadata?: { caseId?: string; title?: string },
): Promise<UploadRecordingResponse> {
	const formData = new FormData();
	formData.append("audio", blob, "recording.webm");

	if (metadata?.caseId) {
		formData.append("caseId", metadata.caseId);
	}
	if (metadata?.title) {
		formData.append("title", metadata.title);
	}

	const response = await fetch(`${API_URL}/api/audio`, {
		method: "POST",
		credentials: "include",
		body: formData,
	});

	if (!response.ok) {
		const errorText = await response.text().catch(() => "Unknown error");
		throw new Error(`Failed to upload recording: ${errorText}`);
	}

	return response.json() as Promise<UploadRecordingResponse>;
}

/**
 * Get the current processing status of a recording.
 * Fetches GET /media/{id} from the Go API and maps the media state to processing steps.
 *
 * @param id - The recording ID
 * @returns Recording status with processing steps
 */
export async function getRecordingStatus(
	id: string,
): Promise<RecordingStatusInterface> {
	const media = await apiFetch<{ url: string; isPublished: boolean }>(`/media/${id}`);
	const isComplete = !!(media.isPublished && media.url);

	const steps: ProcessingStep[] = [
		{ title: "Téléchargement audio", status: "completed" },
		{ title: "Traitement de la transcription", status: isComplete ? "completed" : "processing" },
		{ title: "Génération du résumé", status: isComplete ? "completed" : "processing" },
		{ title: "Terminé", status: isComplete ? "completed" : "processing" },
	];

	return {
		id,
		status: isComplete ? "completed" : "transcribing",
		processingSteps: steps,
	};
}

/**
 * Get the transcript for a recording.
 *
 * @param id - The recording ID
 * @returns Transcript data with text and metadata
 */
export async function getTranscript(
	id: string,
): Promise<TranscriptResponse> {
	return apiFetch<TranscriptResponse>(`/api/recordings/${id}/transcript`);
}

/**
 * Confirm or update the transcript for a recording.
 *
 * @param id - The recording ID
 * @param text - The confirmed/edited transcript text
 * @returns void on success
 */
export async function confirmTranscript(
	id: string,
	text: string,
): Promise<void> {
	await apiFetch<void>(`/api/recordings/${id}/transcript`, {
		method: "PUT",
		body: JSON.stringify({ text }),
	});
}

/**
 * Request analysis of a confirmed transcript.
 *
 * @param id - The recording ID
 * @returns void on success (analysis is async)
 */
export async function analyzeTranscript(
	id: string,
): Promise<void> {
	await apiFetch<void>(`/api/recordings/${id}/analyze`, {
		method: "POST",
	});
}

/**
 * Get the analysis results for a recording.
 *
 * @param id - The recording ID
 * @returns Analysis results with categorized suggestions
 */
export async function getAnalysis(
	id: string,
): Promise<AnalysisResponse> {
	return apiFetch<AnalysisResponse>(`/api/recordings/${id}/analysis`);
}

/**
 * List all recordings for the current user.
 *
 * @param options - Optional query parameters (limit, offset, caseId, status)
 * @returns List of recordings with total count
 */
export async function listRecordings(
	options?: {
		limit?: number;
		offset?: number;
		caseId?: string;
		status?: string;
	},
): Promise<RecordingsListResponse> {
	const params = new URLSearchParams();
	if (options?.limit) params.append("limit", options.limit.toString());
	if (options?.offset) params.append("offset", options.offset.toString());
	if (options?.caseId) params.append("caseId", options.caseId);
	if (options?.status) params.append("status", options.status);

	const queryString = params.toString();
	const path = `/api/recordings${queryString ? `?${queryString}` : ""}`;

	return apiFetch<RecordingsListResponse>(path);
}

/**
 * Delete a recording.
 *
 * @param id - The recording ID
 * @returns void on success
 */
export async function deleteRecording(
	id: string,
): Promise<void> {
	await apiFetch<void>(`/api/recordings/${id}`, {
		method: "DELETE",
	});
}

/**
 * Poll for recording status updates.
 * Returns a promise that resolves when the recording reaches the target status.
 *
 * @param id - The recording ID
 * @param targetStatus - The status to wait for (default: "completed")
 * @param interval - Polling interval in ms (default: 2000)
 * @param timeout - Maximum time to wait in ms (default: 300000 = 5 minutes)
 * @returns Final recording status
 */
export async function pollRecordingStatus(
	id: string,
	targetStatus: Array<"completed" | "failed"> = ["completed", "failed"],
	interval = 2000,
	timeout = 300000,
): Promise<RecordingStatusResponse> {
	const startTime = Date.now();

	while (Date.now() - startTime < timeout) {
		const status = await getRecordingStatus(id);

		if (targetStatus.includes(status.status)) {
			return status;
		}

		// Wait before next poll
		await new Promise((resolve) => setTimeout(resolve, interval));
	}

	throw new Error("Polling timeout: recording did not complete in time");
}

/**
 * Get a single recording with its transcript and analysis.
 *
 * @param id - The recording ID
 * @returns Recording details with transcript and analysis if available
 */
export async function getRecording(id: string): Promise<{
	recording: import("../types/recording").Recording & { audioUrl?: string };
	transcript: import("../types/recording").Transcript | null;
	analysis: import("../types/recording").AnalysisResult | null;
}> {
	// Fetch recording details
	const recording = await apiFetch<import("../types/recording").Recording & { audioUrl?: string }>(
		`/api/recordings/${id}`,
	);

	// Try to fetch transcript (may not exist)
	let transcript: import("../types/recording").Transcript | null = null;
	try {
		transcript = await getTranscript(id);
	} catch {
		// Transcript not available yet
	}

	// Try to fetch analysis (may not exist)
	let analysis: import("../types/recording").AnalysisResult | null = null;
	try {
		analysis = await getAnalysis(id);
	} catch {
		// Analysis not available yet
	}

	return { recording, transcript, analysis };
}

/**
 * Get the audio URL/blob for a recording.
 *
 * @param id - The recording ID
 * @returns URL to the audio file
 */
export async function getAudioUrl(id: string): Promise<string> {
	// The backend should return a signed URL or the audio can be streamed directly
	return `${API_URL}/api/recordings/${id}/audio`;
}
