import { get, writable } from "svelte/store";
import { changes } from "./changes";

export interface Diff {
  Raw: string;
  Hunks: DiffHunk[];
  Binary: boolean;
  SelectableLines: number;
  SelectedLines: number;
  // UI
  File: string;
  Status: string;
  Staged: boolean;
  Committed: boolean;
}

export interface DiffHunk {
  Header: string;
  StartOld: number;
  EndOld: number;
  StartNew: number;
  EndNew: number;
  Heading: string;
  Lines: DiffLine[];
}

export interface DiffLine {
  Line: string;
  Type: string;
  RawLineNo: number;
  OldLineNo: number;
  NewLineNo: number;
  NoNewline: boolean;
  MiniHunk: number;
  // UI
  Selected: boolean;
}

function createCurrentDiff() {
  const { subscribe, set } = writable({} as Diff);

  return {
    subscribe,
    set,
    clear: () => set({} as Diff),
    fetchDiff: async (hash: string, file: string) => {
      // Fetch diff from specific commit.
      //...
    },
    refresh: () => {
      let cd = get(currentDiff);
      if (cd.Committed) {
        // Handle diffs for previously committed files here.
        //...
      }
      else {
        // Let the changes store handle diffs for changed files.
        changes.fetchDiff(cd.Staged ? 'x' : 'y', cd.File);
      }
    },
  }
}
export const currentDiff = createCurrentDiff();
