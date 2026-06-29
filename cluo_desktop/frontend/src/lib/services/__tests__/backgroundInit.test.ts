/**
 * Unit tests for backgroundInit.ts — overdue-invoice notification wiring.
 *
 * The API module is mocked so no real HTTP is made; sessionStorage is stubbed
 * in-memory so the real notification store (which persists to it) works without
 * a DOM. Each test gets a fresh module + singletons via `vi.resetModules()` so
 * the dedup set and the notification store never leak between tests.
 */
import { describe, it, expect, beforeEach, vi } from "vitest";
import type { Invoice, OverdueInvoicesResponse } from "$lib/types/entities";

// ---------------------------------------------------------------------------
// Mocks — must be set up before importing the module under test
// ---------------------------------------------------------------------------

// Mock fetchOverdueInvoices — the overdue-invoice network touchpoint.
const fetchOverdueInvoices = vi.fn();
vi.mock("$lib/services/api", () => ({
	fetchOverdueInvoices: (...args: any[]) => fetchOverdueInvoices(...args),
}));

// Mock the Wails updater binding (imported lazily inside backgroundInit) so the
// update-availability path can be driven without the native bridge.
const checkForUpdate = vi.fn();
vi.mock("$lib/wailsjs/go/updater/Updater", () => ({
	CheckForUpdate: (...args: any[]) => checkForUpdate(...args),
}));

let initBackgroundCheckers: typeof import("../backgroundInit").initBackgroundCheckers;
let pushInvoiceIfNew: typeof import("../backgroundInit").pushInvoiceIfNew;
let pushUpdateIfNew: typeof import("../backgroundInit").pushUpdateIfNew;
let notificationStore: typeof import("$lib/stores/notifications.svelte").notificationStore;

// ---------------------------------------------------------------------------
// Helpers
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

function overdueResponse(data: Invoice[]): OverdueInvoicesResponse {
	return { success: true, data, total: data.length, page: 1, per_page: 20 };
}

function makeInvoice(overrides: Partial<Invoice> = {}): Invoice {
	return {
		id: "inv-1",
		case_id: "case-1",
		client_id: "client-1",
		invoice_number: "F-2026-001",
		issue_date: "2026-05-01",
		due_date: "2026-05-31",
		line_items: [],
		total_amount: 1500,
		tax_rate: 20,
		tax_amount: 300,
		payment_status: "overdue",
		status: "sent",
		created_at: "2026-05-01T00:00:00Z",
		updated_at: "2026-05-01T00:00:00Z",
		...overrides,
	};
}

beforeEach(async () => {
	vi.stubGlobal("sessionStorage", createMemoryStorage());
	fetchOverdueInvoices.mockReset();
	checkForUpdate.mockReset();
	vi.resetModules();
	({ initBackgroundCheckers, pushInvoiceIfNew, pushUpdateIfNew } = await import("../backgroundInit"));
	({ notificationStore } = await import("$lib/stores/notifications.svelte"));
});

// ---------------------------------------------------------------------------
// initBackgroundCheckers — overdue invoices (acceptance criteria)
// ---------------------------------------------------------------------------

describe("initBackgroundCheckers — overdue invoices", () => {
	it("pushes one invoice_overdue notification per overdue invoice returned", async () => {
		fetchOverdueInvoices.mockResolvedValueOnce(
			overdueResponse([makeInvoice({ id: "inv-1" }), makeInvoice({ id: "inv-2" })]),
		);

		await initBackgroundCheckers();

		expect(notificationStore.notifications).toHaveLength(2);
		expect(notificationStore.unreadCount).toBe(2);
		for (const notif of notificationStore.notifications) {
			expect(notif.kind).toBe("invoice_overdue");
			expect(notif.read).toBe(false);
			expect(notif.id).toEqual(expect.any(String));
			expect(notif.createdAt).toBeInstanceOf(Date);
		}
	});

	it("shapes content from the invoice reference and caseId from the dossier", async () => {
		fetchOverdueInvoices.mockResolvedValueOnce(
			overdueResponse([
				makeInvoice({ id: "inv-9", invoice_number: "F-2026-042", case_id: "case-9" }),
			]),
		);

		await initBackgroundCheckers();

		const notif = notificationStore.notifications[0];
		expect(notif.title).toBe("Facture en retard");
		expect(notif.content).toBe("F-2026-042");
		expect(notif.caseId).toBe("case-9");
	});

	it("does not push duplicates when run a second time in the same session", async () => {
		const invoices = [makeInvoice({ id: "inv-1" }), makeInvoice({ id: "inv-2" })];
		fetchOverdueInvoices.mockResolvedValue(overdueResponse(invoices));

		await initBackgroundCheckers();
		await initBackgroundCheckers();

		expect(fetchOverdueInvoices).toHaveBeenCalledTimes(2);
		expect(notificationStore.notifications).toHaveLength(2);
		expect(notificationStore.unreadCount).toBe(2);
	});

	it("pushes no notification when the endpoint returns an empty list", async () => {
		fetchOverdueInvoices.mockResolvedValueOnce(overdueResponse([]));

		await initBackgroundCheckers();

		expect(notificationStore.notifications).toHaveLength(0);
		expect(notificationStore.unreadCount).toBe(0);
	});

	it("swallows API errors and pushes nothing (never breaks app startup)", async () => {
		fetchOverdueInvoices.mockRejectedValueOnce(new Error("network down"));

		await expect(initBackgroundCheckers()).resolves.toBeUndefined();

		expect(notificationStore.notifications).toHaveLength(0);
	});
});

