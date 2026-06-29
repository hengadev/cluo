/**
 * Unit tests for the job tracker (jobTracker.ts).
 *
 * The tracker owns a `setInterval` poll per job, so these tests use vitest's
 * fake timers (`vi.useFakeTimers` + `advanceTimersByTimeAsync`) to step through
 * ticks without real delays. `getTranscriptionJobStatus` is mocked so no HTTP
 * is made; the real notification store (with an in-memory sessionStorage) is
 * used so we assert on the actual notifications that get pushed.
 *
 * Each test gets a fresh tracker + store via `vi.resetModules()` so polled
 * state never leaks between tests.
 */
import { describe, it, expect, beforeEach, vi } from 'vitest';

// ---------------------------------------------------------------------------
// Mocks — `getTranscriptionJobStatus` is the only API surface the tracker
// touches. The factory references the `mock*`-prefixed variable, which vitest
// hoists alongside the mock registration.
// ---------------------------------------------------------------------------

const mockGetStatus = vi.fn();
vi.mock('../api', () => ({
	getTranscriptionJobStatus: mockGetStatus,
}));

// ---------------------------------------------------------------------------
// In-memory sessionStorage so the notification store can persist/rehydrate
// without a DOM.
// ---------------------------------------------------------------------------

function createMemoryStorage(): Storage {
	let store: Record<string, string> = {};
	return {
		getItem: (key: string) => (key in store ? store[key] : null),
		setItem: (key: string, value: string) => {
			store[key] = String(value);
		},
		removeItem: (key: string) => {
			delete store[key];
		},
		clear: () => {
			store = {};
		},
		key: (index: number) => Object.keys(store)[index] ?? null,
		get length() {
			return Object.keys(store).length;
		},
	};
}

let jobTracker: typeof import('../jobTracker').jobTracker;
let notificationStore: typeof import('../../stores/notifications.svelte').notificationStore;

beforeEach(async () => {
	vi.stubGlobal('sessionStorage', createMemoryStorage());
	vi.useFakeTimers();
	vi.resetModules();
	mockGetStatus.mockReset();

	({ jobTracker } = await import('../jobTracker'));
	({ notificationStore } = await import('../../stores/notifications.svelte'));
});

// ---------------------------------------------------------------------------
// Completion
// ---------------------------------------------------------------------------

describe('JobTracker — completion', () => {
	it('pushes exactly one transcription_completed notification when a job completes and then stops polling', async () => {
		mockGetStatus
			.mockResolvedValueOnce({ status: 'processing', mediaFileId: 'media-1' })
			.mockResolvedValueOnce({ status: 'processing', mediaFileId: 'media-1' })
			.mockResolvedValueOnce({ status: 'completed', mediaFileId: 'media-1' })
			// Any further calls (should not happen once polling stops).
			.mockResolvedValue({ status: 'completed', mediaFileId: 'media-1' });

		jobTracker.trackJob('job-1', 'recording-001.mp3', 'case-9');

		// Two processing ticks → no notification yet.
		await vi.advanceTimersByTimeAsync(5000);
		await vi.advanceTimersByTimeAsync(5000);
		expect(notificationStore.notifications).toHaveLength(0);

		// Third tick → completed.
		await vi.advanceTimersByTimeAsync(5000);

		expect(mockGetStatus).toHaveBeenCalledTimes(3);
		expect(notificationStore.notifications).toHaveLength(1);

		const notif = notificationStore.notifications[0];
		expect(notif.kind).toBe('transcription_completed');
		expect(notif.title).toBe('Transcription terminée');
		expect(notif.content).toContain('recording-001.mp3');
		expect(notif.caseId).toBe('case-9');
		expect(notif.mediaFileId).toBe('media-1');
		expect(notif.read).toBe(false);

		// Polling has stopped: advancing several more ticks must not push again
		// or make any further status calls.
		await vi.advanceTimersByTimeAsync(5000);
		await vi.advanceTimersByTimeAsync(5000);
		await vi.advanceTimersByTimeAsync(5000);

		expect(mockGetStatus).toHaveBeenCalledTimes(3);
		expect(notificationStore.notifications).toHaveLength(1);
	});

	it('completes on the very first poll', async () => {
		mockGetStatus.mockResolvedValue({ status: 'completed', mediaFileId: 'media-7' });

		jobTracker.trackJob('job-fast', 'fast.mp3', 'case-1');
		await vi.advanceTimersByTimeAsync(5000);

		expect(notificationStore.notifications).toHaveLength(1);
		expect(notificationStore.notifications[0].kind).toBe('transcription_completed');
		expect(mockGetStatus).toHaveBeenCalledTimes(1);
	});
});

// ---------------------------------------------------------------------------
// Failure / cancellation
// ---------------------------------------------------------------------------

