import { writable } from 'svelte/store';
import { AddRepo, GetRepos, GetSelectedRepo, RemoveRepo, UpdateSelectedRepo } from 'wailsjs/go/main/App';
import { commitData } from 'stores/commit-data';
import { changes } from 'stores/changes';
import { currentTab } from './current-tab';
import { currentCommit } from './commit-details';

export interface Repo {
  Name: string;
  Directory: string;
}

function createRepos() {
  const { subscribe, set } = writable([] as Repo[]);
  
  return {
    subscribe,
    refresh: async () => {
      GetRepos().then(result => {
        set(result as Repo[])
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
          repos.refresh();
        });
      }, 'Remove', 'Cancel');
    }
  };
}
export const repos = createRepos();

let cTab = '';
currentTab.subscribe(t => cTab = t);

function createCurrentRepo() {
  const { subscribe, update, set } = writable("");
  
  return {
    subscribe,
    refresh: async () => {
      GetSelectedRepo().then(result => {
        set(result);
      });
    },
    set: async (r: string) => {
      update(c => {
        if (c === r) {
          return c;
        }
        UpdateSelectedRepo(r).then(() => {
          if (cTab === 'tree') {
            commitData.refresh();
            currentCommit.unset();
          }
          changes.refresh();
        });
        return r;
      });
    },
  };
}
export const currentRepo = createCurrentRepo();
