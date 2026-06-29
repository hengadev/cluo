/**
 * Unit tests for caseActions.ts — case-release notification wiring.
 *
 * The API module is mocked so no real HTTP is made; sessionStorage is stubbed
 * in-memory so the real notification store (which persists to it) works without
 * a DOM. Each test gets a fresh module + singletons via `vi.resetModules()` so
 * the notification store never leaks between tests.
 */
import { describe, it, expect, beforeEach, vi } from "vitest";
import type { ReleaseResponse } from "$lib/types/entities";

// ---------------------------------------------------------------------------
// Mocks — must be set up before importing the module under test
// ---------------------------------------------------------------------------

// Mock releaseCase — the only network touchpoint in caseActions.
const releaseCase = vi.fn();
vi.mock("$lib/services/api", () => ({
	releaseCase: (...args: any[]) => releaseCase(...args),
}));

let releaseCaseAndNotify: typeof import("../caseActions").releaseCaseAndNotify;
let pushCaseReleasedNotification: typeof import("../caseActions").pushCaseReleasedNotification;
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

function makeReleaseResponse(caseId: string): ReleaseResponse {
	return {
		caseId,
		tokenId: "tok-" + caseId,
		rawToken: "raw-" + caseId,
		portalUrl: `https://portal.cluo.local/cases/${caseId}`,
		expiresAt: "2026-12-31T00:00:00Z",
	};
}

beforeEach(async () => {
	vi.stubGlobal("sessionStorage", createMemoryStorage());
	releaseCase.mockReset();
	vi.resetModules();
	({ releaseCaseAndNotify, pushCaseReleasedNotification } = await import("../caseActions"));
	({ notificationStore } = await import("$lib/stores/notifications.svelte"));
});

// ---------------------------------------------------------------------------
// releaseCaseAndNotify — acceptance criteria
// ---------------------------------------------------------------------------

describe("releaseCaseAndNotify", () => {
	it("pushes a single case_released notification on a successful release", async () => {
		releaseCase.mockResolvedValueOnce(makeReleaseResponse("case-7"));

		await releaseCaseAndNotify("case-7", "Surveillance Dupont");

		expect(releaseCase).toHaveBeenCalledWith("case-7");
		expect(notificationStore.notifications).toHaveLength(1);
		expect(notificationStore.unreadCount).toBe(1);

		const notif = notificationStore.notifications[0];
		expect(notif.kind).toBe("case_released");
		expect(notif.read).toBe(false);
		expect(notif.id).toEqual(expect.any(String));
		expect(notif.createdAt).toBeInstanceOf(Date);
	});

	it("carries the dossier title and caseId", async () => {
		releaseCase.mockResolvedValueOnce(makeReleaseResponse("case-9"));

		await releaseCaseAndNotify("case-9", "Vol Propriété Intellectuelle");

		const notif = notificationStore.notifications[0];
		expect(notif.title).toBe("Vol Propriété Intellectuelle");
		expect(notif.caseId).toBe("case-9");
	});

	it("uses content that indicates the portal link was generated", async () => {
		releaseCase.mockResolvedValueOnce(makeReleaseResponse("case-1"));

		await releaseCaseAndNotify("case-1", "Affaire test");

		expect(notificationStore.notifications[0].content).toContain("portail");
	});

	it("returns the release response to the caller", async () => {
		const response = makeReleaseResponse("case-1");
		releaseCase.mockResolvedValueOnce(response);

		const result = await releaseCaseAndNotify("case-1", "Affaire test");

		expect(result).toEqual(response);
	});

	it("does not push a notification when the release fails", async () => {
		releaseCase.mockRejectedValueOnce(new Error("network down"));

		await expect(releaseCaseAndNotify("case-1", "Affaire test")).rejects.toThrow(
			"network down",
		);

		expect(notificationStore.notifications).toHaveLength(0);
		expect(notificationStore.unreadCount).toBe(0);
	});
});

// ---------------------------------------------------------------------------
// pushCaseReleasedNotification — shaping (no network)
// ---------------------------------------------------------------------------

describe("pushCaseReleasedNotification", () => {
	it("shapes a case_released notification with title, content and caseId", () => {
		pushCaseReleasedNotification("case-42", "Harcèlement Travail");

		const notif = notificationStore.notifications[0];
		expect(notif.kind).toBe("case_released");
		expect(notif.title).toBe("Harcèlement Travail");
		expect(notif.caseId).toBe("case-42");
		expect(notif.read).toBe(false);
	});

	it("is additive — multiple releases stack in the bell", () => {
		pushCaseReleasedNotification("case-1", "Dossier A");
		pushCaseReleasedNotification("case-2", "Dossier B");

		expect(notificationStore.notifications).toHaveLength(2);
		expect(notificationStore.unreadCount).toBe(2);
	});
});