describe('JobTracker — failure', () => {
	it('pushes exactly one transcription_failed notification when a job fails and then stops polling', async () => {
		mockGetStatus
			.mockResolvedValueOnce({ status: 'processing', mediaFileId: 'media-2' })
			.mockResolvedValueOnce({ status: 'failed', mediaFileId: 'media-2', errorMessage: 'boom' })
			.mockResolvedValue({ status: 'failed', mediaFileId: 'media-2' });

		jobTracker.trackJob('job-fail', 'broken.mp3', 'case-2');

		await vi.advanceTimersByTimeAsync(5000);
		expect(notificationStore.notifications).toHaveLength(0);

		await vi.advanceTimersByTimeAsync(5000);

		expect(mockGetStatus).toHaveBeenCalledTimes(2);
		expect(notificationStore.notifications).toHaveLength(1);

		const notif = notificationStore.notifications[0];
		expect(notif.kind).toBe('transcription_failed');
		expect(notif.title).toBe('Transcription échouée');
		expect(notif.content).toContain('broken.mp3');
		expect(notif.caseId).toBe('case-2');
		expect(notif.mediaFileId).toBe('media-2');

		// Stopped — no further calls or notifications.
		await vi.advanceTimersByTimeAsync(20000);
		expect(mockGetStatus).toHaveBeenCalledTimes(2);
		expect(notificationStore.notifications).toHaveLength(1);
	});

	it('treats a cancelled job as a transcription_failed notification', async () => {
		mockGetStatus.mockResolvedValueOnce({ status: 'cancelled', mediaFileId: 'media-3' });

		jobTracker.trackJob('job-cancel', 'canceled.mp3', 'case-3');
		await vi.advanceTimersByTimeAsync(5000);

		expect(notificationStore.notifications).toHaveLength(1);
		const notif = notificationStore.notifications[0];
		expect(notif.kind).toBe('transcription_failed');
		expect(notif.title).toBe('Transcription annulée');
		expect(mockGetStatus).toHaveBeenCalledTimes(1);
	});
});

// ---------------------------------------------------------------------------
// Deduplication
// ---------------------------------------------------------------------------

describe('JobTracker — deduplication', () => {
	it('calling trackJob twice with the same jobId registers only one poll interval', async () => {
		mockGetStatus.mockResolvedValue({ status: 'processing', mediaFileId: 'media-dup' });

		jobTracker.trackJob('dup', 'a.mp3', 'case-dup');
		jobTracker.trackJob('dup', 'a.mp3', 'case-dup');

		await vi.advanceTimersByTimeAsync(5000);

		// A single interval means a single status call per tick, not two.
		expect(mockGetStatus).toHaveBeenCalledTimes(1);

		await vi.advanceTimersByTimeAsync(5000);
		expect(mockGetStatus).toHaveBeenCalledTimes(2);
	});

	it('different jobIds each get their own poll', async () => {
		mockGetStatus.mockResolvedValue({ status: 'processing', mediaFileId: 'media-x' });

		jobTracker.trackJob('job-a', 'a.mp3', 'case-a');
		jobTracker.trackJob('job-b', 'b.mp3', 'case-b');

		await vi.advanceTimersByTimeAsync(5000);

		// Two independent jobs → two status calls per tick.
		expect(mockGetStatus).toHaveBeenCalledTimes(2);
	});
});

// ---------------------------------------------------------------------------
// Timeout
// ---------------------------------------------------------------------------

describe('JobTracker — timeout', () => {
	it('stops polling after 60 attempts without a terminal status and pushes no notification', async () => {
		mockGetStatus.mockResolvedValue({ status: 'processing', mediaFileId: 'media-stuck' });

		jobTracker.trackJob('stuck', 'forever.mp3', 'case-stuck');

		// Advance through all 60 ticks (60 × 5 s).
		await vi.advanceTimersByTimeAsync(60 * 5000);

		expect(mockGetStatus).toHaveBeenCalledTimes(60);
		expect(notificationStore.notifications).toHaveLength(0);

		// The 61st tick must not fire — polling has stopped.
		await vi.advanceTimersByTimeAsync(5000);
		await vi.advanceTimersByTimeAsync(5000);

		expect(mockGetStatus).toHaveBeenCalledTimes(60);
		expect(notificationStore.notifications).toHaveLength(0);
	});

	it('a transient poll error does not abort polling and does not push a notification', async () => {
		mockGetStatus
			.mockRejectedValueOnce(new Error('network blip'))
			.mockResolvedValueOnce({ status: 'completed', mediaFileId: 'media-retry' });

		jobTracker.trackJob('flaky', 'flaky.mp3', 'case-flaky');

		await vi.advanceTimersByTimeAsync(5000); // tick 1: throws → swallowed
		expect(notificationStore.notifications).toHaveLength(0);

		await vi.advanceTimersByTimeAsync(5000); // tick 2: completed
		expect(notificationStore.notifications).toHaveLength(1);
		expect(notificationStore.notifications[0].kind).toBe('transcription_completed');
		expect(mockGetStatus).toHaveBeenCalledTimes(2);
	});
});
