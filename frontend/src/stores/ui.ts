import { writable } from 'svelte/store';
import { GetInProgressCommitMessageEither } from 'wailsjs/go/main/App'

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
    fetch: async () => {
      GetInProgressCommitMessageEither().then(result => { console.log(result); set(result) });
    },
    clear: async () => {
      set({} as CommitMessage);
    }
  };
}
export const inProgressCommitMessage = createInProgMsg();
