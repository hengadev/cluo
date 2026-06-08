import { writable } from 'svelte/store';

/** Whether the manual modal UpdateDialog is open (settings, profile menu, etc.) */
export const updateDialogOpen = writable(false);
