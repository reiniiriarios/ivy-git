import { parseResponse } from "scripts/parse-response";
import { RepoState, repoState } from "stores/repo-state";
import { get, writable } from "svelte/store";
import { FetchAll } from "wailsjs/go/main/App";

// How often to check the timer.
const AUTO_FETCH_INTERVAL = 1000; //ms
// How often to auto-fetch.
const AUTO_FETCH_TIMEOUT = 30000; //ms

export const autoFetchTimer = function () {
  const { subscribe, set, update } = writable(AUTO_FETCH_TIMEOUT);

  let interval: NodeJS.Timer;

  return {
    subscribe,
    init: async () => {
      interval = setInterval(() => {
        autoFetchTimer.tick().then(() => {
          if (get(autoFetchTimer) <= 0) {
            // Only fetch if in a neutral state and not in
            // the middle of a merge, rebase, bisect, etc.
            if (get(repoState) === RepoState.None) {
              console.log('Auto fetching...');
              FetchAll().then(r => {
                parseResponse(r, () => {
                  autoFetchTimer.reset();
                });
              });
              autoFetchTimer.reset();
            }
          }
        });
      }, AUTO_FETCH_INTERVAL);
    },
    disable: async () => clearInterval(interval),
    reset: async () => set(AUTO_FETCH_TIMEOUT),
    tick: async () => update(time => time -= 1000),
  };
}();
