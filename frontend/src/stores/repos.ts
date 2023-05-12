import { writable } from 'svelte/store';
import { AddRepo, GetRepos, GetSelectedRepo, RemoveRepo, UpdateSelectedRepo } from '../../wailsjs/go/main/App';

export interface Repo {
  Name: string;
  Directory: string;
}

function createRepos() {
  const { subscribe, update } = writable([] as Repo[]);
  
  return {
    subscribe,
    refresh: async () => {
      GetRepos().then(result => {
        update(_ => result as Repo[])
      });
    },
    add: async () => {
      AddRepo().then((result) => {
        if (result.Response === "error") {
          (window as any).messageModal(result.Message);
        } else {
          repos.refresh();
        }
      });
    },
    delete: async (id: string) => {
      (window as any).confirmModal(`Are you sure you want to remove ${repos[id].Name}?`, () => {
        RemoveRepo(id).then(() => {
          repos.refresh()
        });
      }, 'Remove', 'Cancel');
    }
  };
}
export const repos = createRepos();

function createCurrentRepo() {
  const { subscribe, update } = writable("");
  
  return {
    subscribe,
    refresh: async () => {
      GetSelectedRepo().then(result => {
        update(_ => result)
      });
    },
    set: (r: string) => {
      update(c => {
        if (c !== r) {
          UpdateSelectedRepo(r).then(() => {
            if ((window as any).currentTab == 'tree') {
              (window as any).GetCommitList();
              (window as any).hideCommitDetails();
            }
          });
          c = r;
        }
        return c;
      });
    },
  };
}
export const currentRepo = createCurrentRepo();
