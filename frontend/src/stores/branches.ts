import { writable, get } from 'svelte/store';

import { GetBranches, GetCurrentBranch, SwitchBranch } from 'wailsjs/go/main/App';

import { commitData, commitSignData } from 'stores/commit-data';
import { changes } from 'stores/changes';
import { currentTab } from 'stores/ui';
import { currentCommit } from 'stores/commit-details';
import { remotes } from 'stores/remotes';

import { parseResponse } from 'scripts/parse-response';

let cTab = '';
currentTab.subscribe(t => cTab = t);

interface Branch {
  Name: string;
  Upstream: string;
}

function createBranches() {
  const { subscribe, set } = writable([] as Branch[]);
  
  return {
    subscribe,
    refresh: async () => {
      GetBranches().then(result => {
        set(result.Data as Branch[]);
      });
    },
  };
}
export const branches = createBranches();

function createCurrentBranch() {
  const { subscribe, set } = writable({} as Branch);
  
  return {
    subscribe,
    refresh: async () => {
      GetCurrentBranch().then(result => {
        set(result.Data as Branch);
      });
    },
    set: (b: string) => {
      if (b !== get(currentBranch).Name) {
        SwitchBranch(b).then(result => {
          parseResponse(result, () => {
            if (cTab === 'tree') {
              commitData.refresh();
              commitSignData.refresh();
              currentCommit.unset();
            } else if (cTab === 'details') {
              remotes.refresh();
            }
            changes.refresh();
            set(result.Data);
          });
        });
      }
    },
  };
}
export const currentBranch = createCurrentBranch();
