import { writable } from 'svelte/store';
import { AddRepo, GetRepos, GetSelectedRepo, RemoveRepo, UpdateSelectedRepo } from 'wailsjs/go/main/App';
import { commitData } from 'stores/commit-data';
import { changes } from 'stores/changes';
import { currentTab } from 'stores/current-tab';
import { currentCommit } from 'stores/commit-details';
import { branches, currentBranch } from 'stores/branches';
import { remotes } from 'stores/remotes';
import { parseResponse } from 'scripts/parse-response';

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
      AddRepo().then((result) => parseResponse(result, repos.refresh));
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
        currentBranch.refresh();
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
          } else if (cTab === 'details') {
            remotes.refresh();
          }
          branches.refresh();
          changes.refresh();
        });
        return r;
      });
    },
  };
}
export const currentRepo = createCurrentRepo();
