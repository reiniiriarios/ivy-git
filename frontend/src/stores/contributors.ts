import { parseResponse } from "scripts/parse-response";
import { writable } from "svelte/store";
import { GetCachedContributorsData, ResetContributorsData, UpdateContributorsData } from "wailsjs/go/main/App";

interface Contributor {
	Name: string;
	Email: string;
	Commits: number;
	Insertions: number;
	Deletions: number;
}

type Contributors = Contributor[];

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
  };
}
export const contributors = createContributors();
