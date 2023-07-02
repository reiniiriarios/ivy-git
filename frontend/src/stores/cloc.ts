import { parseResponse } from "scripts/parse-response";
import { get, writable } from "svelte/store";
import { CommitsBehindMain, GetCachedClocData, UpdateClocData, ResetClocData } from "wailsjs/go/main/App";

interface ClocData {
  LastHashParsed: string;
  Languages: LanguageData[];
  Total: LanguageData;
  Error: string;
}

interface LanguageData {
  Name: string;
  Files: number;
  Code: number;
  Comments: number;
  Blanks: number;
  Total: number;
  Bytes: number;
  CodePercent: number;
  TotalPercent: number;
}

function createCloc() {
  const { subscribe, set } = writable({} as ClocData);

  return {
    subscribe,
    fetch: async () => {
      set({} as ClocData);
      GetCachedClocData().then(result => {
        parseResponse(result, () => {
          set(result.Data);
        }, () => {
          set({ Error: result.Message } as ClocData);
        });
      });
    },
    update: async () => {
      await UpdateClocData().then(result =>
        parseResponse(result, () => set(result.Data))
      );
    },
    clear: async () => set({} as ClocData),
    reset: async () => {
      cloc.clear();
      ResetClocData();
    },
    numCommitsBehind: async (): Promise<number> => {
      let behind: number = 0;
      await CommitsBehindMain(get(cloc).LastHashParsed).then(result =>
        parseResponse(result, () => behind = result.Data)
      );
      return behind;
    },
  };
}
export const cloc = createCloc();
