import { writable } from 'svelte/store';

// These stores reflec the current ui state and can be used
// across the app to change the ui state from components
// unrelated in hierarchy, but related in content.
export const repoSelect = writable(false);
export const branchSelect = writable(false);
