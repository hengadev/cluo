/**
 * Recording store for managing audio recording state.
 * Uses Svelte stores with IndexedDB persistence for offline support.
 */

import { writable, type Writable } from "svelte/store";
import type { Recording, RecordingStatus } from "$lib/types/recording";

/**
 * Store for the current recording being captured.
 */
interface CurrentRecording {
	blob: Blob | null;
	duration: number;
	metadata: {
		caseId?: string;
		title?: string;
		createdAt: string;
	};
}

/**
 * IndexedDB configuration for storing recordings offline.
 */
const DB_NAME = "cluo_recordings";
const DB_VERSION = 1;
const STORE_NAME = "recordings";

/**
 * Open IndexedDB connection.
 */
function openDB(): Promise<IDBDatabase> {
	return new Promise((resolve, reject) => {
		const request = indexedDB.open(DB_NAME, DB_VERSION);

		request.onerror = () => reject(request.error);
		request.onsuccess = () => resolve(request.result);

		request.onupgradeneeded = (event) => {
			const db = (event.target as IDBOpenDBRequest).result;
			if (!db.objectStoreNames.contains(STORE_NAME)) {
				const store = db.createObjectStore(STORE_NAME, { keyPath: "id" });
				store.createIndex("caseId", "caseId", { unique: false });
				store.createIndex("createdAt", "metadata.createdAt", { unique: false });
			}
		};
	});
}

/**
 * Save a recording to IndexedDB.
 */
async function saveRecordingToIndexedDB(recording: Recording & { blob: Blob }): Promise<void> {
	const db = await openDB();
	return new Promise((resolve, reject) => {
		const transaction = db.transaction([STORE_NAME], "readwrite");
		const store = transaction.objectStore(STORE_NAME);
		const request = store.put(recording);

		request.onerror = () => reject(request.error);
		request.onsuccess = () => resolve();
	});
}

/**
 * Get a recording from IndexedDB by ID.
 */
async function getRecordingFromIndexedDB(id: string): Promise<(Recording & { blob: Blob }) | null> {
	const db = await openDB();
	return new Promise((resolve, reject) => {
		const transaction = db.transaction([STORE_NAME], "readonly");
		const store = transaction.objectStore(STORE_NAME);
		const request = store.get(id);

		request.onerror = () => reject(request.error);
		request.onsuccess = () => resolve(request.result ?? null);
	});
}

/**
 * Delete a recording from IndexedDB.
 */
async function deleteRecordingFromIndexedDB(id: string): Promise<void> {
	const db = await openDB();
	return new Promise((resolve, reject) => {
		const transaction = db.transaction([STORE_NAME], "readwrite");
		const store = transaction.objectStore(STORE_NAME);
		const request = store.delete(id);

		request.onerror = () => reject(request.error);
		request.onsuccess = () => resolve();
	});
}

/**
 * List all recordings from IndexedDB.
 */
async function listRecordingsFromIndexedDB(): Promise<(Recording & { blob: Blob })[]> {
	const db = await openDB();
	return new Promise((resolve, reject) => {
		const transaction = db.transaction([STORE_NAME], "readonly");
		const store = transaction.objectStore(STORE_NAME);
		const request = store.getAll();

		request.onerror = () => reject(request.error);
		request.onsuccess = () => resolve(request.result ?? []);
	});
}

/**
 * Current recording store - holds the recording being captured.
 */
export const currentRecording: Writable<CurrentRecording> = writable({
	blob: null,
	duration: 0,
	metadata: {
		createdAt: new Date().toISOString(),
	},
});

/**
 * Recordings list store - holds all recordings (from backend or IndexedDB).
 */
export const recordings = writable<Recording[]>([]);

/**
 * Recording state store - holds UI state for recording.
 */
export const recordingState = writable<{
	isRecording: boolean;
	isPlaying: boolean;
	canUpload: boolean;
}>({
	isRecording: false,
	isPlaying: false,
	canUpload: false,
});

/**
 * Save a recording to the store and optionally to IndexedDB.
 */
export async function saveRecording(
	recording: Recording,
	blob?: Blob,
): Promise<void> {
	if (blob) {
		await saveRecordingToIndexedDB({ ...recording, blob });
	}
}

/**
 * Get a recording by ID from store or IndexedDB.
 */
export async function getRecording(id: string): Promise<(Recording & { blob: Blob }) | null> {
	// Try IndexedDB first (for offline support)
	return getRecordingFromIndexedDB(id);
}

/**
 * Delete a recording from store and IndexedDB.
 */
export async function deleteRecording(id: string): Promise<void> {
	await deleteRecordingFromIndexedDB(id);

	// Update recordings list
	recordings.update((current) => current.filter((r) => r.id !== id));
}

/**
 * List all recordings.
 * Returns from IndexedDB if offline, otherwise would fetch from backend.
 */
export async function listRecordings(): Promise<Recording[]> {
	return listRecordingsFromIndexedDB();
}

/**
 * Update the current recording state.
 */
export function setCurrentRecording(
	blob: Blob | null,
	duration: number,
	metadata?: Partial<CurrentRecording["metadata"]>,
): void {
	currentRecording.set({
		blob,
		duration,
		metadata: {
			caseId: metadata?.caseId,
			title: metadata?.title,
			createdAt: metadata?.createdAt ?? new Date().toISOString(),
		},
	});

	// Enable upload button if we have a blob
	recordingState.update((state) => ({
		...state,
		canUpload: blob !== null,
	}));
}

/**
 * Clear the current recording.
 */
export function clearCurrentRecording(): void {
	currentRecording.set({
		blob: null,
		duration: 0,
		metadata: {
			createdAt: new Date().toISOString(),
		},
	});

	recordingState.update((state) => ({
		...state,
		canUpload: false,
		isPlaying: false,
	}));
}

/**
 * Set the recording state.
 */
export function setRecordingState(
	key: keyof typeof recordingState extends string ? "isRecording" | "isPlaying" | "canUpload" : never,
	value: boolean,
): void {
	recordingState.update((state) => ({
		...state,
		[key]: value,
	}));
}
