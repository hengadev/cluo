import { env } from '$env/dynamic/public';

const { PUBLIC_APP_ENV } = env;

const isStaging = PUBLIC_APP_ENV === 'staging';

// Served as `application/manifest+json` (the spec-correct MIME type) rather than
// via SvelteKit's `json()` helper, which would send `application/json` and trip
// up stricter browsers/validators.
const manifest = {
	// Stable identity for the installed app across reinstalls.
	id: '/',
	name: isStaging ? 'Cluo [Staging]' : 'Cluo Mobile',
	short_name: isStaging ? 'Cluo [S]' : 'Cluo',
	description: 'Cluo Mobile Application',
	start_url: '/',
	scope: '/',
	display: 'standalone',
	background_color: isStaging ? '#1a1a1a' : '#ffffff',
	theme_color: isStaging ? '#1a1a1a' : '#000000',
	orientation: 'portrait-primary',
	icons: [
		{
			src: isStaging ? '/icon-staging-192.png' : '/icon-192.png',
			sizes: '192x192',
			type: 'image/png',
			purpose: 'any maskable'
		},
		{
			src: isStaging ? '/icon-staging-512.png' : '/icon-512.png',
			sizes: '512x512',
			type: 'image/png',
			purpose: 'any maskable'
		}
	]
};

export function GET() {
	return new Response(JSON.stringify(manifest), {
		headers: {
			'content-type': 'application/manifest+json; charset=utf-8',
			'cache-control': 'public, max-age=3600'
		}
	});
}
