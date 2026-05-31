/**
 * Snackbar store for error feedback with retry actions.
 */
import { writable } from "svelte/store";

export interface Snackbar {
    id: string;
    message: string;
    action?: {
        label: string;
        onClick: () => void;
    };
}

function createSnackbarStore() {
    const { subscribe, set } = writable<Snackbar | null>(null);

    function show(message: string, action?: Snackbar["action"]): void {
        const id = `snackbar-${Date.now()}`;
        set({ id, message, action });
    }

    function dismiss(): void {
        set(null);
    }

    return {
        subscribe,
        show,
        dismiss,
        error: (message: string, onRetry?: () => void) => {
            show(message, onRetry ? { label: "Réessayer", onClick: onRetry } : undefined);
        },
    };
}

export const snackbar = createSnackbarStore();