// ---------------------------------------------------------------------------
// pushInvoiceIfNew — dedup by invoice id
// ---------------------------------------------------------------------------

describe("pushInvoiceIfNew — dedup by invoice id", () => {
	it("pushing the same invoice id twice produces a single notification", () => {
		const invoice = makeInvoice({ id: "inv-1" });

		pushInvoiceIfNew(invoice);
		pushInvoiceIfNew(invoice);

		expect(notificationStore.notifications).toHaveLength(1);
		expect(notificationStore.unreadCount).toBe(1);
	});

	it("pushing two invoices with different ids produces two notifications", () => {
		pushInvoiceIfNew(makeInvoice({ id: "inv-1", invoice_number: "F-1" }));
		pushInvoiceIfNew(makeInvoice({ id: "inv-2", invoice_number: "F-2" }));

		expect(notificationStore.notifications).toHaveLength(2);
	});
});

// ---------------------------------------------------------------------------
// initBackgroundCheckers — update availability (acceptance criteria)
// ---------------------------------------------------------------------------

describe("initBackgroundCheckers — update availability", () => {
	// Silence the overdue-invoice checker so these tests isolate the update path.
	beforeEach(() => {
		fetchOverdueInvoices.mockResolvedValue(overdueResponse([]));
	});

	it("pushes an update_available notification when CheckForUpdate reports an update", async () => {
		checkForUpdate.mockResolvedValueOnce({ available: true, version: "1.2.3" });

		await initBackgroundCheckers();

		expect(notificationStore.notifications).toHaveLength(1);
		const notif = notificationStore.notifications[0];
		expect(notif.kind).toBe("update_available");
		expect(notif.title).toBe("Mise à jour disponible");
		expect(notif.content).toBe("Version 1.2.3");
		expect(notif.read).toBe(false);
		expect(notif.createdAt).toBeInstanceOf(Date);
	});

	it("does not push when no update is available", async () => {
		checkForUpdate.mockResolvedValueOnce({ available: false, version: "" });

		await initBackgroundCheckers();

		expect(notificationStore.notifications).toHaveLength(0);
		expect(notificationStore.unreadCount).toBe(0);
	});

	it("does not duplicate the same version across repeated checks", async () => {
		checkForUpdate.mockResolvedValue({ available: true, version: "1.2.3" });

		await initBackgroundCheckers();
		await initBackgroundCheckers();

		expect(checkForUpdate).toHaveBeenCalledTimes(2);
		expect(notificationStore.notifications).toHaveLength(1);
		expect(notificationStore.unreadCount).toBe(1);
	});

	it("surfaces a newer version as a second, distinct notification", async () => {
		checkForUpdate.mockResolvedValueOnce({ available: true, version: "1.2.3" });
		await initBackgroundCheckers();
		checkForUpdate.mockResolvedValueOnce({ available: true, version: "1.2.4" });
		await initBackgroundCheckers();

		expect(notificationStore.notifications).toHaveLength(2);
	});

	it("swallows updater errors and pushes nothing (never breaks app startup)", async () => {
		checkForUpdate.mockRejectedValueOnce(new Error("manifest fetch failed"));

		await expect(initBackgroundCheckers()).resolves.toBeUndefined();

		expect(notificationStore.notifications).toHaveLength(0);
	});
});

// ---------------------------------------------------------------------------
// pushUpdateIfNew — dedup by version
// ---------------------------------------------------------------------------

describe("pushUpdateIfNew — dedup by version", () => {
	it("pushing the same version twice produces a single notification", () => {
		pushUpdateIfNew("1.2.3");
		pushUpdateIfNew("1.2.3");

		expect(notificationStore.notifications).toHaveLength(1);
		expect(notificationStore.unreadCount).toBe(1);
	});

	it("pushing two different versions produces two notifications", () => {
		pushUpdateIfNew("1.2.3");
		pushUpdateIfNew("1.2.4");

		expect(notificationStore.notifications).toHaveLength(2);
	});

	it("is a no-op when an update_available for that version is already in the store", () => {
		notificationStore.push({ kind: "update_available", title: "Mise à jour", content: "Version 1.2.3" });

		pushUpdateIfNew("1.2.3");

		expect(notificationStore.notifications).toHaveLength(1);
	});

	it("still pushes for a different version even when another update is in the store", () => {
		notificationStore.push({ kind: "update_available", title: "Mise à jour", content: "Version 1.2.3" });

		pushUpdateIfNew("1.2.4");

		expect(notificationStore.notifications).toHaveLength(2);
	});
});
