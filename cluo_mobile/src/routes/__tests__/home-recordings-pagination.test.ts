import { describe, it, expect, vi, beforeEach, type Mock } from "vitest";

/**
 * Tests for the home-page recordings pagination behaviour.
 *
 * The pagination logic is extracted as a pure state machine so we can test it
 * without mounting Svelte components.
 */

import type { Recording } from "$lib/types/recording";

// ---------------------------------------------------------------------------
// Extracted pagination state machine that mirrors the page logic
// ---------------------------------------------------------------------------

interface PaginationState {
	recordings: Recording[];
	totalCount: number;
	loadingMore: boolean;
	currentCaseId: string;
}

const PAGE_SIZE = 20;

const remainingCount = (s: PaginationState) => Math.max(0, s.totalCount - s.recordings.length);
const hasMore = (s: PaginationState) => remainingCount(s) > 0;

async function loadMore(
	state: PaginationState,
	listFn: (opts?: { caseId?: string; offset?: number; limit?: number }) => Promise<{ recordings: Recording[]; totalCount: number }>,
): Promise<PaginationState> {
	if (state.loadingMore || !hasMore(state) || !state.currentCaseId) return state;

	const next: PaginationState = {
		recordings: [...state.recordings],
		totalCount: state.totalCount,
		loadingMore: true,
		currentCaseId: state.currentCaseId,
	};

	try {
		const res = await listFn({
			caseId: state.currentCaseId,
			offset: state.recordings.length,
			limit: PAGE_SIZE,
		});
		next.recordings = [...state.recordings, ...res.recordings];
		next.totalCount = res.totalCount;
	} finally {
		next.loadingMore = false;
	}

	return next;
}

// ---------------------------------------------------------------------------
// Fixtures
// ---------------------------------------------------------------------------

function makeRecording(id: string): Recording {
	return {
		id,
		title: `Recording ${id}`,
		date: "1 Jan 2026",
		startTime: "10:00",
		duration: 120,
		status: "completed",
		purpose: "general",
	};
}

const firstPage: Recording[] = Array.from({ length: 20 }, (_, i) => makeRecording(`r${i + 1}`));
const secondPage: Recording[] = Array.from({ length: 5 }, (_, i) => makeRecording(`r${i + 21}`));

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

