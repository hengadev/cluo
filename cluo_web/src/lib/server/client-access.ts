import { env } from '$env/dynamic/private';

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

const MOCK_MEDIA: MediaResponse[] = [
	{
		id: 'mock-media-001',
		caseId: 'mock-case-001',
		url: '/mock/photo1.jpg',
		type: 'image',
		mimeType: 'image/jpeg',
		fileName: 'photo1.jpg',
		fileSize: 2_400_000,
		caption: 'Photo de surveillance',
		isPublished: true,
		createdAt: '2026-04-15T00:00:00Z'
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

/**
 * Stream the case archive (zip) for download.
 * The token is used directly — no cookie session.
 */
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
