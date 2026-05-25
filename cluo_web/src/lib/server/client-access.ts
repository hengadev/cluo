import { env } from '$env/dynamic/private';
import sanitizeHtml from 'sanitize-html';

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

export type TokenValidation =
	| { valid: true; caseData: CaseResponse }
	| { valid: false; reason: 'expired' | 'invalid' | 'unavailable' };

/** Shape returned by GET /token/{token} on success (PortalCaseResponse from Go API). */
export type CaseResponse = {
	id: string;
	title: string;
	description: string;
	clientId: string;
	assignedContactID?: string;
	caseSubjectId?: string;
	externalReference?: string;
	caseTypeId?: string;
	status: string;
	placename?: string;
	address1?: string;
	address2?: string;
	city?: string;
	postalCode?: string;
	country?: string;
	latitude?: string;
	longitude?: string;
	locationType?: string;
	locationNotes?: string;
	createdAt: string;
	updatedAt: string;
	tokenExpiresAt: string;
};

/** Shape returned by GET /token/{token}/documents on success. */
export type DocumentSummaryResponse = {
	id: string;
	case_id: string;
	client_id: string;
	type: 'estimate' | 'mandate' | 'contract' | 'invoice' | 'report' | 'other';
	status: 'draft' | 'sent' | 'signed' | 'active' | 'archived' | 'cancelled' | 'rejected' | 'expired';
	document_ref: string;
	created_at: string;
	updated_at: string;
};

/** Shape returned by GET /token/{token}/media on success. */
export type MediaResponse = {
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
};

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

const USE_MOCK_DATA = env.USE_MOCK_DATA === 'true';
const API_BASE = env.API_BASE_URL ?? 'http://localhost:8080';

// getApiBaseUrl returns the configured API base URL.
// Exported for use in server routes that need to proxy to the Go API.
export function getApiBaseUrl(): string {
	return API_BASE;
}

async function apiFetch(path: string): Promise<Response> {
	return fetch(`${API_BASE}${path}`);
}

// ---------------------------------------------------------------------------
// Mock data (development only)
// ---------------------------------------------------------------------------

const MOCK_CASE: CaseResponse = {
	id: 'mock-case-001',
	title: 'Succession — Famille Martin',
	description: 'Tous les documents ont été vérifiés et compilés pour votre dossier.',
	clientId: 'mock-client-001',
	status: 'released',
	createdAt: '2026-04-12T00:00:00Z',
	updatedAt: '2026-05-01T00:00:00Z',
	tokenExpiresAt: '2026-06-12T00:00:00Z'
};

const MOCK_DOCUMENTS: DocumentSummaryResponse[] = [
	{
		id: 'mock-doc-001',
		case_id: 'mock-case-001',
		client_id: 'mock-client-001',
		type: 'estimate',
		status: 'signed',
		document_ref: 'DEV-2026-001',
		created_at: '2026-04-10T00:00:00Z',
		updated_at: '2026-04-12T00:00:00Z'
	},
	{
		id: 'mock-doc-002',
		case_id: 'mock-case-001',
		client_id: 'mock-client-001',
		type: 'mandate',
		status: 'active',
		document_ref: 'MAN-2026-001',
		created_at: '2026-04-10T00:00:00Z',
		updated_at: '2026-04-12T00:00:00Z'
	},
	{
		id: 'mock-doc-003',
		case_id: 'mock-case-001',
		client_id: 'mock-client-001',
		type: 'contract',
		status: 'active',
		document_ref: 'CTR-2026-001',
		created_at: '2026-04-10T00:00:00Z',
		updated_at: '2026-04-12T00:00:00Z'
	},
	{
		id: 'mock-doc-004',
		case_id: 'mock-case-001',
		client_id: 'mock-client-001',
		type: 'invoice',
		status: 'sent',
		document_ref: 'FAC-2026-001',
		created_at: '2026-05-01T00:00:00Z',
		updated_at: '2026-05-01T00:00:00Z'
	}
];

