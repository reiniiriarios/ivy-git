import { derived, get, writable } from 'svelte/store';
import { GetInProgressCommitMessageMerge, GetInProgressCommitMessageEdit } from 'wailsjs/go/main/App'

import { RepoState, repoState } from 'stores/repo-state';
import { parseResponse } from 'scripts/parse-response';

// These stores reflec the current ui state and can be used
// across the app to change the ui state from components
// unrelated in hierarchy, but related in content.
export const currentTab = writable('tree');
export const repoSelect = writable(false);
export const branchSelect = writable(false);
export const commitDetailsWindow = writable(false);

interface CommitMessage {
  Subject: string;
  Body: string;
}

function createInProgMsg() {
  const { subscribe, set } = writable({} as CommitMessage);

  return {
    subscribe,
    // If the repo is in a state where there might be an in-progress message, then fetch it.
    fetch: async () => {
      switch (get(repoState)) {
        case RepoState.Merge:
        case RepoState.Rebase:
        case RepoState.RebaseMerge:
        case RepoState.Apply:
        case RepoState.ApplyOrRebase:
        case RepoState.Interactive:
        case RepoState.Revert:
        case RepoState.RevertSequence:
        case RepoState.CherryPick:
        case RepoState.CherryPickSequence:
          inProgressCommitMessage.fetchMerge();
      }
    },
    fetchMerge: async () => {
      GetInProgressCommitMessageMerge().then(result => {
        parseResponse(result, () => {
          set(result.Data);
        });
      });
    },
    fetchEdit: async () => {
      GetInProgressCommitMessageEdit().then(result => {
        parseResponse(result, () => {
          set(result.Data);
        });
      });
    },
    // Check if there's anything currently loaded or typed into the commit message fields.
    check: async () => {
      if (!get(commitMessageSubject) && !get(commitMessageBody)) {
        inProgressCommitMessage.fetch();
      }
    },
    clear: async () => {
      set({} as CommitMessage);
    },
    refresh: async () => {
      inProgressCommitMessage.clear().then(() => {
        inProgressCommitMessage.fetch();
      });
    }
  };
}
export const inProgressCommitMessage = createInProgMsg();
export const commitMessageSubject = derived(inProgressCommitMessage, $inProgressCommitMessage => $inProgressCommitMessage.Subject);
export const commitMessageBody = derived(inProgressCommitMessage, $inProgressCommitMessage => $inProgressCommitMessage.Body);
