/**
 * Unit tests for the notification store (notifications.svelte.ts).
 *
 * sessionStorage is stubbed with an in-memory implementation so the store's
 * persistence round-trip can be exercised without a DOM. Each test gets a
 * fresh module + singleton via `vi.resetModules()` + dynamic import, so state
 * never leaks between tests.
 */
import { describe, it, expect, beforeEach, vi } from "vitest";

const STORAGE_KEY = "cluo.notifications.v1";

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

let notificationStore: typeof import("../notifications.svelte").notificationStore;

beforeEach(async () => {
	vi.stubGlobal("sessionStorage", createMemoryStorage());
	vi.resetModules();
	({ notificationStore } = await import("../notifications.svelte"));
});

describe("NotificationState.push", () => {
	it("adds a notification with read:false and a generated id/createdAt", () => {
		notificationStore.push({
			kind: "transcription_completed",
			title: "Transcription terminée",
			content: "enregistrement-001.mp3",
		});

		expect(notificationStore.notifications).toHaveLength(1);

		const notif = notificationStore.notifications[0];
		expect(notif.id).toEqual(expect.any(String));
		expect(notif.id).not.toBe("");
		expect(notif.kind).toBe("transcription_completed");
		expect(notif.read).toBe(false);
		expect(notif.createdAt).toBeInstanceOf(Date);
	});

	it("increments unreadCount", () => {
		expect(notificationStore.unreadCount).toBe(0);

		notificationStore.push({ kind: "update_available", title: "Mise à jour", content: "2.4.0" });
		notificationStore.push({ kind: "analysis_completed", title: "Analyse prête", content: "Résumé" });

		expect(notificationStore.unreadCount).toBe(2);
	});

	it("preserves optional navigation fields", () => {
		notificationStore.push({
			kind: "transcription_failed",
			title: "Échec",
			content: "job-123",
			caseId: "case-9",
			mediaFileId: "media-1",
		});

		const notif = notificationStore.notifications[0];
		expect(notif.caseId).toBe("case-9");
		expect(notif.mediaFileId).toBe("media-1");
	});
});

describe("NotificationState.markRead", () => {
	it("marks a single notification read and decrements unreadCount", () => {
		notificationStore.push({ kind: "invoice_overdue", title: "Facture", content: "F-12" });
		const id = notificationStore.notifications[0].id;

		expect(notificationStore.unreadCount).toBe(1);

		notificationStore.markRead(id);

		expect(notificationStore.notifications[0].read).toBe(true);
		expect(notificationStore.unreadCount).toBe(0);
	});

	it("is a no-op for an unknown id", () => {
		notificationStore.push({ kind: "invoice_overdue", title: "Facture", content: "F-12" });
		notificationStore.markRead("does-not-exist");

		expect(notificationStore.unreadCount).toBe(1);
		expect(notificationStore.notifications[0].read).toBe(false);
	});

	it("does not double-count an already-read notification", () => {
		notificationStore.push({ kind: "invoice_overdue", title: "Facture", content: "F-12" });
		const id = notificationStore.notifications[0].id;

		notificationStore.markRead(id);
		expect(notificationStore.unreadCount).toBe(0);

		// Marking the same id again must not push unreadCount negative.
		notificationStore.markRead(id);
		expect(notificationStore.unreadCount).toBe(0);
	});
});

describe("NotificationState.markAllRead", () => {
	it("sets every notification to read and zeros unreadCount", () => {
		notificationStore.push({ kind: "invoice_overdue", title: "A", content: "a" });
		notificationStore.push({ kind: "case_released", title: "B", content: "b" });
		notificationStore.push({ kind: "update_available", title: "C", content: "c" });

		// One already read on its own.
		notificationStore.markRead(notificationStore.notifications[0].id);
		expect(notificationStore.unreadCount).toBe(2);

		notificationStore.markAllRead();

		expect(notificationStore.unreadCount).toBe(0);
		for (const n of notificationStore.notifications) {
			expect(n.read).toBe(true);
		}
	});

	it("reflects the correct unread count when pushing after markAllRead", () => {
		notificationStore.push({ kind: "invoice_overdue", title: "A", content: "a" });
		notificationStore.markAllRead();
		expect(notificationStore.unreadCount).toBe(0);

		notificationStore.push({ kind: "case_released", title: "B", content: "b" });
		expect(notificationStore.unreadCount).toBe(1);
	});
});

describe("NotificationState sessionStorage persistence", () => {
	it("round-trips: pushed + marked-read notifications rehydrate into a fresh store", async () => {
		notificationStore.push({
			kind: "transcription_failed",
			title: "Échec",
			content: "job-123",
			caseId: "case-9",
			mediaFileId: "media-1",
		});
		notificationStore.push({ kind: "update_available", title: "MAJ", content: "2.5.0" });
		notificationStore.markRead(notificationStore.notifications[0].id);

		// A brand-new store (fresh module) must rebuild from the same sessionStorage.
		vi.resetModules();
		const { notificationStore: freshStore } = await import("../notifications.svelte");

		expect(freshStore.notifications).toHaveLength(2);

		// First notification kept its kind + navigation anchors + read state.
		const first = freshStore.notifications[0];
		expect(first.kind).toBe("transcription_failed");
		expect(first.caseId).toBe("case-9");
		expect(first.mediaFileId).toBe("media-1");
		expect(first.read).toBe(true);
		expect(first.createdAt).toBeInstanceOf(Date);

		// Second notification is still unread → unreadCount survives the trip.
		expect(freshStore.unreadCount).toBe(1);
	});

	it("writes a snapshot to sessionStorage on every mutation", () => {
		expect(sessionStorage.getItem(STORAGE_KEY)).toBeNull();

		notificationStore.push({ kind: "update_available", title: "MAJ", content: "2.5.0" });
		const afterPush = sessionStorage.getItem(STORAGE_KEY);
		expect(afterPush).not.toBeNull();
		expect(JSON.parse(afterPush!)).toHaveLength(1);

		notificationStore.markAllRead();
		const snapshot = JSON.parse(sessionStorage.getItem(STORAGE_KEY)!);
		expect(snapshot.every((n: { read: boolean }) => n.read)).toBe(true);
	});

	it("starts empty when sessionStorage has no snapshot", () => {
		expect(notificationStore.notifications).toHaveLength(0);
		expect(notificationStore.unreadCount).toBe(0);
	});

	it("ignores a corrupt snapshot and starts fresh", async () => {
		sessionStorage.setItem(STORAGE_KEY, "{ not valid json");

		vi.resetModules();
		const { notificationStore: freshStore } = await import("../notifications.svelte");

		expect(freshStore.notifications).toHaveLength(0);
	});
});
