/**
 * Snackbar store for feedback messages with optional retry actions.
 */
import { writable } from "svelte/store";

export interface Snackbar {
    id: string;
    message: string;
    type: "info" | "error";
    action?: {
        label: string;
        onClick: () => void;
    };
}

function createSnackbarStore() {
    const { subscribe, set } = writable<Snackbar | null>(null);
    let dismissTimer: ReturnType<typeof setTimeout> | null = null;

    function clearTimer() {
        if (dismissTimer !== null) {
            clearTimeout(dismissTimer);
            dismissTimer = null;
        }
    }

    function dismiss(): void {
        clearTimer();
        set(null);
    }

    function show(message: string, action?: Snackbar["action"], duration = 5000): void {
        clearTimer();
        set({ id: `snackbar-${Date.now()}`, message, type: "info", action });
        if (duration > 0) {
            dismissTimer = setTimeout(dismiss, duration);
        }
    }

    return {
        subscribe,
        show,
        dismiss,
        error: (message: string, onRetry?: () => void) => {
            clearTimer();
            const action = onRetry ? { label: "Réessayer", onClick: onRetry } : undefined;
            set({ id: `snackbar-${Date.now()}`, message, type: "error", action });
            // Errors with a retry action stay until dismissed manually; others auto-close
            if (!onRetry) {
                dismissTimer = setTimeout(dismiss, 8000);
            }
        },
    };
}

export const snackbar = createSnackbarStore();
