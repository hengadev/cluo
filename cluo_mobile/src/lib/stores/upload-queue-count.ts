import { writable } from "svelte/store";
import { list } from "$lib/upload-queue";

function createQueueCountStore() {
	const { subscribe, set } = writable(0);

	async function refresh() {
		try {
			const entries = await list();
			set(entries.length);
		} catch {
			// IndexedDB may be unavailable in SSR or private browsing
		}
	}

	return { subscribe, refresh };
}

export const queueCount = createQueueCountStore();
