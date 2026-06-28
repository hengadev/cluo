/**
 * Client-side API wrappers for the Golang backend.
 *
 * ID chain problem: the backend uses separate IDs at each stage of the recording
 * pipeline (mediaId → jobId → transcriptionId → analysisId). Since the mobile
 * app routes use only the mediaId, we track the chain in localStorage.
 */

import type {
	UploadRecordingResponse,
	RecordingStatusResponse,
	TranscriptResponse,
	AnalysisResponse,
	RecordingsListResponse,
	ProcessingStep,
	Recording,
	Transcript,
	AnalysisResult,
	RecordingStatus,
} from "../types/recording";
import type { Case, CaseStatus } from "../types/case";
import { apiFetch, fetchWithRetry } from "./apiFetch";

const API_URL = import.meta.env.VITE_API_URL ?? "";

// ---------------------------------------------------------------------------
// localStorage ID chain
// ---------------------------------------------------------------------------

interface RecordingChain {
	jobId?: string;
	transcriptionId?: string;
	analysisId?: string;
}

const CHAIN_KEY = "cluo_recording_chain";
const CASE_KEY = "cluo_current_case_id";

function getChain(mediaId: string): RecordingChain {
	try {
		const raw = localStorage.getItem(CHAIN_KEY);
		if (!raw) return {};
		const all: Record<string, RecordingChain> = JSON.parse(raw);
		return all[mediaId] ?? {};
	} catch {
		return {};
	}
}

function updateChain(mediaId: string, update: Partial<RecordingChain>): void {
	try {
		const raw = localStorage.getItem(CHAIN_KEY);
		const all: Record<string, RecordingChain> = raw ? JSON.parse(raw) : {};
		all[mediaId] = { ...all[mediaId], ...update };
		localStorage.setItem(CHAIN_KEY, JSON.stringify(all));
	} catch {
		// localStorage unavailable (SSR or private browsing)
	}
}

function clearChain(mediaId: string): void {
	try {
		const raw = localStorage.getItem(CHAIN_KEY);
		if (!raw) return;
		const all: Record<string, RecordingChain> = JSON.parse(raw);
		delete all[mediaId];
		localStorage.setItem(CHAIN_KEY, JSON.stringify(all));
	} catch {}
}

// ---------------------------------------------------------------------------
// HTTP helper
// ---------------------------------------------------------------------------


// ---------------------------------------------------------------------------
// Backend response types (Go handler shapes)
// ---------------------------------------------------------------------------

