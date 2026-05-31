import { PUBLIC_APP_ENV } from '$env/static/public';
import { json } from '@sveltejs/kit';

const isStaging = PUBLIC_APP_ENV === 'staging';

export function GET() {
	return json({
		name: isStaging ? 'Cluo [Staging]' : 'Cluo Mobile',
		short_name: isStaging ? 'Cluo [S]' : 'Cluo',
		description: 'Cluo Mobile Application',
		start_url: '/',
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
	});
}
