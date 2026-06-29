import { getTranscriptionJobStatus } from './api';
import { notificationStore } from '../stores/notifications.svelte';

/**
 * JobTracker — module-level singleton that owns transcription job polling for
 * the whole application lifetime.
 *
 * Before this existed, the recordings page ran its own `setInterval` poll.
 * That poll died the moment the investigator navigated away from the page, so a
 * transcription submitted on dossier A would silently complete (or fail) with
 * nobody watching. The job tracker moves that poll out of page-local scope so
 * it survives navigation, and reports the terminal result through the
 * notification store instead of a page-local toast.
 *
 * Public surface is intentionally tiny: callers only `trackJob(...)`. The
 * tracker deduplicates by `jobId`, polls on a fixed interval, pushes exactly
 * one notification on a terminal status, and gives up silently after a bounded
 * number of attempts.
 */

/** Interval between status polls — matches the previous per-page poll. */
const POLL_INTERVAL_MS = 5000;

/**
 * Maximum number of polls before the tracker gives up.
 *
 * 60 attempts at 5 s ≈ 5 minutes. At that point the job is presumed stuck and
 * the tracker stops polling silently (no notification) rather than spinning
 * forever.
 */
const MAX_ATTEMPTS = 60;

/** Mutable bookkeeping for a single tracked job. */
interface TrackedJob {
	/** Handle returned by `setInterval`, used to stop the poll. */
	intervalId: ReturnType<typeof setInterval> | undefined;
	/** How many polls have run for this job so far. */
	attempts: number;
}

export class JobTracker {
	/** Active jobs keyed by `jobId`. Presence here is what deduplicates polls. */
	private tracked = new Map<string, TrackedJob>();

	/**
	 * Begin polling a transcription job and report its result through the
	 * notification store.
	 *
	 * Calling this with a `jobId` that is already tracked is a no-op, so it is
	 * safe to call repeatedly (e.g. on every navigation to the recordings
	 * page) without stacking duplicate polls.
	 *
	 * `fileName` and `caseId` are captured at submit time so the eventual
	 * notification can describe the recording and deep-link back to the
	 * dossier even if the caller has long since navigated away.
	 */
	trackJob(jobId: string, fileName: string, caseId: string): void {
		if (this.tracked.has(jobId)) return;

		const entry: TrackedJob = { intervalId: undefined, attempts: 0 };
		entry.intervalId = setInterval(() => {
			void this.poll(jobId, entry, fileName, caseId);
		}, POLL_INTERVAL_MS);
		this.tracked.set(jobId, entry);
	}

	/**
	 * One poll tick for a tracked job.
	 *
	 * On a terminal status (`completed` / `failed` / `cancelled`) the job is
	 * untracked and exactly one notification is pushed. After `MAX_ATTEMPTS`
	 * ticks without a terminal status the job is untracked silently. Transient
	 * errors from `getTranscriptionJobStatus` are swallowed so a single flaky
	 * response does not kill the poll.
	 */
	private async poll(
		jobId: string,
		entry: TrackedJob,
		fileName: string,
		caseId: string,
	): Promise<void> {
		entry.attempts++;

		try {
			const job = await getTranscriptionJobStatus(jobId);

			if (job.status === 'completed') {
				this.stop(jobId);
				notificationStore.push({
					kind: 'transcription_completed',
					title: 'Transcription terminée',
					content: `« ${fileName} » a été transcrit.`,
					caseId,
					mediaFileId: job.mediaFileId,
				});
				return;
			}

			if (job.status === 'failed' || job.status === 'cancelled') {
				this.stop(jobId);
				const cancelled = job.status === 'cancelled';
				notificationStore.push({
					kind: 'transcription_failed',
					title: cancelled ? 'Transcription annulée' : 'Transcription échouée',
					content: `La transcription de « ${fileName} » ${
						cancelled ? 'a été annulée.' : 'a échoué.'
					}`,
					caseId,
					mediaFileId: job.mediaFileId,
				});
				return;
			}
		} catch {
			// Transient poll error — retry on the next tick.
		}

		if (entry.attempts >= MAX_ATTEMPTS) {
			this.stop(jobId);
		}
	}

	/** Clear the poll interval and forget the job. Idempotent. */
	private stop(jobId: string): void {
		const entry = this.tracked.get(jobId);
		if (!entry) return;
		if (entry.intervalId !== undefined) clearInterval(entry.intervalId);
		this.tracked.delete(jobId);
	}
}

/** Shared job tracker. Import this directly from anywhere. */
export const jobTracker = new JobTracker();
