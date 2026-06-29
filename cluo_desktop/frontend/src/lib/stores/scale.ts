import { writable } from 'svelte/store';
import { browser } from '$app/environment';

const STORAGE_KEY = 'ui-scale';
export const DEFAULT_SCALE = 1.0;
export const MIN_SCALE = 0.75;
export const MAX_SCALE = 1.25;

function createScaleStore() {
    const getInitialScale = (): number => {
        if (!browser) return DEFAULT_SCALE;
        const stored = localStorage.getItem(STORAGE_KEY);
        if (stored) {
            const parsed = parseFloat(stored);
            if (!isNaN(parsed) && parsed >= MIN_SCALE && parsed <= MAX_SCALE) return parsed;
        }
        return DEFAULT_SCALE;
    };

    const { subscribe, set } = writable<number>(getInitialScale());

    return {
        subscribe,
        set: (scale: number) => {
            if (!browser) return;
            document.documentElement.style.zoom = String(scale);
            localStorage.setItem(STORAGE_KEY, String(scale));
            set(scale);
        },
        init: () => {
            if (!browser) return;
            const scale = getInitialScale();
            document.documentElement.style.zoom = String(scale);
            set(scale);
        }
    };
}

export const uiScale = createScaleStore();
