import { writable } from "svelte/store";

interface Conflicts {[file: string]: Conflict}[];

interface Conflict {
  ConflictSelections: number[];
}

function createConflicts() {
  const { subscribe, set, update } = writable({} as Conflicts);

  return {
    subscribe,
    clear: () => set({} as Conflicts),
    setFiles: (files: string[]) => {
      let conflicts: Conflicts = {};
      for (let i = 0; i < files.length; i++) {
        conflicts[files[i]] = {
          ConflictSelections: [],
        };
      }
      set(conflicts);
    },
    setConflictResolution(file: string, minihunk: number, oursTheirs: number) {
      update(c => {
        if (c[file]) {
          c[file].ConflictSelections[minihunk] = oursTheirs;
        }
        return c;
      });
    },
  }
}
export const conflicts = createConflicts();
