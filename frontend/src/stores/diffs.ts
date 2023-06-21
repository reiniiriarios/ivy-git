import { derived, get, writable } from "svelte/store";
import { changes } from "./changes";
import {
  GetCommitFileParsedDiff,
  GetHighlightedFileRange,
  ResolveDiffConflicts,
  StageFile,
  RemoveFile,
} from "wailsjs/go/main/App";
import { parseResponse } from "scripts/parse-response";

type OursTheirs = number;

export const DiffConflictType = {
  Ours: -1,
  Neither: 0,
  Theirs: 1,
  Both: 2,
  BothInverse: -2,
}

export interface Diff {
  File: string;
  Raw: string;
  Flags: string[];
  Hunks: DiffHunk[];
  Binary: boolean;
  Empty: boolean;
  SelectableLines: number;
  SelectedLines: number;
  NumConflicts: number;
  Conflicts: DiffConflict[];
  // UI
  TooLarge: boolean;
  Status: string;
  Staged: boolean;
  Committed: boolean;
  Conflict: boolean;
  Hash: string;
  Resolved: boolean;
  // Separate fetch
  Lang: string;
  Highlight: HighlightedLines;
}

export interface DiffHunk {
  Header: string;
  StartOld: number;
  EndOld: number;
  StartNew: number;
  EndNew: number;
  StartCur: number;
  EndCur: number;
  Heading: string;
  Lines: DiffLine[];
}

export interface DiffLine {
  Line: string;
  Type: string;
  RawLineNo: number;
  OldLineNo: number;
  NewLineNo: number;
  CurLineNo: number;
  NoNewline: boolean;
  MiniHunk: number;
  OursTheirs: OursTheirs;
  // UI
  Selected: boolean;
}

interface DiffConflict {
	Ours: DiffLine[];
	Theirs: DiffLine[];
	Resolution: OursTheirs;
}

type HighlightedLines = {[line: number]: string}[];

interface HighlightedFile {
	Size: number;
	TooLarge: boolean;
	Lang: string;
	Highlight: HighlightedLines;
}

function createCurrentDiff() {
  const { subscribe, set, update } = writable({} as Diff);

  return {
    subscribe,
    set,
    clear: () => set({} as Diff),
    // Fetch diff from specific commit.
    fetchDiff: async (hash: string, file: string, oldfile: string) => {
      GetCommitFileParsedDiff(hash, file, oldfile, false).then(result => {
        parseResponse(result, () => {
          let diff = {} as Diff;
          if (result.Response === 'too-large') {
            diff.TooLarge = true;
          } else {
            diff = result.Data;
          }
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
    // Set the resolution for a specific conflict.
    setConflictResolution: (conflict: number, resolution: number) => {
      update(d => {
        d.Conflicts[conflict].Resolution = resolution;
        return d;
      });
    },
    // Call the backend to edit the file, applying the conflict resolutions.
    resolveConflicts: async () => {
      // call backend
      ResolveDiffConflicts(get(currentDiff)).then(result => {
        parseResponse(result, () => {
          changes.refresh();
        });
      });
    },
    // Delete file relating to this diff.
    delete: async () => {
      RemoveFile(get(currentDiff).File).then(result => {
        parseResponse(result, () => {
          changes.refresh();
        });
      });
    },
    // Stage a conflicted file to keep it.
    keep: async () => {
      StageFile(get(currentDiff).File).then(result => {
        parseResponse(result, () => {
          changes.refresh();
        });
      });
    },
  }
}
export const currentDiff = createCurrentDiff();
export const currentDiffResolved = derived(currentDiff, $currentDiff =>
  $currentDiff?.Conflicts && Object.values($currentDiff.Conflicts).every(c => c.Resolution !== 0)
);