describe("Home — recordings pagination", () => {
	type ListFn = (opts?: { caseId?: string; offset?: number; limit?: number }) => Promise<{ recordings: Recording[]; totalCount: number }>;
	let mockListRecordings: Mock<ListFn>;

	beforeEach(() => {
		mockListRecordings = vi.fn();
	});

	describe("remainingCount and hasMore", () => {
		it("shows button when totalCount > loaded recordings", () => {
			const state: PaginationState = { recordings: firstPage, totalCount: 25, loadingMore: false, currentCaseId: "case-1" };
			expect(hasMore(state)).toBe(true);
			expect(remainingCount(state)).toBe(5);
		});

		it("shows correct remaining count", () => {
			const state: PaginationState = { recordings: firstPage, totalCount: 50, loadingMore: false, currentCaseId: "case-1" };
			expect(remainingCount(state)).toBe(30);
		});

		it("hides button when all recordings are loaded", () => {
			const state: PaginationState = { recordings: firstPage, totalCount: 20, loadingMore: false, currentCaseId: "case-1" };
			expect(hasMore(state)).toBe(false);
			expect(remainingCount(state)).toBe(0);
		});

		it("hides button when recordings exceed totalCount (defensive)", () => {
			const state: PaginationState = { recordings: firstPage, totalCount: 15, loadingMore: false, currentCaseId: "case-1" };
			expect(hasMore(state)).toBe(false);
			expect(remainingCount(state)).toBe(0);
		});

		it("hides button when totalCount is 0", () => {
			const state: PaginationState = { recordings: [], totalCount: 0, loadingMore: false, currentCaseId: "case-1" };
			expect(hasMore(state)).toBe(false);
		});
	});

	describe("loadMore", () => {
		it("appends results to the existing list", async () => {
			mockListRecordings.mockResolvedValue({ recordings: secondPage, totalCount: 25 });

			const state: PaginationState = { recordings: firstPage, totalCount: 25, loadingMore: false, currentCaseId: "case-1" };
			const result = await loadMore(state, mockListRecordings);

			expect(result.recordings).toHaveLength(25);
			expect(result.recordings.slice(0, 20)).toEqual(firstPage);
			expect(result.recordings.slice(20)).toEqual(secondPage);
			expect(result.totalCount).toBe(25);
			expect(result.loadingMore).toBe(false);
		});

		it("calls listRecordings with correct offset and limit", async () => {
			mockListRecordings.mockResolvedValue({ recordings: secondPage, totalCount: 25 });

			const state: PaginationState = { recordings: firstPage, totalCount: 25, loadingMore: false, currentCaseId: "case-1" };
			await loadMore(state, mockListRecordings);

			expect(mockListRecordings).toHaveBeenCalledWith({
				caseId: "case-1",
				offset: 20,
				limit: 20,
			});
		});

		it("does not fetch when all recordings are already loaded", async () => {
			const state: PaginationState = { recordings: firstPage, totalCount: 20, loadingMore: false, currentCaseId: "case-1" };
			const result = await loadMore(state, mockListRecordings);

			expect(mockListRecordings).not.toHaveBeenCalled();
			expect(result.recordings).toEqual(firstPage);
		});

		it("does not fetch when already loading", async () => {
			const state: PaginationState = { recordings: firstPage, totalCount: 25, loadingMore: true, currentCaseId: "case-1" };
			await loadMore(state, mockListRecordings);

			expect(mockListRecordings).not.toHaveBeenCalled();
		});

		it("does not fetch when there is no current case", async () => {
			const state: PaginationState = { recordings: firstPage, totalCount: 25, loadingMore: false, currentCaseId: "" };
			await loadMore(state, mockListRecordings);

			expect(mockListRecordings).not.toHaveBeenCalled();
		});

		it("re-throws on error so the caller can handle it", async () => {
			mockListRecordings.mockRejectedValue(new Error("Network error"));

			const state: PaginationState = { recordings: firstPage, totalCount: 25, loadingMore: false, currentCaseId: "case-1" };
			await expect(loadMore(state, mockListRecordings)).rejects.toThrow("Network error");
		});

		it("button disappears after all pages loaded", async () => {
			mockListRecordings.mockResolvedValue({ recordings: secondPage, totalCount: 25 });

			const state: PaginationState = { recordings: firstPage, totalCount: 25, loadingMore: false, currentCaseId: "case-1" };
			const result = await loadMore(state, mockListRecordings);

			expect(hasMore(result)).toBe(false);
			expect(remainingCount(result)).toBe(0);
		});

		it("supports multiple sequential loads", async () => {
			const page2 = Array.from({ length: 20 }, (_, i) => makeRecording(`r${i + 21}`));
			const page3 = Array.from({ length: 3 }, (_, i) => makeRecording(`r${i + 41}`));

			mockListRecordings.mockResolvedValueOnce({ recordings: page2, totalCount: 43 });

			let state: PaginationState = { recordings: firstPage, totalCount: 43, loadingMore: false, currentCaseId: "case-1" };
			state = await loadMore(state, mockListRecordings);

			expect(state.recordings).toHaveLength(40);
			expect(hasMore(state)).toBe(true);
			expect(remainingCount(state)).toBe(3);

			mockListRecordings.mockResolvedValueOnce({ recordings: page3, totalCount: 43 });
			state = await loadMore(state, mockListRecordings);

			expect(state.recordings).toHaveLength(43);
			expect(hasMore(state)).toBe(false);
		});
	});
});
