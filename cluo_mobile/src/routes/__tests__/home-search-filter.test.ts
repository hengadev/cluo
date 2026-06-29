import { describe, it, expect } from "vitest";

/**
 * Tests for the home-page search/filter behaviour.
 *
 * The filter logic is extracted as a pure function so we can test it
 * without mounting Svelte components.
 */

import type { Recording } from "$lib/types/recording";

// ---------------------------------------------------------------------------
// Extracted filter function that mirrors the page derived logic
// ---------------------------------------------------------------------------

function filterRecordings(recordings: Recording[], query: string): Recording[] {
	const trimmed = query.trim();
	if (trimmed === "") return recordings;
	return recordings.filter((r) =>
		r.title.toLowerCase().includes(trimmed.toLowerCase()),
	);
}

// ---------------------------------------------------------------------------
// Fixtures
// ---------------------------------------------------------------------------

const recordings: Recording[] = [
	{ id: "r1", title: "Entretien client Dupont", date: "1 Jan 2026", startTime: "10:00", duration: 120, status: "completed", purpose: "general" },
	{ id: "r2", title: "Note sur le dossier Martin", date: "2 Jan 2026", startTime: "14:00", duration: 90, status: "completed", purpose: "general" },
	{ id: "r3", title: "Audition témoin — affaire Dupont", date: "3 Jan 2026", startTime: "09:00", duration: 200, status: "transcribing", purpose: "witness_interview" },
	{ id: "r4", title: "Résumé expertise médicale", date: "4 Jan 2026", startTime: "11:30", duration: 60, status: "completed", purpose: "general" },
];

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

describe("Home — search filter", () => {
	it("returns all recordings when query is empty", () => {
		expect(filterRecordings(recordings, "")).toEqual(recordings);
	});

	it("returns all recordings when query is only whitespace", () => {
		expect(filterRecordings(recordings, "   ")).toEqual(recordings);
	});

	it("filters by partial substring match", () => {
		const result = filterRecordings(recordings, "Dupont");
		expect(result).toHaveLength(2);
		expect(result.map((r) => r.id)).toEqual(["r1", "r3"]);
	});

	it("filters by case-insensitive match", () => {
		const result = filterRecordings(recordings, "dupont");
		expect(result).toHaveLength(2);
		expect(result.map((r) => r.id)).toEqual(["r1", "r3"]);
	});

	it("filters with mixed-case query", () => {
		const result = filterRecordings(recordings, "dUpOnT");
		expect(result).toHaveLength(2);
		expect(result.map((r) => r.id)).toEqual(["r1", "r3"]);
	});

	it("returns empty array when query matches nothing", () => {
		const result = filterRecordings(recordings, "xyznonexistent");
		expect(result).toEqual([]);
	});

	it("matches a single character", () => {
		const result = filterRecordings(recordings, "é");
		// "Résumé expertise médicale" contains "é"
		expect(result.length).toBeGreaterThanOrEqual(1);
		expect(result.some((r) => r.id === "r4")).toBe(true);
	});

	it("clearing the query restores the full list", () => {
		// Simulate: type a query, then clear it
		const afterSearch = filterRecordings(recordings, "Dupont");
		expect(afterSearch).toHaveLength(2);

		const afterClear = filterRecordings(recordings, "");
		expect(afterClear).toEqual(recordings);
		expect(afterClear).toHaveLength(4);
	});

	it("trims whitespace from the query before matching", () => {
		const result = filterRecordings(recordings, "  Dupont  ");
		expect(result).toHaveLength(2);
		expect(result.map((r) => r.id)).toEqual(["r1", "r3"]);
	});

	it("handles an empty recordings list gracefully", () => {
		expect(filterRecordings([], "Dupont")).toEqual([]);
		expect(filterRecordings([], "")).toEqual([]);
	});
});
