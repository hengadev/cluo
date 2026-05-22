import { writable, derived } from 'svelte/store';
import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import { apiFetch } from '$lib/services/apiFetch';
import type { AuthUser } from '$lib/types/entities';

interface AuthState {
	user: AuthUser | null;
	isAuthenticated: boolean;
}

function createAuthStore() {
	const getInitialState = (): AuthState => ({
		user: null,
		isAuthenticated: false
	});

	const { subscribe, set, update } = writable<AuthState>(getInitialState());

	return {
		subscribe,
		setUser: (user: AuthUser | null) => {
			update(state => ({
				...state,
				user,
				isAuthenticated: user !== null
			}));
		},
		clear: () => {
			set(getInitialState());
		},
		logout: async () => {
			try {
				await apiFetch('/auth/logout', { method: 'POST', skipRefresh: true });
			} catch {
				// Ignore logout errors - we're clearing local state anyway
			}
			set(getInitialState());
			if (browser) {
				await goto('/login');
			}
		},
		isAuthenticated: derived({ subscribe }, ($state) => $state.isAuthenticated)
	};
}

export const auth = createAuthStore();
