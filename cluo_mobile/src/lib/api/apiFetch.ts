const API_URL = import.meta.env.VITE_API_URL ?? "";

let isRefreshing = false;
let refreshPromise: Promise<boolean> | null = null;

// Retries once on network error — guards against stale connections after idle
export async function fetchWithRetry(
	input: RequestInfo | URL,
	init?: RequestInit,
): Promise<Response> {
	try {
		return await fetch(input, init);
	} catch {
		return await fetch(input, init);
	}
}

async function attemptRefresh(): Promise<boolean> {
	try {
		const res = await fetch(`${API_URL}/auth/refresh`, {
			method: "POST",
			credentials: "include",
		});
		return res.ok;
	} catch {
		return false;
	}
}

// Raw helper: prepends API_URL, sets JSON headers, returns Response as-is.
// On 401, attempts one token refresh then retries before returning.
export async function apiFetchRaw(
	path: string,
	options?: RequestInit & { skipRefresh?: boolean },
): Promise<Response> {
	const url = `${API_URL}${path}`;
	const init: RequestInit = {
		credentials: "include",
		...options,
		headers: {
			"Content-Type": "application/json",
			...options?.headers,
		},
	};

	let response = await fetchWithRetry(url, init);

	if (response.status === 401 && !options?.skipRefresh) {
		if (!isRefreshing) {
			isRefreshing = true;
			refreshPromise = attemptRefresh();
		}
		const refreshed = await refreshPromise;
		isRefreshing = false;
		refreshPromise = null;

		if (refreshed) {
			response = await fetchWithRetry(url, init);
		} else {
			throw new Error("Session expirée. Veuillez vous reconnecter.");
		}
	}

	return response;
}

// JSON helper: prepends API_URL, sets JSON headers, throws on non-ok
export async function apiFetch<T>(
	path: string,
	options?: RequestInit,
): Promise<T> {
	const response = await apiFetchRaw(path, options);
	if (!response.ok) {
		const errorText = await response.text().catch(() => "Unknown error");
		throw new Error(`API error (${response.status}): ${errorText}`);
	}
	return response.json() as Promise<T>;
}
