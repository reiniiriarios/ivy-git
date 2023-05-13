import { writable } from 'svelte/store';
import { GetCommitDetails, GetCommitDiffSummary } from 'wailsjs/go/main/App';
import type { Commit } from 'stores/commit-data';

interface CommitDetails {
  Hash: string;
  Body: string;
  CommitterName: string;
  CommitterEmail: string;
  CommitterTimestamp: number;
  CommitterDatetime: string;
}

interface FileStat {
  File: string;
  Name: string;
  Dir: string;
  Path: string[];
  OldFile: string;
  OldName: string;
  OldDir: string;
  Added: number;
  Deleted: number;
  Status: string;
}

export interface FileStatDir {
  Name: string;
  Path: string[];
  Files: FileStat[];
  Dirs: FileStatDir[];
}

function createCurrentCommit() {
  const { subscribe, update, set } = writable({} as Commit);

  return {
    subscribe,
    toggle: (commit: Commit) => update(c => {
      if (commit.Hash === c.Hash) {
        commitDetails.set({} as CommitDetails);
        commitDiffSummary.set({} as FileStatDir);
        return {} as Commit;
      }
      commitDetails.fetch(commit.Hash);
      commitDiffSummary.fetch(commit.Hash);
      return commit;
    }),
    unset: () => {
      commitDetails.set({} as CommitDetails);
      commitDiffSummary.set({} as FileStatDir);
      set({} as Commit);
    },
  }
}
export const currentCommit = createCurrentCommit();

function createDetails() {
  const { subscribe, set } = writable({} as CommitDetails);

  return {
    subscribe,
    set,
    fetch: async (hash: string) => {
      GetCommitDetails(hash).then(result => {
        if (result.Response === "error") {
          (window as any).messageModal(result.Message);
          set({} as CommitDetails);
        }
        set(result.Commit);
      });
    }
  }
}
export const commitDetails = createDetails();

function createSummary() {
  const { subscribe, set } = writable({} as FileStatDir);

  return {
    subscribe,
    set,
    fetch: async (hash: string) => {
      GetCommitDiffSummary(hash).then(result => {
        if (result.Response === "error") {
          (window as any).messageModal(result.Message);
          set({} as FileStatDir);
        }
        set(result.Files);
      });
    }
  }
}
export const commitDiffSummary = createSummary();