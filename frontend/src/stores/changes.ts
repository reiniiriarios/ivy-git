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
        if (result.Response === "error") {
          (window as any).messageModal(result.Message);
        } else {
          set({
            x: result.ChangesX ?? [],
            y: result.ChangesY ?? [],
          });
        }
      });      
    },
  };
}
export const changes = createChanges();