const MOCK_MEDIA: MediaResponse[] = [
	{
		id: 'mock-media-001',
		caseId: 'mock-case-001',
		url: 'https://images.unsplash.com/photo-1506905925346-21bda4d32df4?w=1200',
		type: 'image',
		mimeType: 'image/jpeg',
		fileName: 'surveillance-001.jpg',
		fileSize: 2_400_000,
		caption: 'Photo de surveillance — entrée principale',
		isPublished: true,
		createdAt: '2026-04-15T00:00:00Z'
	},
	{
		id: 'mock-media-002',
		caseId: 'mock-case-001',
		url: 'https://images.unsplash.com/photo-1470071459604-3b5ec3a7fe05?w=1200',
		type: 'image',
		mimeType: 'image/jpeg',
		fileName: 'surveillance-002.jpg',
		fileSize: 3_100_000,
		caption: 'Photo de surveillance — vue générale',
		isPublished: true,
		createdAt: '2026-04-15T01:00:00Z'
	},
	{
		id: 'mock-media-003',
		caseId: 'mock-case-001',
		url: 'https://images.unsplash.com/photo-1441974231531-c6227db76b6e?w=1200',
		type: 'image',
		mimeType: 'image/jpeg',
		fileName: 'surveillance-003.jpg',
		fileSize: 1_800_000,
		caption: 'Photo de surveillance — zone arrière',
		isPublished: true,
		createdAt: '2026-04-16T00:00:00Z'
	},
	{
		id: 'mock-media-004',
		caseId: 'mock-case-001',
		url: 'https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4',
		type: 'video',
		mimeType: 'video/mp4',
		fileName: 'video-surveillance.mp4',
		fileSize: 15_000_000,
		caption: 'Vidéo de surveillance — entrée',
		isPublished: true,
		createdAt: '2026-04-17T00:00:00Z'
	},
	{
		id: 'mock-media-005',
		caseId: 'mock-case-001',
		url: 'https://www.soundhelix.com/examples/mp3/SoundHelix-Song-1.mp3',
		type: 'audio',
		mimeType: 'audio/mpeg',
		fileName: 'enregistrement-audio.mp3',
		fileSize: 5_200_000,
		caption: 'Enregistrement audio — témoin A',
		isPublished: true,
		createdAt: '2026-04-18T00:00:00Z'
	}
];

// ---------------------------------------------------------------------------
// Public API
// ---------------------------------------------------------------------------

/**
 * Validate a client-access token by calling GET /token/{token}.
 * On success the API returns the full CaseResponse.
 * On failure we classify the error as expired / invalid / unavailable.
 */
export async function validateClientToken(token: string): Promise<TokenValidation> {
	if (USE_MOCK_DATA) {
		if (token === 'expired') return { valid: false, reason: 'expired' };
		if (token === 'unavailable') return { valid: false, reason: 'unavailable' };
		if (token === 'invalid') return { valid: false, reason: 'invalid' };
		return { valid: true, caseData: MOCK_CASE };
	}

	try {
		const res = await apiFetch(`/token/${encodeURIComponent(token)}`);

		if (res.ok) {
			const caseData: CaseResponse = await res.json();
			return { valid: true, caseData };
		}

		// Map HTTP status → error reason
		if (res.status === 401) {
			return { valid: false, reason: 'expired' };
		}
		if (res.status === 404) {
			return { valid: false, reason: 'invalid' };
		}
		return { valid: false, reason: 'unavailable' };
	} catch {
		// Network / server unreachable
		return { valid: false, reason: 'unavailable' };
	}
}

/**
 * Check whether a token has published media by calling GET /token/{token}/media.
 * Returns the media array, or null when the call fails (fail-open: show the tab).
 */
export type DocumentsResult =
	| { status: 'ok'; documents: DocumentSummaryResponse[] }
	| { status: 'error' };

/** Statuses that are visible to the client — drafts are excluded. */
const VISIBLE_STATUSES: ReadonlySet<string> = new Set(['sent', 'signed', 'active', 'archived']);

/**
 * Fetch document summaries for a given token.
 * Filters out draft documents on the client side.
 * Returns the filtered list, or an error status when the call fails.
 */
export async function getDocumentsByToken(token: string): Promise<DocumentsResult> {
	if (USE_MOCK_DATA) {
		const filtered = MOCK_DOCUMENTS.filter((d) => VISIBLE_STATUSES.has(d.status));
		return { status: 'ok', documents: filtered };
	}

	try {
		const res = await apiFetch(`/token/${encodeURIComponent(token)}/documents`);
		if (!res.ok) return { status: 'error' };
		const all: DocumentSummaryResponse[] = await res.json();
		const filtered = all.filter((d) => VISIBLE_STATUSES.has(d.status));
		return { status: 'ok', documents: filtered };
	} catch {
		return { status: 'error' };
	}
}

