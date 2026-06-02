import { describe, it, expect } from "vitest";

/**
 * Pure functions extracted from Footer.svelte recording-title logic.
 * These mirror the component's generateTimestampTitle and effectiveTitle helpers.
 */

function generateTimestampTitle(date: Date = new Date()): string {
    const hours = date.getHours().toString().padStart(2, "0");
    const minutes = date.getMinutes().toString().padStart(2, "0");
    return `Enregistrement ${hours}h${minutes}`;
}

function effectiveTitle(currentTitle: string, defaultTitle: string): string {
    return currentTitle.trim() || defaultTitle;
}

describe("Footer — recording title", () => {
    describe("generateTimestampTitle", () => {
        it("produces the expected format at a known time", () => {
            // 14:32 local time
            const date = new Date(2025, 0, 15, 14, 32);
            expect(generateTimestampTitle(date)).toBe("Enregistrement 14h32");
        });

        it("zero-pads hours and minutes", () => {
            const date = new Date(2025, 0, 15, 2, 5);
            expect(generateTimestampTitle(date)).toBe("Enregistrement 02h05");
        });

        it("handles midnight", () => {
            const date = new Date(2025, 0, 15, 0, 0);
            expect(generateTimestampTitle(date)).toBe("Enregistrement 00h00");
        });

        it("handles 23h59", () => {
            const date = new Date(2025, 0, 15, 23, 59);
            expect(generateTimestampTitle(date)).toBe("Enregistrement 23h59");
        });
    });

    describe("effectiveTitle", () => {
        const defaultTitle = "Enregistrement 14h32";

        it("returns custom title when provided", () => {
            expect(effectiveTitle("Mon enregistrement", defaultTitle)).toBe("Mon enregistrement");
        });

        it("falls back to default when current title is empty string", () => {
            expect(effectiveTitle("", defaultTitle)).toBe(defaultTitle);
        });

        it("falls back to default when current title is only whitespace", () => {
            expect(effectiveTitle("   ", defaultTitle)).toBe(defaultTitle);
        });

        it("uses custom title even if it differs from default", () => {
            expect(effectiveTitle("Interview client", defaultTitle)).toBe("Interview client");
        });

        it("trims whitespace from custom title before deciding", () => {
            expect(effectiveTitle("  Titre custom  ", defaultTitle)).toBe("Titre custom");
        });
    });
});
