/**
 * Offline upload queue backed by IndexedDB.
 *
 * Extends the existing `cluo_recordings` database (bumps to version 2)
 * with a new `upload_queue` object store.
 *
 * All functions are pure — no UI or component dependencies.
 */

// ---------------------------------------------------------------------------
// Types
// ---------------------------------------------------------------------------

export interface QueuedUpload {
	id: string;
	blob: Blob;
	metadata: { caseId: string; title: string; purpose?: string };
	enqueuedAt: string; // ISO 8601
	attemptCount: number;
}

export interface FlushResult {
	succeeded: string[]; // queue IDs removed
	failed: string[]; // queue IDs that remain
}

// ---------------------------------------------------------------------------
// IndexedDB helpers
// ---------------------------------------------------------------------------

const DB_NAME = "cluo_recordings";
const DB_VERSION = 2;
const STORE_NAME = "upload_queue";

let dbInstance: IDBDatabase | null = null;
let dbReady: Promise<IDBDatabase> | null = null;

function openDB(): Promise<IDBDatabase> {
	if (dbInstance) return Promise.resolve(dbInstance);
	if (dbReady) return dbReady;

	dbReady = new Promise<IDBDatabase>((resolve, reject) => {
		const request = indexedDB.open(DB_NAME, DB_VERSION);

		request.onerror = () => {
			dbReady = null;
			reject(request.error);
		};

		request.onsuccess = () => {
			dbInstance = request.result;
			resolve(request.result);
		};

		request.onupgradeneeded = (event) => {
			const db = (event.target as IDBOpenDBRequest).result;

			// Create upload_queue store if it doesn't exist (new in v2)
			if (!db.objectStoreNames.contains(STORE_NAME)) {
				db.createObjectStore(STORE_NAME, { keyPath: "id" });
			}
		};

		request.onblocked = () => {
			dbReady = null;
			reject(new Error("IndexedDB upgrade blocked — close other tabs to continue"));
		};
	});

	return dbReady;
}

function txStore(mode: IDBTransactionMode): Promise<IDBObjectStore> {
	return openDB().then((db) => {
		const transaction = db.transaction([STORE_NAME], mode);
		return transaction.objectStore(STORE_NAME);
	});
}

// ---------------------------------------------------------------------------
// Public API
// ---------------------------------------------------------------------------

/**
 * Enqueue a failed upload so it survives page refresh / app restart.
 * Returns the generated queue ID.
 */
export async function enqueue(
	blob: Blob,
	metadata: { caseId: string; title: string; purpose?: string },
): Promise<string> {
	const id = crypto.randomUUID();
	const entry: QueuedUpload = {
		id,
		blob,
		metadata,
		enqueuedAt: new Date().toISOString(),
		attemptCount: 0,
	};

	const store = await txStore("readwrite");

	return new Promise<string>((resolve, reject) => {
		const req = store.add(entry);
		req.onsuccess = () => resolve(id);
		req.onerror = () => reject(req.error);
	});
}

/**
 * List all entries currently in the queue.
 */
export async function list(): Promise<QueuedUpload[]> {
	const store = await txStore("readonly");

	return new Promise<QueuedUpload[]>((resolve, reject) => {
		const req = store.getAll();
		req.onsuccess = () => resolve(req.result ?? []);
		req.onerror = () => reject(req.error);
	});
}

/**
 * Remove a single entry by its queue ID.
 */
export async function remove(queueId: string): Promise<void> {
	const store = await txStore("readwrite");

	return new Promise<void>((resolve, reject) => {
		const req = store.delete(queueId);
		req.onsuccess = () => resolve();
		req.onerror = () => reject(req.error);
	});
}

/**
 * Attempt to upload every queued entry using the provided `uploadFn`.
 *
 * - Succeeded entries are removed from the queue.
 * - Failed entries stay with an incremented `attemptCount`.
 * - Idempotent: safe to call concurrently (uses an internal lock).
 */
export async function flush(
	uploadFn: (blob: Blob, metadata: { caseId?: string; title?: string; purpose?: string }) => Promise<unknown>,
): Promise<FlushResult> {
	// Simple lock to prevent concurrent flushes
	if (flushInProgress) {
		return { succeeded: [], failed: [] };
	}
	flushInProgress = true;

	try {
		const entries = await list();
		if (entries.length === 0) {
			return { succeeded: [], failed: [] };
		}

		const succeeded: string[] = [];
		const failed: string[] = [];

		// Process entries sequentially to avoid overwhelming the network
		for (const entry of entries) {
			try {
				await uploadFn(entry.blob, entry.metadata);
				await remove(entry.id);
				succeeded.push(entry.id);
			} catch {
				// Increment attemptCount
				const store = await txStore("readwrite");
				await new Promise<void>((resolve, reject) => {
					const updated: QueuedUpload = {
						...entry,
						attemptCount: entry.attemptCount + 1,
					};
					const req = store.put(updated);
					req.onsuccess = () => resolve();
					req.onerror = () => reject(req.error);
				});
				failed.push(entry.id);
			}
		}

		return { succeeded, failed };
	} finally {
		flushInProgress = false;
	}
}

let flushInProgress = false;

/**
 * Reset internal state (for testing only).
 * Closes the cached DB connection if open.
 */
export function _reset(): void {
	if (dbInstance) {
		dbInstance.close();
		dbInstance = null;
	}
	dbReady = null;
	flushInProgress = false;
}
