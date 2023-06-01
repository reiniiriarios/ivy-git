import { writable, get } from 'svelte/store';

import { AddRepo, CreateRepo, GetRepos, GetSelectedRepo, RemoveRepo, UpdateMain, UpdateSelectedRepo } from 'wailsjs/go/main/App';

import { commitData, commitSignData } from 'stores/commit-data';
import { changes } from 'stores/changes';
import { currentCommit } from 'stores/commit-details';
import { branches, currentBranch } from 'stores/branches';
import { remoteData } from 'stores/remotes';
import { currentTab, repoSelect } from 'stores/ui';
import { messageDialog } from 'stores/message-dialog';
import { cloc, numBranches, numCommits, numTags } from 'stores/repo-info';
import { currentDiff } from 'stores/diffs';

import { parseResponse } from 'scripts/parse-response';

export interface Repo {
  Name: string;
  Directory: string;
  Main: string;
}

let cTab = '';
currentTab.subscribe(t => cTab = t);

function createRepos() {
  const { subscribe, set, update } = writable([] as Repo[]);
  
  return {
    subscribe,
    refresh: async () => {
      GetRepos().then((result: Repo[]) => {
        // Sort
        if (Object.entries(result).length) {
          set(Object.entries(result).sort(([_aId, aRepo], [_bId, bRepo]) => {
            if (aRepo.Name < bRepo.Name) return -1;
            if (aRepo.Name > bRepo.Name) return 1;
            return 0;
          }).reduce(
            (obj, [id, repo]) => { 
              obj[id] = repo; 
              return obj;
            }, 
            {} as Repo[]
          ));
        }
      });
    },
    add: async () => {
      AddRepo().then((result) => parseResponse(result, () => {
        if (result.Response !== 'none') {
          repos.refresh();
          currentRepo.switch(result.Id);
          repoSelect.set(false);
          messageDialog.clear();
        }
      }));
    },
    create: async () => {
      messageDialog.addRepo({
        callbackConfirm: () => {
          let data = messageDialog.addRepoData();
          CreateRepo(data.name, data.location).then(result => {
            parseResponse(result, () => {
              repos.refresh();
              currentRepo.switch(result.Id);
              repoSelect.set(false);
              messageDialog.clear();
            });
          });
        },
      });
    },
    delete: async (id: string) => {
      let name = get(repos)[id]?.Name ?? 'this repo';
      messageDialog.confirm({
        heading: 'Remove Repo',
        message: `Are you sure you want to remove <strong>${name}</strong>?<br><br>This will not affect the repo or its files except to remove it from this app.`,
        confirm: 'Remove',
        okay: 'Cancel',
        callbackConfirm: () => {
          RemoveRepo(id).then(() => {
            repos.refresh();
            if (id === get(currentRepo)) {
              currentRepo.clear();
            }
          })
        },
      });
    },
    updateMain: async (branch: string) => {
      UpdateMain(branch).then(r => {
        update(repos => {
          repos[get(currentRepo)].Main = branch;
          return repos;
        });
        if (cTab === 'details') {
          remoteData.refresh();
          numCommits.fetch();
          cloc.fetch();
        }
      });
    }
  };
}
export const repos = createRepos();

function createCurrentRepo() {
  const { subscribe, set } = writable("");
  
  return {
    subscribe,
    refresh: async () => {
      GetSelectedRepo().then((result: string) => {
        set(result);
        currentBranch.refresh();
        branches.refresh();
        changes.refresh();
      });
    },
    clear: async () => {
      UpdateSelectedRepo("").then(result => {
        parseResponse(result, () => {
          if (cTab === 'tree') {
            commitData.refresh();
            commitSignData.refresh();
            currentCommit.unset();
          } else if (cTab === 'details') {
            remoteData.refresh();
            numCommits.fetch();
            numBranches.fetch();
            numTags.fetch();
            cloc.fetch();
          }
          currentDiff.clear();
          branches.refresh();
          currentBranch.refresh();
          changes.refresh();
          set("");
        });
      });
    },
    switch: async (repo_id: string) => {
      let current_repo_id = get(currentRepo);
      if (current_repo_id === repo_id) {
        return;
      }
      UpdateSelectedRepo(repo_id).then(result => {
        parseResponse(result, () => {
          if (cTab === 'tree') {
            commitData.refresh();
            commitSignData.refresh();
            currentCommit.unset();
          } else if (cTab === 'details') {
            remoteData.refresh();
            numCommits.fetch();
            numBranches.fetch();
            numTags.fetch();
            cloc.fetch();
          }
          currentDiff.clear();
          branches.refresh();
          currentBranch.refresh();
          changes.refresh();
          set(repo_id);
        });
      });
    },
  };
}
export const currentRepo = createCurrentRepo();
