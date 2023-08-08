import { get, writable } from 'svelte/store';
import { GetInProgressCommitMessageMerge, GetInProgressCommitMessageEdit, GetLastCommitMessage } from 'wailsjs/go/ivy/App'

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
    fetch: async () => {
      switch (get(repoState)) {
        // If the repo is in a state where there might be an in-progress message, then fetch it.
        case RepoState.Merge:
        case RepoState.RebaseMerge:
        case RepoState.Interactive:
        case RepoState.RevertSequence:
          commitMessage.fetchMerge();
          break;
        // Otherwise, clear it. This method should only run when the repo state changes, and so
        // this should ensure that any message loaded during a repo state change is cleared, as
        // it would no longer be relevant/correct to display.
        default:
          commitMessage.clear();
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
