import { describe, it, expect, vi } from "vitest";

// Mock import.meta.env before importing the module
vi.stubEnv("VITE_API_URL", "https://api.test.local");

// Mock localStorage for the module-level chain helpers
const store: Record<string, string> = {};
vi.stubGlobal("localStorage", {
	getItem: vi.fn((key: string) => store[key] ?? null),
	setItem: vi.fn((key: string, value: string) => {
		store[key] = value;
	}),
	removeItem: vi.fn((key: string) => {
		delete store[key];
	}),
	clear: vi.fn(() => {
		Object.keys(store).forEach((k) => delete store[k]);
	}),
});

const { formatDate, formatTime } = await import("../audio");

describe("formatDate", () => {
	it("formats an ISO timestamp in fr-FR locale", () => {
		// Fixed timestamp: 2026-06-02T14:32:00Z
		const result = formatDate("2026-06-02T14:32:00Z");
		// fr-FR short month for June is "juin"
		expect(result).toContain("juin");
		expect(result).toContain("2026");
		expect(result).toContain("02");
	});

	it("formats a winter date correctly", () => {
		const result = formatDate("2026-01-15T09:00:00Z");
		// "janv" covers both "janv." (full ICU) and "janv" (minimal ICU)
		expect(result).toMatch(/janv\.?/);
		expect(result).toContain("2026");
		expect(result).toContain("15");
	});

	it("always uses 2-digit day", () => {
		const result = formatDate("2026-03-05T08:00:00Z");
		expect(result).toContain("05");
	});
});

describe("formatTime", () => {
	it("returns a 24h HH:MM string", () => {
		const result = formatTime("2026-06-02T14:32:00Z");
		expect(result).toMatch(/^\d{2}:\d{2}$/);
	});

	it("pads single-digit hours with a leading zero", () => {
		const result = formatTime("2026-06-02T09:05:00Z");
		expect(result).toMatch(/^\d{2}:\d{2}$/);
	});

	it("formats midnight correctly", () => {
		const result = formatTime("2026-06-02T00:00:00Z");
		expect(result).toMatch(/^\d{2}:\d{2}$/);
	});

	it("never contains AM/PM — always 24h format", () => {
		const result = formatTime("2026-06-02T14:32:00Z");
		expect(result).not.toMatch(/AM|PM/i);
	});
});
