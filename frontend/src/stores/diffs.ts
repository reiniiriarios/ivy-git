import { get, writable } from "svelte/store";
import { parseResponse } from "scripts/parse-response";
import { GetUnstagedFileParsedDiff } from 'wailsjs/go/main/App';
import { changes } from "stores/changes";

interface Diff {
  Raw: string;
  Hunks: DiffHunk[];
  Binary: boolean;
  // set here
  File: string;
  Status: string;
}

interface DiffHunk {
  Header: string;
  StartOld: number;
  EndOld: number;
  StartNew: number;
  EndNew: number;
  Heading: string;
  Lines: DiffLine[];
}

interface DiffLine {
  Line: string;
  Type: string;
  RawLineNo: number;
  OldLineNo: number;
  NewLineNo: number;
  NoNewline: boolean;
}

function createUnstagedDiff() {
  const { subscribe, set } = writable({} as Diff);

  return {
    subscribe,
    set,
    fetch: async (file: string, status: string) => {
      GetUnstagedFileParsedDiff(file, status).then(result => {
        parseResponse(result, () => {
          result.Data.File = file;
          result.Data.Status = status;
          set(result.Data)
        }, () => set({} as Diff));
      });
    },
    clear: async () => set({} as Diff),
    refresh: async () => {
      // If the file that needs to be refreshed is in the changed files list, refresh it.
      if (get(changes).y.filter(c => c.File === get(unstagedFileDiff).File).length === 1) {
        unstagedFileDiff.fetch(get(unstagedFileDiff).File, get(unstagedFileDiff).Status);
      } else {
        unstagedFileDiff.clear();
      }
    },
  };
}
export const unstagedFileDiff = createUnstagedDiff();

export const fetchDiff = (file: string, status: string, staged: boolean, committed: boolean) => {
  if (!staged && !committed) {
    unstagedFileDiff.fetch(file, status);
  }
}
