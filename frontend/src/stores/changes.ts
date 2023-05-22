import { parseResponse } from 'scripts/parse-response';
import { get, writable } from 'svelte/store';
import { GitListChanges } from 'wailsjs/go/main/App';
import { currentRepo } from 'stores/repos';

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
      if (get(currentRepo)) {
        GitListChanges().then(result => {
          parseResponse(result, () => set({
            x: result.Data.ChangesX ?? [],
            y: result.Data.ChangesY ?? [],
          }));
        });
      } else {
        set({x: [], y: []});
      }
    },
  };
}
export const changes = createChanges();
