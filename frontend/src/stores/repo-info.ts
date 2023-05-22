import { parseResponse } from "scripts/parse-response";
import { writable } from "svelte/store";
import { NumBranches, NumMainBranchCommits, NumTags } from "wailsjs/go/main/App";

function createNumBranches() {
  const { subscribe, set } = writable(0);

  return {
    subscribe,
    fetch: async () => {
      NumBranches().then(result => {
        set(result.Data);
      });
    },
  };
}
export const numBranches = createNumBranches();

function createNumTags() {
  const { subscribe, set } = writable(0);

  return {
    subscribe,
    fetch: async () => {
      NumTags().then(result => {
        set(result.Data);
      });
    },
  };
}
export const numTags = createNumTags();

function createMainInfo() {
  const { subscribe, set } = writable({
    Count: 0,
    Name: '',
  });

  return {
    subscribe,
    fetch: async () => {
      NumMainBranchCommits().then(result => {
        console.log(result)
        parseResponse(result, () => {
          set(result.Data);
        });
      });
    },
  };
}
export const mainInfo = createMainInfo();
