import { writable } from 'svelte/store';
import { GetBranches, GetCurrentBranch, SwitchBranch } from 'wailsjs/go/main/App';
import { commitData } from 'stores/commit-data';
import { changes } from 'stores/changes';
import { currentTab } from 'stores/current-tab';
import { currentCommit } from 'stores/commit-details';

let cTab = '';
currentTab.subscribe(t => cTab = t);

interface Branch {
  Name: string;
}

function createBranches() {
  const { subscribe, set } = writable([] as Branch[]);
  
  return {
    subscribe,
    refresh: async () => {
      GetBranches().then(result => {
        set(result as Branch[]);
      });
    },
  };
}
export const branches = createBranches();

function createCurrentBranch() {
  const { subscribe, update, set } = writable({} as Branch);
  
  return {
    subscribe,
    refresh: async () => {
      GetCurrentBranch().then(result => {
        set(result.Branch as Branch);
      });
    },
    set: (b: string) => {
      update(c => {
        if (c.Name === b) {
          return c;
        }
        SwitchBranch(b).then(result => {
          if (result.Response === "error") {
            (window as any).messageModal(result.Message);
          } else {
            if (cTab === 'tree') {
              commitData.refresh();
              currentCommit.unset();
            }
            changes.refresh();
            c = { Name: b };
          }
          return c;
        });
      });
    },
  };
}
export const currentBranch = createCurrentBranch();
