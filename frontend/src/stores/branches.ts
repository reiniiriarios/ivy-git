import { writable, get, derived } from 'svelte/store';

import { GetBranches, GetCurrentBranch, SwitchBranch, GetRemoteBranches } from 'wailsjs/go/main/App';

import { commitData, commitSignData } from 'stores/commits';
import { changes } from 'stores/changes';
import { currentTab } from 'stores/ui';
import { currentCommit } from 'stores/commit-details';
import { remoteData } from 'stores/remotes';
import { currentRepo } from 'stores/repos';
import { cloc } from 'stores/cloc';

import { parseResponse } from 'scripts/parse-response';

let cTab = '';
currentTab.subscribe(t => cTab = t);

export interface Branch {
  Name: string;
  Upstream: string;
  Remote: string;
  // ui
  Local: boolean;
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
export const upstreams = derived(branches, $branches => $branches.map(branch => branch.Upstream));

function createRemoteBranches() {
  const { subscribe, set } = writable([] as Branch[]);
  
  return {
    subscribe,
    refresh: async () => {
      GetRemoteBranches().then(result => {
        set(result.Data as Branch[]);
      });
    },
  };
}
export const remoteBranches = createRemoteBranches();
export const remoteOnlyBranches = derived([upstreams, remoteBranches], ([$upstreams, $remoteBranches]) => {
  return $remoteBranches.filter(branch => (
    !$upstreams.includes(branch.Remote+'/'+branch.Name)
  ));
});

function createCurrentBranch() {
  const { subscribe, set } = writable({} as Branch);
  
  return {
    subscribe,
    set,
    refresh: async () => {
      if (get(currentRepo)) {
        GetCurrentBranch().then(result => {
          // Ignore errors here. Sometimes there isn't a current branch.
          set(result.Data as Branch);
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
              currentCommit.clear();
            } else if (cTab === 'details') {
              remoteData.refresh();
              cloc.fetch();
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
    },
  };
}
export const currentBranch = createCurrentBranch();
export const detachedHead = derived(currentBranch, $currentBranch => $currentBranch?.Name && $currentBranch.Name === 'HEAD');
export const noBranchSelected = derived(currentBranch, $currentBranch => !$currentBranch?.Name?.length);
