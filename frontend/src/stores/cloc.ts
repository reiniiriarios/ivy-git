import { parseResponse } from "scripts/parse-response";
import { get, writable } from "svelte/store";
import { CommitsBehindMain, GetCachedClocData, UpdateClocData, ResetClocData } from "wailsjs/go/main/App";
import { currentRepo } from "./repos";

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
      let repo = get(currentRepo);
      await UpdateClocData().then(result => {
        // This might take a while. If the repo isn't the same as when we started, don't update.
        if (get(currentRepo) !== repo) return;
        parseResponse(result, () => set(result.Data));
      });
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
export const clocRunning = writable(false);
