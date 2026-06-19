import adapter from '@sveltejs/adapter-node';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Consult https://svelte.dev/docs/kit/integrations
	// for more information about preprocessors
	preprocess: vitePreprocess(),

	kit: {
		adapter: adapter(),
		// Prerender the offline fallback so it's always available as a static
		// asset for the service worker to cache, even under server load.
		prerender: {
			entries: ['/offline']
		}
	}
};

export default config;
