import { get, writable } from "svelte/store";
import { changes } from "./changes";
import { GetCommitFileParsedDiff } from "wailsjs/go/main/App";
import { parseResponse } from "scripts/parse-response";

export const DiffConflict = {
  Ours: -1,
  Neither: 0,
  Theirs: 1,
  Both: 2,
  BothInverse: 3,
}

export interface Diff {
  Raw: string;
  Hunks: DiffHunk[];
  Binary: boolean;
  SelectableLines: number;
  SelectedLines: number;
  NumConflicts: number;
  // UI
  File: string;
  Status: string;
  Staged: boolean;
  Committed: boolean;
  Conflict: boolean;
  Hash: string;
  ConflictSelections: number[];
  Resolved: boolean;
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
  OursTheirs: number; // -1 ours, 0 neither, 1 theirs
  // UI
  Selected: boolean;
}

function createCurrentDiff() {
  const { subscribe, set, update } = writable({} as Diff);

  return {
    subscribe,
    set,
    clear: () => set({} as Diff),
    // Fetch diff from specific commit.
    fetchDiff: async (hash: string, file: string, oldfile: string) => {
      GetCommitFileParsedDiff(hash, file, oldfile).then(result => {
        parseResponse(result, () => {
          let diff = result.Data;
          diff.Staged = false;
          diff.Committed = true;
          diff.Hash = hash;
          diff.File = oldfile ? `${file} -> ${oldfile}` : file; // ???
          set(diff);
        });
      });
    },
    // Refetch the current diff.
    refresh: () => {
      let cd = get(currentDiff);
      if (!cd.Committed) {
        // Let the changes store handle diffs for changed files.
        changes.fetchDiff(cd.Staged ? 'x' : 'y', cd.File, true);
      }
    },
    setConflictResolution: (conflict: number, resolution: number) => {
      update(d => {
        d.ConflictSelections[conflict] = resolution;
        changes.setResolved(d.File, d.Resolved);
        return d;
      });
    },
    setResolved: (resolved: boolean) => {
      update(d => {
        d.Resolved = resolved;
        changes.setResolved(d.File, resolved);
        return d;
      })
    },
  }
}
export const currentDiff = createCurrentDiff();
