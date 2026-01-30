import { writable } from 'svelte/store';
import { browser } from '$app/environment';

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
