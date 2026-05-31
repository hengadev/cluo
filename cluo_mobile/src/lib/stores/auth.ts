import { writable } from 'svelte/store';
import { goto } from '$app/navigation';
import { browser } from '$app/environment';

export interface AuthUser {
    id: string;
    email: string;
    role: 'admin' | 'investigator' | 'viewer';
    name?: string;
}

interface AuthState {
    user: AuthUser | null;
    isAuthenticated: boolean;
}

function createAuthStore() {
    const { subscribe, set, update } = writable<AuthState>({ user: null, isAuthenticated: false });

    return {
        subscribe,
        setUser: (user: AuthUser | null) => {
            update(state => ({ ...state, user, isAuthenticated: user !== null }));
        },
        clear: () => set({ user: null, isAuthenticated: false }),
        logout: async () => {
            set({ user: null, isAuthenticated: false });
            if (browser) await goto('/auth');
        }
    };
}

export const auth = createAuthStore();
