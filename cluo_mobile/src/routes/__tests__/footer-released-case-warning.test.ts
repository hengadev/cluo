import { describe, it, expect } from "vitest";
import type { Case } from "$lib/types/case";

/**
 * Pure predicate & decision logic extracted from Footer's keepRecording flow.
 *
 * shouldShowReleasedWarning returns true when the confirm dialog must appear.
 * resolveKeepAction returns the action the Footer should take:
 *   - "dialog"  → show confirmation dialog
 *   - "upload"  → proceed with upload directly
 *   - "noop"    → no blob, do nothing
 */

type KeepAction = "dialog" | "upload" | "noop";

function shouldShowReleasedWarning(
    currentCase: Case | null | undefined
): boolean {
    return currentCase?.status === "released";
}

function resolveKeepAction(
    hasRecordedBlob: boolean,
    currentCase: Case | null | undefined
): KeepAction {
    if (!hasRecordedBlob) return "noop";
    if (shouldShowReleasedWarning(currentCase)) return "dialog";
    return "upload";
}

const baseCase: Case = {
    id: "case-1",
    title: "Test Affair",
    status: "in_progress",
};

describe("Footer — released case warning", () => {
    describe("shouldShowReleasedWarning", () => {
        it("returns true when currentCase status is released", () => {
            expect(
                shouldShowReleasedWarning({ ...baseCase, status: "released" })
            ).toBe(true);
        });

        it("returns false when currentCase status is in_progress", () => {
            expect(shouldShowReleasedWarning(baseCase)).toBe(false);
        });

        it("returns false when currentCase status is ready", () => {
            expect(
                shouldShowReleasedWarning({ ...baseCase, status: "ready" })
            ).toBe(false);
        });

        it("returns false when currentCase is null", () => {
            expect(shouldShowReleasedWarning(null)).toBe(false);
        });

        it("returns false when currentCase is undefined", () => {
            expect(shouldShowReleasedWarning(undefined)).toBe(false);
        });
    });

    describe("resolveKeepAction", () => {
        it("returns 'noop' when there is no recorded blob", () => {
            expect(resolveKeepAction(false, baseCase)).toBe("noop");
        });

        it("returns 'dialog' for a released case with a blob", () => {
            expect(
                resolveKeepAction(true, { ...baseCase, status: "released" })
            ).toBe("dialog");
        });

        it("returns 'upload' for an in_progress case with a blob", () => {
            expect(resolveKeepAction(true, baseCase)).toBe("upload");
        });

        it("returns 'upload' for a ready case with a blob", () => {
            expect(
                resolveKeepAction(true, { ...baseCase, status: "ready" })
            ).toBe("upload");
        });

        it("returns 'noop' for no blob even with a released case", () => {
            expect(
                resolveKeepAction(false, { ...baseCase, status: "released" })
            ).toBe("noop");
        });
    });

    describe("confirm dialog interaction simulation", () => {
        it("cancel dismisses dialog without triggering upload", () => {
            const actions: string[] = [];

            const state = {
                confirmDialog: { show: false } as
                    | { show: true; onConfirm: () => void }
                    | { show: false },
                uploadCalled: false,
            };

            // Simulate keepRecording for released case
            function keepRecording() {
                state.confirmDialog = {
                    show: true,
                    onConfirm: () => {
                        state.confirmDialog = { show: false };
                        state.uploadCalled = true;
                    },
                };
            }

            keepRecording();
            expect(state.confirmDialog.show).toBe(true);
            expect(state.uploadCalled).toBe(false);

            // Cancel: just dismiss
            state.confirmDialog = { show: false };
            expect(state.uploadCalled).toBe(false);
        });

        it("confirm triggers upload", () => {
            const state = {
                confirmDialog: { show: false } as
                    | { show: true; onConfirm: () => void }
                    | { show: false },
                uploadCalled: false,
            };

            function keepRecording() {
                state.confirmDialog = {
                    show: true,
                    onConfirm: () => {
                        state.confirmDialog = { show: false };
                        state.uploadCalled = true;
                    },
                };
            }

            keepRecording();
            expect(state.confirmDialog.show).toBe(true);

            // Confirm
            (state.confirmDialog as { show: true; onConfirm: () => void }).onConfirm();
            expect(state.confirmDialog.show).toBe(false);
            expect(state.uploadCalled).toBe(true);
        });
    });
});
