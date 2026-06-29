/**
 * caseActions — user-triggered case lifecycle actions that produce
 * notifications.
 *
 * Unlike the background checkers in `backgroundInit.ts` (overdue invoices,
 * updater), case lifecycle transitions are initiated by the investigator from
 * the UI. Each action calls the corresponding API function and, on success,
 * pushes a notification so the investigator has a persistent record in the bell
 * — for example, that the client portal link has been generated.
 *
 * The notification shaping is extracted into `pushCaseReleasedNotification` so
 * it can be unit-tested without going through the network, mirroring the
 * `pushInvoiceIfNew` split in `backgroundInit.ts`.
 */
import { releaseCase } from "$lib/services/api";
import { notificationStore } from "$lib/stores/notifications.svelte";
import type { ReleaseResponse } from "$lib/types/entities";

/**
 * Push a `case_released` notification for a dossier.
 *
 * The dossier title becomes the notification heading so the investigator can
 * see *which* case was released at a glance, and the content notes that the
 * portal link has been generated (per the PRD). `caseId` anchors
 * click-navigation to the dossier's informations page.
 */
export function pushCaseReleasedNotification(caseId: string, title: string): void {
	notificationStore.push({
		kind: "case_released",
		title,
		content: "Le lien du portail client a été généré.",
		caseId,
	});
}

/**
 * Release a dossier to the client portal and surface a `case_released`
 * notification on success.
 *
 * Any API error is re-thrown unchanged — the caller is responsible for
 * surfacing a toast — but the notification is pushed only when the release
 * resolves. This mirrors how transcription completion is handled by the job
 * tracker: the notification is a consequence of success, never of failure.
 */
export async function releaseCaseAndNotify(
	caseId: string,
	title: string,
): Promise<ReleaseResponse> {
	const response = await releaseCase(caseId);
	pushCaseReleasedNotification(caseId, title);
	return response;
}
