const API_URL = import.meta.env.VITE_API_URL ?? "";

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

// JSON helper: prepends API_URL, sets JSON headers, throws on non-ok
export async function apiFetch<T>(
	path: string,
	options?: RequestInit,
): Promise<T> {
	const response = await fetchWithRetry(`${API_URL}${path}`, {
		credentials: "include",
		...options,
		headers: {
			"Content-Type": "application/json",
			...options?.headers,
		},
	});
	if (!response.ok) {
		const errorText = await response.text().catch(() => "Unknown error");
		throw new Error(`API error (${response.status}): ${errorText}`);
	}
	return response.json() as Promise<T>;
}

// Raw helper: prepends API_URL, sets JSON headers, returns Response as-is
export async function apiFetchRaw(
	path: string,
	options?: RequestInit,
): Promise<Response> {
	return fetchWithRetry(`${API_URL}${path}`, {
		credentials: "include",
		...options,
		headers: {
			"Content-Type": "application/json",
			...options?.headers,
		},
	});
}
