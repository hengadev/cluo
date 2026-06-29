/**
 * Type definitions for application notifications.
 *
 * Notifications are a frontend-only concern. They are produced by background
 * services (job tracker, updater, overdue-invoice checker, case-release call
 * sites) and surfaced in the header notification centre. They are never
 * persisted to the backend — the only persistence layer is sessionStorage, so
 * they survive in-app navigation but are cleared when the app is closed.
 */

// =============================================================================
// NOTIFICATION KINDS
// =============================================================================

/**
 * The six notification categories surfaced in the notification centre.
 *
 * - `transcription_completed` — an audio transcription job finished successfully.
 * - `transcription_failed`     — a transcription job failed or was cancelled.
 * - `analysis_completed`       — an AI transcript analysis is ready to review.
 * - `invoice_overdue`          — an invoice has passed its due date unpaid.
 * - `update_available`         — a newer app version can be installed.
 * - `case_released`            — a dossier transitioned to the `released` state.
 */
export type NotificationKind =
	| 'transcription_completed'
	| 'transcription_failed'
	| 'analysis_completed'
	| 'invoice_overdue'
	| 'update_available'
	| 'case_released';

// =============================================================================
// NOTIFICATION
// =============================================================================

/**
 * A single notification item shown in the header bell.
 *
 * `kind` discriminates the category and drives the row icon (and, later,
 * click navigation). `caseId` / `mediaFileId` are optional navigation anchors
 * — when present, a future click handler can deep-link the investigator to the
 * relevant dossier or recording.
 */
export interface AppNotification {
	id: string;
	kind: NotificationKind;
	title: string;
	content: string;
	createdAt: Date;
	read: boolean;
	caseId?: string;
	mediaFileId?: string;
}
