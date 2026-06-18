import { goto } from '$app/navigation';
import { browser } from '$app/environment';
import { API_BASE_URL } from '../config';

let isRefreshing = false;
let refreshPromise: Promise<boolean> | null = null;

// WebKitGTK can silently hang forever on a stale keep-alive connection instead
// of throwing — without this, a hung request leaves callers stuck loading.
const REQUEST_TIMEOUT_MS = 15_000;

interface ApiFetchOptions extends RequestInit {
	// Allow skipping refresh for the refresh endpoint itself
	skipRefresh?: boolean;
}

function fetchWithTimeout(url: RequestInfo | URL, options: RequestInit): Promise<Response> {
	const controller = new AbortController();

	// Racing against a timer (rather than relying solely on the abort causing
	// fetch() to reject) is required: this WebKitGTK build doesn't reliably
	// reject an in-flight fetch when its AbortSignal fires, so a stale
	// keep-alive connection can leave the fetch promise pending forever even
	// after abort() is called. The race forces the await to settle regardless.
	let timeoutId: ReturnType<typeof setTimeout>;
	const timeoutPromise = new Promise<never>((_, reject) => {
		timeoutId = setTimeout(() => {
			controller.abort();
			reject(new Error('Request timed out'));
		}, REQUEST_TIMEOUT_MS);
	});

	return Promise.race([
		fetch(url, { ...options, signal: controller.signal }),
		timeoutPromise,
	]).finally(() => clearTimeout(timeoutId));
}

/**
 * Wrapper around fetch that handles authentication and token refresh.
 * Includes credentials: 'include' for httpOnly cookies.
 * On 401, attempts one refresh before failing.
 */
export async function apiFetch(
	input: RequestInfo | URL,
	init: ApiFetchOptions = {}
): Promise<Response> {
	const url = typeof input === 'string' && input.startsWith('/')
		? `${API_BASE_URL}${input}`
		: input;

	const options: RequestInit = {
		...init,
		credentials: 'include',
		headers: {
			'Content-Type': 'application/json',
			...init.headers,
		},
	};

	let response: Response;
	try {
		response = await fetchWithTimeout(url, options);
	} catch {
		// WebKitGTK may fail or hang on stale connections after idle — retry once on a fresh connection
		try {
			response = await fetchWithTimeout(url, options);
		} catch {
			throw new Error('Impossible de joindre le serveur. Vérifiez votre connexion réseau.');
		}
	}

	// If we get a 401 and we're not already refreshing, try to refresh
	if (response.status === 401 && !init.skipRefresh) {
		if (!isRefreshing) {
			isRefreshing = true;
			refreshPromise = attemptRefresh();
		}

		const refreshed = await refreshPromise;
		isRefreshing = false;
		refreshPromise = null;

		if (refreshed) {
			try {
				response = await fetchWithTimeout(url, options);
			} catch {
				try {
					response = await fetchWithTimeout(url, options);
				} catch {
					throw new Error('Impossible de joindre le serveur. Vérifiez votre connexion réseau.');
				}
			}
		} else {
			if (browser) {
				goto('/login');
			}
			throw new Error('Session expirée. Veuillez vous reconnecter.');
		}
	}

	return response;
}

async function attemptRefresh(): Promise<boolean> {
	const opts: RequestInit = { method: 'POST', credentials: 'include' };
	try {
		const response = await fetchWithTimeout(`${API_BASE_URL}/auth/refresh`, opts);
		return response.ok;
	} catch {
		// WebKitGTK stale-connection — retry once on a fresh connection
		try {
			const response = await fetchWithTimeout(`${API_BASE_URL}/auth/refresh`, opts);
			return response.ok;
		} catch {
			return false;
		}
	}
}

/**
 * Type-safe API request helpers
 */
export async function apiGet<T>(url: string, options?: Omit<ApiFetchOptions, 'method'>): Promise<T> {
	const response = await apiFetch(url, { ...options, method: 'GET' });
	if (!response.ok) {
		throw new Error(`API error: ${response.status} ${response.statusText}`);
	}
	return response.json();
}

export async function apiPost<T>(url: string, data?: unknown, options?: Omit<ApiFetchOptions, 'method' | 'body'>): Promise<T> {
	const response = await apiFetch(url, {
		...options,
		method: 'POST',
		body: JSON.stringify(data),
	});
	if (!response.ok) {
		throw new Error(`API error: ${response.status} ${response.statusText}`);
	}
	return response.json();
}

export async function apiDelete<T>(url: string, options?: Omit<ApiFetchOptions, 'method'>): Promise<T> {
	const response = await apiFetch(url, { ...options, method: 'DELETE' });
	if (!response.ok) {
		throw new Error(`API error: ${response.status} ${response.statusText}`);
	}
	return response.json();
}
