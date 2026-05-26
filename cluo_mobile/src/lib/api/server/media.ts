/**
 * Server-side API wrappers for media endpoints.
 * These functions use SvelteKit's fetch helper and run only on the server.
 */

import { error } from "@sveltejs/kit";
import type { RecordingStatus } from "$lib/types/recording";

const API_URL = import.meta.env.VITE_API_URL ?? "";

// ---------- Types ----------

/** Response shape from GET /media/{id} */
export interface MediaApiResponse {
	id: string;
	caseId: string;
	url: string;
	type: string;
	mimeType: string;
	fileName: string;
	fileSize: number;
	caption: string;
	isPublished: boolean;
	createdAt: string;
}

/** Mapped recording shape for the recording detail page */
export interface RecordingPageData {
	id: string;
	title: string;
	date: string; // "DD Mon, YYYY"
	startTime: string;
	duration: string; // "MM:SS"
	audioUrl: string;
	status: RecordingStatus;
}

/** Processing step for the processing page */
export interface ProcessingStepData {
	title: string;
	status: "pending" | "processing" | "completed" | "failed";
}

/** Result for the processing page */
export interface ProcessingPageData {
	recordingId: string;
	steps: ProcessingStepData[];
	error: string | null;
}

/** Upload response from POST /media */
export interface UploadMediaResponse {
	id: string;
	caseId: string;
	url: string;
	type: string;
	mimeType: string;
	fileName: string;
	fileSize: number;
	caption: string;
	isPublished: boolean;
	createdAt: string;
}

// ---------- Helpers ----------

/**
 * Format a Date or ISO string to "DD Mon, YYYY" (e.g. "01 Jan, 2025").
 */
function formatDate(date: Date | string): string {
	const d = typeof date === "string" ? new Date(date) : date;
	return d.toLocaleDateString("en-GB", {
		day: "2-digit",
		month: "short",
		year: "numeric",
	});
}

/**
 * Format milliseconds to "MM:SS".
 */
function formatDuration(ms: number): string {
	const totalSeconds = Math.floor(ms / 1000);
	const mins = Math.floor(totalSeconds / 60);
	const secs = totalSeconds % 60;
	return `${mins.toString().padStart(2, "0")}:${secs.toString().padStart(2, "0")}`;
}

function deriveStatus(media: MediaApiResponse): RecordingStatus {
	if (media.isPublished && media.url) return "completed";
	if (media.url) return "transcribing";
	return "uploading";
}

// ---------- API Wrappers ----------

/**
 * Fetch a media record by ID from the Go backend.
 * Uses SvelteKit's fetch helper (passed from the load function).
 *
 * @throws {import('@sveltejs/kit').HttpError} 404 if media not found, 500 for other failures.
 */
export async function fetchMediaById(
	fetch: typeof globalThis.fetch,
	id: string,
): Promise<RecordingPageData> {
	const res = await fetch(`${API_URL}/media/${id}`, {
		headers: { Accept: "application/json" },
	});

	if (res.status === 404) {
		throw error(404, { message: "Recording not found" });
	}

	if (!res.ok) {
		const body = await res.text().catch(() => "Unknown error");
		console.error(`API error fetching media ${id}: ${res.status} ${body}`);
		throw error(500, { message: "Failed to load recording" });
	}

	const media: MediaApiResponse = await res.json();

	// Map API response → page shape
	const createdAt = new Date(media.createdAt);

	return {
		id: media.id,
		title: media.caption || media.fileName,
		date: formatDate(createdAt),
		startTime: createdAt.toLocaleTimeString("en-GB", {
			hour: "2-digit",
			minute: "2-digit",
		}),
		duration: "00:00", // Duration not available in media response; will be enriched later
		audioUrl: media.url || "",
		status: deriveStatus(media),
	};
}

/**
 * Fetch processing status for a media record.
 * Maps the media status to the four-step processing array the page expects.
 *
 * @throws {import('@sveltejs/kit').HttpError} 404 if media not found, 500 for other failures.
 */
export async function fetchProcessingStatus(
	fetch: typeof globalThis.fetch,
	id: string,
): Promise<ProcessingPageData> {
	const res = await fetch(`${API_URL}/media/${id}`, {
		headers: { Accept: "application/json" },
	});

	if (res.status === 404) {
		throw error(404, { message: "Recording not found" });
	}

	if (!res.ok) {
		const body = await res.text().catch(() => "Unknown error");
		console.error(`API error fetching media ${id}: ${res.status} ${body}`);
		throw error(500, { message: "Failed to load processing status" });
	}

	const media: MediaApiResponse = await res.json();

	// Map to four-step array
	// Step 1: Téléchargement audio — always completed (we got here, so upload is done)
	// Step 2: Traitement de la transcription — processing or completed
	// Step 3: Génération du résumé — processing or completed
	// Step 4: Terminé — processing or completed
	const steps: ProcessingStepData[] = [
		{
			title: "Téléchargement audio",
			status: "completed",
		},
		{
			title: "Traitement de la transcription",
			status: "processing",
		},
		{
			title: "Génération du résumé",
			status: "processing",
		},
		{
			title: "Terminé",
			status: "processing",
		},
	];

	// If media has a URL and is published, transcription might be complete
	// This is a heuristic — the actual transcription status comes from the
	// transcription job API, but for the initial wiring we map from media state.
	if (media.url && media.isPublished) {
		steps[1].status = "completed";
		steps[2].status = "completed";
		steps[3].status = "completed";
	}

	return {
		recordingId: id,
		steps,
		error: null,
	};
}

/**
 * Upload an audio blob to the Go backend via multipart relay.
 *
 * @param fetch - SvelteKit fetch helper
 * @param blob - Audio blob
 * @param caseId - Case ID to associate the recording with
 * @returns The created media ID
 * @throws {Error} on API failure
 */
export async function uploadAudio(
	fetch: typeof globalThis.fetch,
	blob: Blob,
	caseId: string,
): Promise<string> {
	const formData = new FormData();
	formData.append("file", blob, "recording.webm");
	formData.append("caseId", caseId);

	const res = await fetch(`${API_URL}/media`, {
		method: "POST",
		body: formData,
		// Do NOT set Content-Type — fetch sets multipart boundaries automatically
	});

	if (!res.ok) {
		const body = await res.text().catch(() => "Unknown error");
		throw new Error(`Upload failed (${res.status}): ${body}`);
	}

	const created: UploadMediaResponse = await res.json();
	return created.id;
}
