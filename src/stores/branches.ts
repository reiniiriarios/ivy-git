import { writable, get } from 'svelte/store';

import { GetBranches, GetCurrentBranch, SwitchBranch } from 'src/_tmp';

import { commitData, commitSignData } from 'stores/commit-data';
import { changes } from 'stores/changes';
import { currentTab } from 'stores/ui';
import { currentCommit } from 'stores/commit-details';
import { remoteData } from 'stores/remotes';
import { currentRepo } from 'stores/repos';

import { parseResponse } from 'scripts/parse-response';

let cTab = '';
currentTab.subscribe(t => cTab = t);

export interface Branch {
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
    set,
    refresh: async () => {
      if (get(currentRepo)) {
        GetCurrentBranch().then(result => {
          parseResponse(result, () => {
            set(result.Data as Branch);
          });
        });
      }
      else {
        set({} as Branch);
      }
    },
    switch: (b: string, r: string = "") => {
      if (b !== get(currentBranch)?.Name) {
        SwitchBranch(b, r).then(result => {
          parseResponse(result, () => {
            if (cTab === 'tree') {
              commitData.refresh();
              commitSignData.refresh();
              currentCommit.unset();
            } else if (cTab === 'details') {
              remoteData.refresh();
            }
            changes.refresh();
            set(result.Data);
          });
        });
      }
    },
    detach: () => {
      set({Name: 'HEAD'} as Branch);
    },
    clear: () => {
      set({} as Branch);
    }
  };
}
export const currentBranch = createCurrentBranch();
