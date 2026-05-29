import type { Handle } from '@sveltejs/kit';

// ---------------------------------------------------------------------------
// Security headers
// ---------------------------------------------------------------------------
// Applied to every response via a SvelteKit handle hook.
//
// CSP design notes:
//   • The Rapport is rendered via `{@html …}` after sanitisation through
//     sanitize-html (see client-access.ts → sanitizeReportHtml).  Only text
//     formatting tags are allowed — no <script>, <style>, or event handlers.
//     This means we can keep a strict CSP without 'unsafe-inline' for scripts.
//   • SvelteKit/Vite extracts component styles to linked CSS files in production,
//     so 'unsafe-inline' in style-src is not strictly required.  It is kept as
//     a safeguard against edge cases (e.g. Svelte transition directives that
//     emit inline style attributes at runtime).
//   • Media URLs (img, video, audio) originate from the API as absolute URLs
//     which may point to an external storage bucket (S3, etc.).  We whitelist
//     'self' plus https: so that any HTTPS media URL works.  In production you
//     may tighten this to the specific storage origin.
//   • connect-src is left open ('self' only) — the portal does not make
//     client-side XHR/fetch calls; everything is server-rendered.
//   • frame-ancestors 'none' mirrors X-Frame-Options: DENY.  Modern browsers
//     treat frame-ancestors as authoritative over X-Frame-Options; both are
//     kept for maximum compatibility.
// ---------------------------------------------------------------------------

const SECURITY_HEADERS: Record<string, string> = {
	'x-frame-options': 'DENY',
	'x-content-type-options': 'nosniff',
	'referrer-policy': 'strict-origin-when-cross-origin',
	'content-security-policy': [
		"default-src 'self'",
		"script-src 'self'",
		"style-src 'self' 'unsafe-inline'",
		"img-src 'self' https: data:",
		"media-src 'self' https:",
		"font-src 'self'",
		"connect-src 'self'",
		"frame-src 'none'",
		"frame-ancestors 'none'",
		"object-src 'none'",
		"base-uri 'self'",
		"form-action 'self'"
	].join('; ')
};

export const handle: Handle = async ({ event, resolve }) => {
	const response = await resolve(event);

	for (const [key, value] of Object.entries(SECURITY_HEADERS)) {
		response.headers.set(key, value);
	}

	return response;
};
