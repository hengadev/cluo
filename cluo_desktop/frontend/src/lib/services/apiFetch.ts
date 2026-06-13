import { API_BASE_URL } from '../config';

let isRefreshing = false;
let refreshPromise: Promise<boolean> | null = null;

interface ApiFetchOptions extends RequestInit {
	// Allow skipping refresh for the refresh endpoint itself
	skipRefresh?: boolean;
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
		response = await fetch(url, options);
	} catch {
		// WebKitGTK may fail on stale connections after idle — retry once on a fresh connection
		try {
			response = await fetch(url, options);
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
				response = await fetch(url, options);
			} catch {
				try {
					response = await fetch(url, options);
				} catch {
					throw new Error('Impossible de joindre le serveur. Vérifiez votre connexion réseau.');
				}
			}
		} else {
			throw new Error('Session expirée. Veuillez vous reconnecter.');
		}
	}

	return response;
}

/**
 * Attempt to refresh the access token using the refresh token cookie.
 */
async function attemptRefresh(): Promise<boolean> {
	try {
		const response = await fetch(`${API_BASE_URL}/auth/refresh`, {
			method: 'POST',
			credentials: 'include',
		});

		return response.ok;
	} catch {
		return false;
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
