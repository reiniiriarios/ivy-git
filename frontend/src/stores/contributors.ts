import { parseResponse } from "scripts/parse-response";
import { writable, get } from "svelte/store";
import { CommitsBehindMain, GetCachedContributorsData, ResetContributorsData, UpdateContributorsData } from "wailsjs/go/main/App";

interface Contributor {
	Name: string;
	Email: string;
	Commits: number;
	Insertions: number;
	Deletions: number;
}

interface Contributors {
  LastHashParsed: string;
  Contributors: Contributor[];
}

function createContributors() {
  const { subscribe, set } = writable({} as Contributors);

  return {
    subscribe,
    set,
    fetch: async () => {
      GetCachedContributorsData().then(result =>
        parseResponse(result, () => set(result.Data))
      );
    },
    update: async () => {
      await UpdateContributorsData().then(result =>
        parseResponse(result, () => set(result.Data))
      );
    },
    clear: async () => set({} as Contributors),
    reset: async () => {
      contributors.clear();
      ResetContributorsData();
    },
    numCommitsBehind: async (): Promise<number> => {
      let behind: number = 0;
      await CommitsBehindMain(get(contributors).LastHashParsed).then(result =>
        parseResponse(result, () => behind = result.Data)
      );
      return behind;
    },
  };
}
export const contributors = createContributors();
