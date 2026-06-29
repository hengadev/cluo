import type { AppNotification } from "$lib/types/notifications";

/**
 * NotificationState — module-level singleton backing the header notification
 * bell.
 *
 * Mirrors the ToastState pattern but, unlike toasts, is exported as a single
 * shared instance so background services (job tracker, updater listener,
 * overdue-invoice checker, case-release call sites) can push notifications
 * without going through Svelte context.
 *
 * All mutations are mirrored to sessionStorage so the list survives in-app
 * navigation but is cleared when the app is closed.
 */

const STORAGE_KEY = "cluo.notifications.v1";

/** AppNotification with its Date serialised — the shape found in sessionStorage. */
type SerializedNotification = Omit<AppNotification, "createdAt"> & {
	createdAt: string;
};

/**
 * Returns the sessionStorage instance, or null when unavailable (SSR, private
 * mode, or disabled storage). Feature-detected rather than relying on
 * `$app/environment` so the store stays testable without a DOM.
 */
function getStorage(): Storage | null {
	try {
		if (typeof sessionStorage === "undefined") return null;
		return sessionStorage;
	} catch {
		return null;
	}
}

export class NotificationState {
	notifications = $state<AppNotification[]>([]);

	/**
	 * Number of notifications with `read === false`.
	 *
	 * Implemented as a getter over the `$state` array rather than `$derived`:
	 * reading it inside a component template still tracks `notifications` (so
	 * the bell re-renders on change), but unlike `$derived` it recomputes on
	 * every plain read — which keeps it correct and trivially unit-testable
	 * outside of a reactive effect.
	 */
	get unreadCount(): number {
		return this.notifications.filter((n) => !n.read).length;
	}

	constructor() {
		this.rehydrate();
	}

	/**
	 * Append a new notification. `id`, `createdAt`, and `read` are generated
	 * by the store — callers only describe *what* happened.
	 */
	push(notification: Omit<AppNotification, "id" | "createdAt" | "read">): void {
		this.notifications.push({
			...notification,
			id: crypto.randomUUID(),
			createdAt: new Date(),
			read: false,
		});
		this.persist();
	}

	/** Mark a single notification as read by id. No-op if already read or unknown. */
	markRead(id: string): void {
		const target = this.notifications.find((n) => n.id === id);
		if (!target || target.read) return;
		target.read = true;
		this.persist();
	}

	/** Mark every notification as read. */
	markAllRead(): void {
		let changed = false;
		for (const n of this.notifications) {
			if (!n.read) {
				n.read = true;
				changed = true;
			}
		}
		if (changed) this.persist();
	}

	/** Serialise the current list to sessionStorage. */
	private persist(): void {
		const storage = getStorage();
		if (!storage) return;
		try {
			storage.setItem(STORAGE_KEY, JSON.stringify(this.notifications));
		} catch {
			// Ignore quota / serialisation failures — notifications are best-effort.
		}
	}

	/** Restore the list from sessionStorage, if a valid snapshot exists. */
	private rehydrate(): void {
		const storage = getStorage();
		if (!storage) return;
		try {
			const raw = storage.getItem(STORAGE_KEY);
			if (!raw) return;
			const parsed = JSON.parse(raw) as SerializedNotification[];
			this.notifications = parsed.map((n) => ({
				...n,
				createdAt: new Date(n.createdAt),
			}));
		} catch {
			// Corrupt snapshot — start fresh rather than crash.
		}
	}
}

/** Shared notification store. Import this directly from anywhere. */
export const notificationStore = new NotificationState();
