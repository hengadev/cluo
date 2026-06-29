/**
 * backgroundInit — one-shot background checkers run on app launch.
 *
 * `initBackgroundCheckers()` is called once from the `(app)/+layout.svelte`
 * `onMount`. Each checker is best-effort: a failure in one must never prevent
 * the others from running or crash the app, so every checker swallows its own
 * errors.
 *
 * v1 scope:
 * - Overdue invoices: query the `OverdueInvoicesResponse` endpoint once and
 *   push an `invoice_overdue` notification per result, deduplicated by invoice
 *   id so repeat init calls (and component hot-reloads that re-run the layout's
 *   onMount) do not double-push.
 * - Update available: call the Wails `CheckForUpdate` binding once and, when it
 *   reports a newer version, push an `update_available` notification carrying
 *   the version number, deduplicated against the store so the same version is
 *   only ever surfaced once.
 *
 * Later issues (case_released call site) plug into this same entry point
 * without changing its signature.
 */
import { fetchOverdueInvoices } from "$lib/services/api";
import { notificationStore } from "$lib/stores/notifications.svelte";
import type { Invoice } from "$lib/types/entities";

/**
 * Invoice ids already surfaced this session.
 *
 * Module-scoped so it survives repeat `initBackgroundCheckers()` calls and
 * component hot-reloads within the same module lifetime. It is intentionally
 * *not* persisted: the notifications themselves are already persisted to
 * sessionStorage by the store; this guard only exists to avoid duplicate pushes
 * during a single module lifetime.
 */
const notifiedInvoiceIds = new Set<string>();

/**
 * Push a single overdue-invoice notification, unless one has already been pushed
 * for this invoice id this session.
 *
 * Extracted from the network path so the dedup + notification-shaping logic can
 * be unit-tested without going through the API.
 */
export function pushInvoiceIfNew(invoice: Invoice): void {
	if (notifiedInvoiceIds.has(invoice.id)) return;
	notifiedInvoiceIds.add(invoice.id);

	notificationStore.push({
		kind: "invoice_overdue",
		title: "Facture en retard",
		content: invoice.invoice_number,
		caseId: invoice.case_id,
	});
}

/**
 * Query the overdue-invoices endpoint once and push an `invoice_overdue`
 * notification per result. Empty responses produce no notifications; a network
 * or API failure is swallowed so the caller never sees it.
 */
async function pushOverdueInvoiceNotifications(): Promise<void> {
	const response = await fetchOverdueInvoices();
	for (const invoice of response.data) {
		pushInvoiceIfNew(invoice);
	}
}

// ---------------------------------------------------------------------------
// Update available
// ---------------------------------------------------------------------------

/**
 * How availability is actually detected (and why this is not an
 * `EventsOn("updater:status", …)` listener).
 *
 * The notification PRD described this checker as listening to the Wails
 * `updater:status` event for an `"available"` status. That does not match the
 * current backend: `Updater.emitStatus` only ever emits `"downloading"`,
 * `"installing"` and `"ready"` (see cluo_desktop/updater/updater.go) — never
 * `"available"` — and the status event is a bare string that carries no
 * version. The only thing that knows whether an update is available, and what
 * its version is, is `Updater.CheckForUpdate`, which returns an
 * `UpdateInfo { available, version, … }`.
 *
 * So instead of waiting on an event that never fires, the checker calls
 * `CheckForUpdate` directly. This keeps detection frontend-owned, needs no
 * backend change, and leaves `UpdatePrompt` (which makes its own
 * `CheckForUpdate` call on mount) completely untouched.
 */

/**
 * Push a single `update_available` notification for `version`, unless one for
 * that exact version already lives in the store.
 *
 * Dedup is store-backed rather than a module-scoped set: the notification is
 * already persisted to sessionStorage by the store, so checking the store also
 * covers re-runs of `initBackgroundCheckers` across an in-app navigation or a
 * dev hot-reload (where a module-scoped set would have been reset but the
 * store would still hold the earlier notification).
 *
 * Extracted from the updater path so the dedup + shaping logic can be
 * unit-tested without the Wails binding.
 */
export function pushUpdateIfNew(version: string): void {
	const content = `Version ${version}`;
	const alreadyNotified = notificationStore.notifications.some(
		(n) => n.kind === "update_available" && n.content === content,
	);
	if (alreadyNotified) return;

	notificationStore.push({
		kind: "update_available",
		title: "Mise à jour disponible",
		content,
	});
}

/**
 * Ask the Wails updater whether a newer version exists and, if so, surface it
 * as a notification. The updater binding is imported lazily so dev builds and
 * the test environment (where the Wails bridge is absent) do not crash on
 * import; any rejection is swallowed by the caller so this never blocks app
 * startup.
 */
async function pushUpdateNotification(): Promise<void> {
	const { CheckForUpdate } = await import("$lib/wailsjs/go/updater/Updater");
	const info = await CheckForUpdate();
	if (info.available) {
		pushUpdateIfNew(info.version);
	}
}

/**
 * Run every background checker once. Safe to call multiple times — each checker
 * is idempotent within a session. Errors are swallowed so a failing endpoint
 * never breaks app startup.
 */
export async function initBackgroundCheckers(): Promise<void> {
	await pushOverdueInvoiceNotifications().catch(() => {
		// Overdue-invoice endpoint unavailable — nothing to surface. The app
		// still starts normally.
	});
	await pushUpdateNotification().catch(() => {
		// Updater binding missing (dev/test), manifest not configured, or a
		// network error — nothing to surface. `UpdatePrompt` still handles the
		// interactive install flow independently of this notification.
	});
}
