import { get, writable } from 'svelte/store';
import { GetInProgressCommitMessageMerge, GetInProgressCommitMessageEdit, GetLastCommitMessage } from 'wailsjs/go/main/App'

import { RepoState, repoState } from 'stores/repo-state';
import { parseResponse } from 'scripts/parse-response';

// These stores reflec the current ui state and can be used
// across the app to change the ui state from components
// unrelated in hierarchy, but related in content.
export const currentTab = writable('tree');
export const repoSelect = writable(false);
export const branchSelect = writable(false);
export const commitDetailsWindow = writable(false);
export const linkPreviewHref = writable('');

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
        case RepoState.RebaseMerge:
        case RepoState.Interactive:
        case RepoState.RevertSequence:
          commitMessage.fetchMerge();
      }
    },
    fetchMerge: async () => {
      GetInProgressCommitMessageMerge().then(result => {
        parseResponse(result, () => {
          set(result.Data);
          commitMessageSubject.set(result.Data.Subject ?? "");
          commitMessageBody.set(result.Data.Body ?? "");
        });
      });
    },
    fetchEdit: async () => {
      GetInProgressCommitMessageEdit().then(result => {
        parseResponse(result, () => {
          set(result.Data);
          commitMessageSubject.set("");
          commitMessageBody.set("");
        });
      });
    },
    fetchLast: async () => {
      GetLastCommitMessage().then(result => {
        parseResponse(result, () => {
          commitMessageSubject.set(result.Data.Subject ?? "");
          commitMessageBody.set(result.Data.Body ?? "");
        });
      });
    },
    // Check if there's anything currently loaded or typed into the commit message fields.
    check: async () => {
      if (!get(commitMessageSubject) && !get(commitMessageBody)) {
        commitMessage.fetch();
      }
    },
    clear: async () => {
      set({} as CommitMessage);
      commitMessageSubject.set("");
      commitMessageBody.set("");
    },
    refresh: async () => {
      commitMessage.clear().then(() => {
        commitMessage.fetch();
      });
    }
  };
}
export const commitMessage = createInProgMsg();
// These two must be writable strings and not derived in order to bind them.
export const commitMessageSubject = writable("");
export const commitMessageBody = writable("");