interface MediaResponse {
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

interface SubmitJobApiResponse {
	jobId: string;
	mediaFileId: string;
	status: string;
	progress: number;
	createdAt: string;
}

interface JobStatusApiResponse {
	jobId: string;
	mediaFileId: string;
	status: string; // "pending" | "running" | "completed" | "failed" | "cancelled"
	progress: number;
	errorMessage?: string;
	transcriptionId?: string;
	createdAt: string;
	startedAt?: string;
	completedAt?: string;
}

interface TranscriptionApiResponse {
	id: string;
	jobId: string;
	mediaFileId: string;
	audioUrl: string;
	transcript: string;
	confidenceScore: number;
	language: string;
	duration: number; // milliseconds
	modelName: string;
	processingTimeMs: number;
	createdAt: string;
}

interface AnalysisApiResponse {
	id: string;
	transcriptionId: string;
	keyFindings: string;
	summary: string;
	sentiment: string;
	topics: string; // JSON-encoded array string from backend
	suggestedActions: string;
	modelUsed?: string;
	processingTimeMs?: number;
	createdAt: string;
}

interface ListMediaApiResponse {
	media: MediaResponse[];
	pagination: { page: number; pageSize: number; totalItems: number; totalPages: number };
}

interface CaseApiResponse {
	id: string;
	title: string;
	description: string;
	clientId: string;
	clientName?: string;
	clientNumber?: string;
	status: string;
	externalReference?: string;
	createdAt: string;
	updatedAt: string;
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

export function formatDate(iso: string): string {
	return new Date(iso).toLocaleDateString("fr-FR", {
		day: "2-digit",
		month: "short",
		year: "numeric",
	});
}

export function formatTime(iso: string): string {
	return new Date(iso).toLocaleTimeString("fr-FR", {
		hour: "2-digit",
		minute: "2-digit",
		hour12: false,
	});
}

function deriveStatus(isPublished: boolean, url: string): RecordingStatus {
	if (isPublished && url) return "completed";
	if (url) return "transcribing";
	return "uploading";
}

function mediaToRecording(m: MediaResponse): Recording & { audioUrl?: string } {
	return {
		id: m.id,
		caseId: m.caseId,
		title: m.caption || m.fileName,
		date: formatDate(m.createdAt),
		startTime: formatTime(m.createdAt),
		duration: 0,
		fileSize: m.fileSize,
		status: deriveStatus(m.isPublished, m.url),
		audioUrl: m.url || undefined,
	};
}

// ---------------------------------------------------------------------------
// API functions
// ---------------------------------------------------------------------------

/**
 * Upload a recording:
 * 1. POST /media  — stores the audio file
 * 2. POST /ai/speech/jobs — submits it for transcription
 * Stores the jobId in localStorage keyed by mediaId for later polling.
 */
export async function uploadRecording(
	blob: Blob,
	metadata?: { caseId?: string; title?: string },
): Promise<UploadRecordingResponse> {
	const caseId = metadata?.caseId ?? localStorage.getItem(CASE_KEY) ?? "";
	if (!caseId) {
		throw new Error("Le caseId est requis pour envoyer un enregistrement");
	}

	// Upload media
	const uploadForm = new FormData();
	uploadForm.append("file", blob, "recording.webm");
	uploadForm.append("caseId", caseId);
	if (metadata?.title) {
		uploadForm.append("caption", metadata.title);
	}

	const uploadRes = await fetchWithRetry(`${API_URL}/media`, {
		method: "POST",
		credentials: "include",
		body: uploadForm,
	});
	if (!uploadRes.ok) {
		const err = await uploadRes.text().catch(() => "Unknown error");
		throw new Error(`Échec de l'envoi : ${err}`);
	}
	const media: MediaResponse = await uploadRes.json();

	// Store the caseId for future listRecordings calls
	try {
		localStorage.setItem(CASE_KEY, caseId);
	} catch {}

	// Submit transcription job with the same audio blob
	const jobForm = new FormData();
	jobForm.append("file", blob, "recording.webm");
	jobForm.append("mediaFileId", media.id);

	const jobRes = await fetchWithRetry(`${API_URL}/ai/speech/jobs`, {
		method: "POST",
		credentials: "include",
		body: jobForm,
	});
	if (!jobRes.ok) {
		// Roll back the media upload so the queue entry can retry the full
		// flow without leaving an orphaned media record in the DB.
		await fetchWithRetry(`${API_URL}/media/${media.id}`, {
			method: "DELETE",
			credentials: "include",
		}).catch(() => {});
		const err = await jobRes.text().catch(() => "Unknown error");
		throw new Error(`Échec de la soumission de la tâche de transcription : ${err}`);
	}
	const job: SubmitJobApiResponse = await jobRes.json();

	updateChain(media.id, { jobId: job.jobId });

	return { id: media.id, status: "uploading" };
}

/**
 * Poll the transcription job status for a given mediaId.
 * Falls back to GET /media/{id} if no jobId is stored yet.
 */
export async function getRecordingStatus(id: string): Promise<RecordingStatusResponse> {
	const { jobId } = getChain(id);

	if (!jobId) {
		// No job ID yet — derive status from media record
		const media = await apiFetch<MediaResponse>(`/media/${id}`);
		const isComplete = !!(media.isPublished && media.url);
		const steps: ProcessingStep[] = [
			{ title: "Téléchargement audio", status: "completed" },
			{
				title: "Traitement de la transcription",
				status: isComplete ? "completed" : "processing",
			},
			{
				title: "Génération du résumé",
				status: isComplete ? "completed" : "processing",
			},
			{ title: "Terminé", status: isComplete ? "completed" : "processing" },
		];
		return {
			id,
			status: isComplete ? "completed" : "transcribing",
			processingSteps: steps,
		};
	}

	let job: JobStatusApiResponse;
	try {
		job = await apiFetch<JobStatusApiResponse>(`/ai/speech/jobs/${jobId}`);
	} catch (err) {
		// Job no longer exists (stale localStorage ID from a previous session or DB reset).
		// Clear the dead ID so we don't keep polling for it, then surface a failed state.
		if (err instanceof Error && err.message.includes("404")) {
			updateChain(id, { jobId: undefined });
			return {
				id,
				status: "failed",
				processingSteps: [
					{ title: "Téléchargement audio", status: "completed" },
					{ title: "Traitement de la transcription", status: "failed" },
					{ title: "Génération du résumé", status: "failed" },
					{ title: "Terminé", status: "failed" },
				],
				error: "L'enregistrement a été supprimé suite à une erreur de traitement",
			};
		}
		throw err;
	}

	if (job.transcriptionId) {
		updateChain(id, { transcriptionId: job.transcriptionId });
	}

	const statusMap: Record<string, RecordingStatus> = {
		pending: "transcribing",
		running: "transcribing",
		completed: "completed",
		failed: "failed",
		cancelled: "failed",
	};
	const recordingStatus: RecordingStatus = statusMap[job.status] ?? "transcribing";

	const isComplete = job.status === "completed";
	const isFailed = job.status === "failed" || job.status === "cancelled";

	const steps: ProcessingStep[] = [
		{ title: "Téléchargement audio", status: "completed" },
		{
			title: "Traitement de la transcription",
			status: isComplete ? "completed" : isFailed ? "failed" : "processing",
		},
		{
			title: "Génération du résumé",
			status: isComplete ? "completed" : isFailed ? "failed" : "pending",
		},
		{
			title: "Terminé",
			status: isComplete ? "completed" : isFailed ? "failed" : "pending",
		},
	];

	return {
		id,
		status: recordingStatus,
		processingSteps: steps,
		error: job.errorMessage,
	};
}

/**
 * Fetch the transcript for a recording by looking up the transcription via mediaFileId.
 * Stores the transcriptionId in the chain for later analysis calls.
 */
export async function getTranscript(id: string): Promise<TranscriptResponse> {
	const res = await apiFetch<{ transcriptions: TranscriptionApiResponse[]; total: number }>(
		`/ai/speech/transcriptions?mediaFileId=${id}`,
	);

	if (!res.transcriptions.length) {
		throw new Error("Transcription pas encore disponible");
	}

	const t = res.transcriptions[0];
	updateChain(id, { transcriptionId: t.id });

	return {
		recordingId: id,
		text: t.transcript,
		confidence: t.confidenceScore,
		isConfirmed: false, // backend has no confirm concept; set false so the user reviews before analyzing
		createdAt: t.createdAt,
		updatedAt: t.createdAt,
	};
}

/**
 * No-op: the backend has no transcript confirm endpoint.
 * The caller sets isConfirmed locally to gate the analyze action.
 */
export async function confirmTranscript(_id: string, _text: string): Promise<void> {}

/**
 * Request AI analysis of the transcript.
 * Looks up transcriptionId from the chain, then calls POST /ai/analysis/analyze.
 * Stores the resulting analysisId for getAnalysis.
 */
export async function analyzeTranscript(id: string): Promise<void> {
	let { transcriptionId } = getChain(id);

	if (!transcriptionId) {
		const res = await apiFetch<{ transcriptions: TranscriptionApiResponse[]; total: number }>(
			`/ai/speech/transcriptions?mediaFileId=${id}`,
		);
		if (!res.transcriptions.length) {
			throw new Error("Transcription non disponible — impossible d'analyser");
		}
		transcriptionId = res.transcriptions[0].id;
		updateChain(id, { transcriptionId });
	}

	const analysis = await apiFetch<AnalysisApiResponse>(`/ai/analysis/analyze`, {
		method: "POST",
		body: JSON.stringify({ transcriptionId }),
	});

	updateChain(id, { analysisId: analysis.id });
}

/**
 * Fetch the AI analysis for a recording.
 * Requires that analyzeTranscript has been called first (analysisId stored in chain).
 */
export async function getAnalysis(id: string): Promise<AnalysisResponse> {
	const { analysisId } = getChain(id);

	if (!analysisId) {
		throw new Error("Aucune analyse trouvée. Veuillez d'abord analyser la transcription.");
	}

	const a = await apiFetch<AnalysisApiResponse>(`/ai/analysis/${analysisId}`);

	return {
		id: a.id,
		transcriptionId: a.transcriptionId,
		keyFindings: a.keyFindings,
		summary: a.summary,
		sentiment: a.sentiment,
		topics: a.topics,
		suggestedActions: a.suggestedActions,
		createdAt: a.createdAt,
	};
}

/**
 * List recordings for a given caseId.
 * Falls back to the last-used caseId stored in localStorage.
 */
export async function listRecordings(options?: {
	limit?: number;
	offset?: number;
	caseId?: string;
	status?: string;
}): Promise<RecordingsListResponse> {
	const caseId = options?.caseId ?? localStorage.getItem(CASE_KEY) ?? "";

	if (!caseId) {
		return { recordings: [], totalCount: 0 };
	}

	const pageSize = options?.limit ?? 20;
	const page = options?.offset ? Math.floor(options.offset / pageSize) + 1 : 1;
	const params = new URLSearchParams({ page: String(page), pageSize: String(pageSize), type: "audio" });

	const res = await apiFetch<ListMediaApiResponse>(`/case/${caseId}/media?${params}`);

	const recordings = res.media.map(mediaToRecording) as Recording[];

	return {
		recordings,
		totalCount: res.pagination.totalItems,
	};
}

/**
 * Delete a media record and clear its local ID chain.
 */
export async function deleteRecording(id: string): Promise<void> {
	await apiFetch<void>(`/media/${id}`, { method: "DELETE" });
	clearChain(id);
}

/**
 * Fetch a recording with its transcript and analysis in one call.
 */
export async function getRecording(id: string): Promise<{
	recording: Recording & { audioUrl?: string };
	transcript: Transcript | null;
	analysis: AnalysisResult | null;
}> {
	const media = await apiFetch<MediaResponse>(`/media/${id}`);
	const recording = mediaToRecording(media);

	let transcript: Transcript | null = null;
	try {
		const res = await apiFetch<{ transcriptions: TranscriptionApiResponse[]; total: number }>(
			`/ai/speech/transcriptions?mediaFileId=${id}`,
		);
		if (res.transcriptions.length) {
			const t = res.transcriptions[0];
			updateChain(id, { transcriptionId: t.id });
			transcript = {
				recordingId: id,
				text: t.transcript,
				confidence: t.confidenceScore,
				isConfirmed: false,
				createdAt: t.createdAt,
				updatedAt: t.createdAt,
			};
		}
	} catch {
		// transcript not yet available
	}

	let analysis: AnalysisResult | null = null;
	try {
		const { analysisId } = getChain(id);
		if (analysisId) {
			const a = await apiFetch<AnalysisApiResponse>(`/ai/analysis/${analysisId}`);
			analysis = {
				id: a.id,
				transcriptionId: a.transcriptionId,
				keyFindings: a.keyFindings,
				summary: a.summary,
				sentiment: a.sentiment,
				topics: a.topics,
				suggestedActions: a.suggestedActions,
				createdAt: a.createdAt,
			};
		}
	} catch {
		// analysis not yet available
	}

	return { recording, transcript, analysis };
}

/**
 * Return the direct audio URL for a media record.
 */
export async function getAudioUrl(id: string): Promise<string> {
	const media = await apiFetch<MediaResponse>(`/media/${id}`);
	return media.url;
}

/**
 * Fetch all cases for the current user.
 */
export async function getCases(): Promise<Case[]> {
	const res = await apiFetch<{ cases: CaseApiResponse[] }>("/cases");
	return res.cases.map((c) => ({
		id: c.id,
		title: c.title,
		status: c.status as CaseStatus,
		externalReference: c.externalReference,
		clientId: c.clientId,
		clientName: c.clientName,
		clientNumber: c.clientNumber,
	}));
}

/**
 * Fetch the active case by ID.
 * Falls back to the last-used caseId stored in localStorage.
 */
export async function getCurrentCase(caseId?: string): Promise<Case | null> {
	let id = caseId;
	if (!id) {
		try {
			id = localStorage.getItem(CASE_KEY) ?? "";
		} catch {
			return null;
		}
	}
	if (!id) return null;

	try {
		const res = await apiFetch<CaseApiResponse>(`/case/${id}`);
		return {
			id: res.id,
			title: res.title,
			status: res.status as CaseStatus,
			externalReference: res.externalReference,
			clientId: res.clientId,
			clientName: res.clientName,
			clientNumber: res.clientNumber,
		};
	} catch {
		return null;
	}
}

/**
 * Poll until the recording reaches one of the target statuses.
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
		if (targetStatus.includes(status.status as "completed" | "failed")) {
			return status;
		}
		await new Promise((resolve) => setTimeout(resolve, interval));
	}
	throw new Error("Délai d'attente dépassé : l'enregistrement n'a pas été traité dans le temps imparti");
}
