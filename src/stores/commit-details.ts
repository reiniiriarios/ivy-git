import { writable } from "svelte/store";
import { GetCommitDetails, GetCommitDiffSummary, GetCommitSignature } from "src/_tmp";
import type { Commit } from "stores/commit-data";
import { parseResponse } from "scripts/parse-response";

interface CommitDetails {
  Hash: string;
  Body: string;
  BodyHtml: string;
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
  OldRel: string;
  Added: number;
  Deleted: number;
  Binary: boolean;
  Status: string;
}

export interface FileStatDir {
  Name: string;
  Path: string[];
  Files: FileStat[];
  Dirs: FileStatDir[];
}

interface CommitSignature {
  Status: string;
  Name: string;
  Key: string;
  Description: string;
}

function createCurrentCommit() {
  const { subscribe, update, set } = writable({} as Commit);

  return {
    subscribe,
    toggle: (commit: Commit) =>
      update((c) => {
        if (commit.Hash === c.Hash) {
          commitDetails.set({} as CommitDetails);
          commitDiffSummary.set({} as FileStatDir);
          commitSignature.set({} as CommitSignature);
          return {} as Commit;
        }
        commitDetails.fetch(commit.Hash);
        // Clear first, wait for data to display.
        commitDiffSummary.set({} as FileStatDir);
        commitDiffSummary.fetch(commit.Stash ? commit.Parents[1] : commit.Hash);
        // Clear first, wait for data to display.
        commitSignature.set({} as CommitSignature);
        commitSignature.fetch(commit.Hash);
        return commit;
      }),
    unset: () => {
      commitDetails.set({} as CommitDetails);
      commitDiffSummary.set({} as FileStatDir);
      commitSignature.set({} as CommitSignature);
      set({} as Commit);
    },
  };
}
export const currentCommit = createCurrentCommit();

function createDetails() {
  const { subscribe, set } = writable({} as CommitDetails);

  return {
    subscribe,
    set,
    fetch: async (hash: string) => {
      GetCommitDetails(hash).then(result =>
        parseResponse(result, () => set(result.Data), () => set({} as CommitDetails))
      );
    },
  };
}
export const commitDetails = createDetails();

function createSummary() {
  const { subscribe, set } = writable({} as FileStatDir);

  return {
    subscribe,
    set,
    fetch: async (hash: string) => {
      GetCommitDiffSummary(hash).then(result =>
        parseResponse(result, () => set(result.Data), () => set({} as FileStatDir))
      );
    },
  };
}
export const commitDiffSummary = createSummary();

function createSignData() {
  const { subscribe, set } = writable({} as CommitSignature);

  return {
    subscribe,
    set,
    fetch: async (hash: string) => {
      GetCommitSignature(hash).then(result => {
        parseResponse(result, () => set(result.Data), () => set({} as CommitSignature))
      }
      );
    },
  };
}
export const commitSignature = createSignData();
