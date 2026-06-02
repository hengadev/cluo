import { describe, it, expect, beforeEach } from "vitest";
import { get } from "svelte/store";
import { currentCase } from "$lib/stores/current-case";
import type { Case } from "$lib/types/case";

const mockCase: Case = {
    id: "case-1",
    title: "Test Affair",
    status: "in_progress",
};

describe("currentCase store", () => {
    beforeEach(() => {
        currentCase.set(null);
    });

    it("initialises to null", () => {
        expect(get(currentCase)).toBeNull();
    });

    it("holds a Case after set", () => {
        currentCase.set(mockCase);
        expect(get(currentCase)).toEqual(mockCase);
    });

    it("can be reset to null", () => {
        currentCase.set(mockCase);
        currentCase.set(null);
        expect(get(currentCase)).toBeNull();
    });

    it("can be updated to a different Case", () => {
        currentCase.set(mockCase);
        const other: Case = { id: "case-2", title: "Other", status: "ready" };
        currentCase.set(other);
        expect(get(currentCase)).toEqual(other);
    });
});
