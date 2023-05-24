import { parseResponse } from 'scripts/parse-response';
import { writable } from 'svelte/store';

interface File {
  File: string;
  Basename: string;
  Dir: string;
  Staged: boolean;
}

function createCurrentFile() {
  const { subscribe, set } = writable({} as File);
  
  return {
    subscribe,
    select: async (file: string, staged: boolean) => {
      // Set the minimum first to affect ui changes.
      set({
        File: file,
        Staged: staged,
      } as File);
      // Then fetch data.

    },
    reset: async () => {
      set({} as File);
    }
  };
}
export const currentFile = createCurrentFile();
