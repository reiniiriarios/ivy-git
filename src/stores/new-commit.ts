import { writable } from 'svelte/store';

interface NewCommit {
  X: CommitChanges;
  Y: CommitChanges;
}

interface CommitChanges {
  Files: {[file: string]: string}[];
  Partials: Partial[];
}

interface Partial {
  File: string;
  //...
}

function createCommit() {
  const { subscribe, update } = writable({
    X: {
      Files: [],
      Partials: [],
    },
    Y: {
      Files: [],
      Partials: [],
    },
  } as NewCommit);
  
  return {
    subscribe,
    setFile: ((set: string, file: string, state: string) => {
      update(new_commit => {
        new_commit[set].Files[file] = state;
        return new_commit;
      });
    }),
  };
}
export const newCommit = createCommit();
