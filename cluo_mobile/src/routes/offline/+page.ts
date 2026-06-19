// Static fallback served by the service worker when the network is down.
// Kept self-contained (inline styles, no app state) so it renders even when
// the rest of the app shell isn't cached. Listed explicitly in
// `kit.prerender.entries` so it's always emitted as a static asset.
export const prerender = true;
