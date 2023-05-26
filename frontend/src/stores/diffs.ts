import { writable } from "svelte/store";
import { parseResponse } from "scripts/parse-response";
import { GetUnstagedFileParsedDiff } from 'wailsjs/go/main/App';

interface Diff {
  Raw: string;
  Hunks: DiffHunk[];
  Binary: boolean;
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
      console.log('fetching', file, status)
      GetUnstagedFileParsedDiff(file, status).then(result => {
        console.log(result)
        parseResponse(result, () => set(result.Data), () => set({} as Diff));
    });
    },
  };
}
export const unstagedFileDiff = createUnstagedDiff();

export const fetchDiff = (file: string, status: string, staged: boolean, committed: boolean) => {
  if (!staged && !committed) {
    unstagedFileDiff.fetch(file, status);
  }
}
