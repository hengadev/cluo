/// <reference types="@sveltejs/kit" />
import { build, files, version } from '$service-worker';

// Unique cache name for this deployment
const CACHE = `cache-${version}`;

// Precached build assets (app shell) + static files
const ASSETS = [...build, ...files];

// Navigation fallback shown when the network is unreachable
const OFFLINE_FALLBACK = '/offline';

self.addEventListener('install', (event) => {
	async function addFilesToCache() {
		const cache = await caches.open(CACHE);
		await cache.addAll(ASSETS);
		// Cache the offline page separately so a failure here can't abort install
		await cache.add(OFFLINE_FALLBACK).catch(() => {
			console.warn('[sw] could not precache offline fallback');
		});
	}

	event.waitUntil(addFilesToCache());
});

self.addEventListener('activate', (event) => {
	// Remove previous cached data from disk
	async function deleteOldCaches() {
		for (const key of await caches.keys()) {
			if (key !== CACHE) await caches.delete(key);
		}
	}

	event.waitUntil(deleteOldCaches());
});

self.addEventListener('fetch', (event) => {
	// ignore POST requests etc
	if (event.request.method !== 'GET') return;

	async function respond() {
		const url = new URL(event.request.url);
		const sameOrigin = url.origin === self.location.origin;
		const cache = await caches.open(CACHE);

		// `build`/`files` are versioned app assets: cache-first
		if (ASSETS.includes(url.pathname)) {
			const cached = await cache.match(url.pathname);
			if (cached) return cached;
		}

		// Cross-origin requests (e.g. the Go backend API): network-only.
		// Never read from or write to the cache — avoids serving stale or
		// leaking authenticated responses while offline.
		if (!sameOrigin) {
			return fetch(event.request);
		}

		const isNavigation = event.request.mode === 'navigate';

		// Same-origin requests: network-first, cache fallback for offline use
		try {
			const response = await fetch(event.request);

			// if we're offline, fetch can return a value that is not a Response
			// instead of throwing - and we can't pass this non-Response to respondWith
			if (!(response instanceof Response)) {
				throw new Error('invalid response from fetch');
			}

			if (response.status === 200) {
				cache.put(event.request, response.clone());
			}

			return response;
		} catch (err) {
			const cached = await cache.match(event.request);
			if (cached) return cached;

			// No cached copy: show the offline page for navigations
			if (isNavigation) {
				const offline = await cache.match(OFFLINE_FALLBACK);
				if (offline) return offline;
			}

			// if there's no cache, then just error out
			// as there is nothing we can do to respond to this request
			throw err;
		}
	}

	event.respondWith(respond());
});
