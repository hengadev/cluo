import { describe, it, expect, vi, beforeEach } from "vitest";

/**
 * Tests for the home-page case-switch + recordings-reload behaviour.
 *
 * The fetch logic is extracted as a pure async function so we can test it
 * without mounting Svelte components.
 */

import type { Recording } from "$lib/types/recording";
import type { Case } from "$lib/types/case";

// ---------------------------------------------------------------------------
// Extracted state machine that mirrors the page logic
// ---------------------------------------------------------------------------

interface RecordingsState {
    recordings: Recording[];
    error: string | null;
    loading: boolean;
}

async function fetchRecordingsForCase(
    caseId: string,
    listFn: (opts?: { caseId?: string }) => Promise<{ recordings: Recording[]; totalCount: number }>,
): Promise<RecordingsState> {
    const state: RecordingsState = { recordings: [], error: null, loading: true };
    try {
        const res = await listFn({ caseId });
        state.recordings = res.recordings;
    } catch (e) {
        state.error = e instanceof Error ? e.message : "Échec du chargement des enregistrements";
    } finally {
        state.loading = false;
    }
    return state;
}

// ---------------------------------------------------------------------------
// Fixtures
// ---------------------------------------------------------------------------

const caseA: Case = { id: "case-a", title: "Case A", status: "in_progress" };
const caseB: Case = { id: "case-b", title: "Case B", status: "ready" };

const recordingsA: Recording[] = [
    { id: "r1", title: "Recording A1", date: "1 Jan 2026", startTime: "10:00", duration: 120, status: "completed" },
];
const recordingsB: Recording[] = [
    { id: "r2", title: "Recording B1", date: "2 Jan 2026", startTime: "14:00", duration: 90, status: "completed" },
    { id: "r3", title: "Recording B2", date: "3 Jan 2026", startTime: "09:00", duration: 200, status: "transcribing" },
];

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

describe("Home — case switch recordings reload", () => {
    let mockListRecordings: ReturnType<typeof vi.fn>;

    beforeEach(() => {
        mockListRecordings = vi.fn();
    });

    it("fetches recordings for the new case", async () => {
        mockListRecordings.mockResolvedValue({ recordings: recordingsB, totalCount: 2 });

        const state = await fetchRecordingsForCase(caseB.id, mockListRecordings);

        expect(mockListRecordings).toHaveBeenCalledWith({ caseId: "case-b" });
        expect(state.recordings).toEqual(recordingsB);
        expect(state.loading).toBe(false);
        expect(state.error).toBeNull();
    });

    it("replaces the list when switching cases", async () => {
        // First fetch returns case A recordings
        mockListRecordings.mockResolvedValueOnce({ recordings: recordingsA, totalCount: 1 });
        const stateA = await fetchRecordingsForCase(caseA.id, mockListRecordings);
        expect(stateA.recordings).toEqual(recordingsA);

        // Second fetch returns case B recordings
        mockListRecordings.mockResolvedValueOnce({ recordings: recordingsB, totalCount: 2 });
        const stateB = await fetchRecordingsForCase(caseB.id, mockListRecordings);
        expect(stateB.recordings).toEqual(recordingsB);
        expect(stateB.recordings).not.toEqual(recordingsA);
    });

    it("shows error message on fetch failure", async () => {
        mockListRecordings.mockRejectedValue(new Error("Network error"));

        const state = await fetchRecordingsForCase(caseA.id, mockListRecordings);

        expect(state.error).toBe("Network error");
        expect(state.recordings).toEqual([]);
        expect(state.loading).toBe(false);
    });

    it("shows fallback error for non-Error throws", async () => {
        mockListRecordings.mockRejectedValue("string error");

        const state = await fetchRecordingsForCase(caseA.id, mockListRecordings);

        expect(state.error).toBe("Échec du chargement des enregistrements");
        expect(state.loading).toBe(false);
    });

    it("returns loading false even when fetch succeeds", async () => {
        mockListRecordings.mockResolvedValue({ recordings: [], totalCount: 0 });

        const state = await fetchRecordingsForCase(caseA.id, mockListRecordings);

        expect(state.loading).toBe(false);
        expect(state.recordings).toEqual([]);
    });

    it("calls listRecordings with the correct caseId each time", async () => {
        mockListRecordings.mockResolvedValue({ recordings: [], totalCount: 0 });

        await fetchRecordingsForCase("case-x", mockListRecordings);
        await fetchRecordingsForCase("case-y", mockListRecordings);
        await fetchRecordingsForCase("case-z", mockListRecordings);

        expect(mockListRecordings).toHaveBeenCalledTimes(3);
        expect(mockListRecordings).toHaveBeenNthCalledWith(1, { caseId: "case-x" });
        expect(mockListRecordings).toHaveBeenNthCalledWith(2, { caseId: "case-y" });
        expect(mockListRecordings).toHaveBeenNthCalledWith(3, { caseId: "case-z" });
    });
});
