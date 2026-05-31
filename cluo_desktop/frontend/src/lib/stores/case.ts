import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import type { CaseStatus } from '../types/entities';

type CaseContext = {
    id: string | null;
};

function createCaseStore() {
    const getInitialCase = (): CaseContext => {
        if (!browser) return { id: null };

        const stored = localStorage.getItem('currentCaseId');
        if (stored) return { id: stored };

        return { id: null };
    };

    const { subscribe, set, update } = writable<CaseContext>(getInitialCase());

    return {
        subscribe,
        setCase: (caseId: string) => {
            if (browser) {
                localStorage.setItem('currentCaseId', caseId);
            }
            set({ id: caseId });
        },
        clearCase: () => {
            if (browser) {
                localStorage.removeItem('currentCaseId');
            }
            set({ id: null });
        },
        getCaseId: () => {
            let caseId: string | null = null;
            subscribe(ctx => { caseId = ctx.id; })();
            return caseId;
        }
    };
}

export const currentCase = createCaseStore();

// =============================================================================
// RECENT CASES STORE
// =============================================================================

export interface RecentCaseEntry {
	id: string;
	title: string;
	status: CaseStatus;
}

const MAX_RECENT = 5;
const RECENT_KEY = 'recentCaseIds';

function createRecentCasesStore() {
	const getInitial = (): RecentCaseEntry[] => {
		if (!browser) return [];
		try {
			const stored = localStorage.getItem(RECENT_KEY);
			return stored ? JSON.parse(stored) : [];
		} catch {
			return [];
		}
	};

	const { subscribe, set } = writable<RecentCaseEntry[]>(getInitial());

	return {
		subscribe,
		push: (entry: RecentCaseEntry) => {
			const current = getInitial().filter(e => e.id !== entry.id);
			const next = [entry, ...current].slice(0, MAX_RECENT);
			if (browser) localStorage.setItem(RECENT_KEY, JSON.stringify(next));
			set(next);
		},
	};
}

export const recentCases = createRecentCasesStore();
