import { writable, type Writable } from "svelte/store";
import type { Case } from "$lib/types/case";

/**
 * Shared store for the currently active Case.
 * Written by the home page (via CasePicker), read by the layout/Footer.
 */
export const currentCase: Writable<Case | null> = writable(null);
