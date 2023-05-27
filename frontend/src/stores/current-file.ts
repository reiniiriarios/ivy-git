import { writable } from 'svelte/store';
import { fetchDiff, unstagedFileDiff } from 'stores/diffs';

interface File {
  File: string;
  Basename: string;
  Dir: string;
  Staged: boolean;
  Status: string;
  Committed: boolean;
}

function createCurrentFile() {
  const { subscribe, set } = writable({} as File);
  
  return {
    subscribe,
    select: async (file: string, status: string, staged: boolean, committed: boolean) => {
      // Set the minimum first to affect ui changes.
      set({
        File: file,
        Staged: staged,
        Status: status,
        Committed: committed,
      } as File);
      // Then fetch data.
      fetchDiff(file, status, staged, committed);
    },
    clear: async () => {
      set({} as File);
      unstagedFileDiff.clear();
    }
  };
}
export const currentFile = createCurrentFile();
