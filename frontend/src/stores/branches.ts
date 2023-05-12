import { writable } from 'svelte/store';
import { GetBranches, GetCurrentBranch, SwitchBranch } from '../../wailsjs/go/main/App';

interface Branch {
  Name: string;
}

function createBranches() {
  const { subscribe, update } = writable([] as Branch[]);
  
  return {
    subscribe,
    refresh: async () => {
      GetBranches().then(result => {
        update(_ => result as Branch[])
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
            if ((window as any).currentTab == 'tree') {
              (window as any).GetCommitList();
              (window as any).hideCommitDetails();
            }
            (window as any).getChanges();
            c = { Name: b };
          }
          return c;
        });
      });
    },
  };
}
export const currentBranch = createCurrentBranch();