export async function getTokenMedia(token: string): Promise<MediaResponse[] | null> {
	if (USE_MOCK_DATA) {
		// Simulate "no media" with the token value "no-media"
		if (token === 'no-media') return [];
		return MOCK_MEDIA;
	}

	try {
		const res = await apiFetch(`/token/${encodeURIComponent(token)}/media`);
		if (!res.ok) return null; // fail open — caller should still show Médias tab
		return await res.json();
	} catch {
		return null; // fail open
	}
}

export type ReportHtmlResult =
	| { status: 'ok'; html: string }
	| { status: 'not_found' }
	| { status: 'error' };

const ALLOWED_REPORT_TAGS = [
	'h1', 'h2', 'h3',
	'p', 'strong', 'em', 'u',
	'ul', 'ol', 'li',
	'blockquote',
];

function sanitizeReportHtml(raw: string): string {
	return sanitizeHtml(raw, {
		allowedTags: ALLOWED_REPORT_TAGS,
		allowedAttributes: {},
	});
}

/**
 * Fetch the Rapport rendered as HTML for a given token.
 * Sanitizes the HTML before returning it.
 * Distinguishes "no rapport" (not_found) from server errors (error).
 */
export async function getReportHtml(token: string): Promise<ReportHtmlResult> {
	if (USE_MOCK_DATA) {
		if (token === 'no-rapport') return { status: 'not_found' };
		const raw = `<h1>Rapport d'investigation</h1>
<p>Le <strong>suspect</strong> a été observé à plusieurs reprises dans le quartier de la Gare du Nord.</p>
<h2>Observations</h2>
<blockquote><p>Note importante : toutes les observations ont été réalisées dans le respect de la légalité.</p></blockquote>
<ul><li><p>Point de surveillance A</p></li><li><p>Point de surveillance B</p></li></ul>
<h2>Conclusion</h2>
<p><em>Fin du rapport de surveillance.</em></p>`;
		return { status: 'ok', html: sanitizeReportHtml(raw) };
	}

	try {
		const res = await apiFetch(`/token/${encodeURIComponent(token)}/report/html`);
		if (res.status === 404) return { status: 'not_found' };
		if (!res.ok) return { status: 'error' };
		const raw = await res.text();
		return { status: 'ok', html: sanitizeReportHtml(raw) };
	} catch {
		return { status: 'error' };
	}
}

export async function streamMediaFile(
	token: string,
	mediaId: string
): Promise<{ stream: ReadableStream; fileName: string; mimeType: string; fileSize: number } | null> {
	if (USE_MOCK_DATA) {
		const media = MOCK_MEDIA.find((m) => m.id === mediaId);
		if (!media) return null;
		const res = await fetch(media.url);
		if (!res.ok || !res.body) return null;
		return { stream: res.body, fileName: media.fileName, mimeType: media.mimeType, fileSize: media.fileSize };
	}

	try {
		const metaRes = await apiFetch(`/token/${encodeURIComponent(token)}/media/${encodeURIComponent(mediaId)}`);
		if (!metaRes.ok) return null;
		const media: MediaResponse = await metaRes.json();

		const fileRes = await fetch(media.url);
		if (!fileRes.ok || !fileRes.body) return null;

		return {
			stream: fileRes.body,
			fileName: media.fileName,
			mimeType: media.mimeType,
			fileSize: media.fileSize
		};
	} catch {
		return null;
	}
}

export async function streamCaseArchive(token: string): Promise<ReadableStream> {
	if (USE_MOCK_DATA) {
		const content = `Mock archive — token: ${token}\nDevelopment mode only.`;
		const bytes = new TextEncoder().encode(content);
		return new ReadableStream({
			start(controller) {
				controller.enqueue(bytes);
				controller.close();
			}
		});
	}

	const res = await apiFetch(`/token/${encodeURIComponent(token)}/archive`);
	if (!res.ok || !res.body) {
		throw new Error(`Failed to download archive: ${res.status}`);
	}
	return res.body;
}
