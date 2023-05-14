import { parseResponse } from 'scripts/parse-response';
import { writable } from 'svelte/store';
import { GitListChanges } from 'wailsjs/go/main/App';

interface Change {
  File: string;
  Basename: string;
  Dir: string;
  Letter: string;
  Flag: string;
}

function createChanges() {
  const { subscribe, set } = writable({
    x: [] as Change[],
    y: [] as Change[],
  });
  
  return {
    subscribe,
    refresh: async () => {
      GitListChanges().then(result => {
        parseResponse(result, () => set({
          x: result.ChangesX ?? [],
          y: result.ChangesY ?? [],
        }));
      });      
    },
  };
}
export const changes = createChanges();
