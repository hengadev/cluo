import { describe, it, expect } from "vitest";
import type { Case } from "$lib/types/case";

/**
 * Pure predicate extracted from Footer's handleDragStart guard logic.
 * This mirrors the condition: `if (!currentCase) return;` which prevents
 * drag-start when no active Case is set.
 */
function canStartDrag(currentCase: Case | null | undefined): boolean {
    return currentCase != null;
}

const mockCase: Case = {
    id: "case-1",
    title: "Test Affair",
    status: "in_progress",
};

describe("Footer — active case guard", () => {
    describe("canStartDrag", () => {
        it("returns false when currentCase is null", () => {
            expect(canStartDrag(null)).toBe(false);
        });

        it("returns false when currentCase is undefined", () => {
            expect(canStartDrag(undefined)).toBe(false);
        });

        it("returns true when currentCase is set", () => {
            expect(canStartDrag(mockCase)).toBe(true);
        });

        it("returns true for a case with 'released' status", () => {
            const releasedCase: Case = {
                ...mockCase,
                status: "released",
            };
            expect(canStartDrag(releasedCase)).toBe(true);
        });

        it("returns true for a case with minimal fields", () => {
            const minimalCase: Case = {
                id: "x",
                title: "",
                status: "in_progress",
            };
            expect(canStartDrag(minimalCase)).toBe(true);
        });
    });
});
