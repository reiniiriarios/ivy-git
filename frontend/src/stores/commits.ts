import { drawGraph, getSVGWidth } from 'scripts/graph';
import { parseResponse } from 'scripts/parse-response';
import { derived, writable, get } from 'svelte/store';
import { GetCommitList, GetCommitsSignStatus } from 'wailsjs/go/ivy/App';
import { settings } from 'stores/settings';

const COMMIT_LIST_PAGING = 500;

export interface Commit {
  Id: number;
  Hash: string;
  Parents: string[];
	RefName: string;
  AuthorName: string;
  AuthorEmail: string;
  AuthorTimestamp: number;
  AuthorDatetime: string;
  Subject: string;
  Branches: Ref[];
  Tags: Ref[];
  RemoteBranches: Ref[];
  Heads: Ref[];
  Stash: boolean;
  Labeled: boolean;
  Color: number;
  X: number;
  Merge: boolean;
}

export interface Ref {
  Hash: string;
  Name: string;
  Branch: string;
  Remote: string;
  Upstream: string;
  AbbrName: string;
  SyncedRemotes: string[];
  SyncedLocally: boolean;
  Head: boolean;
}

interface Line {
  P1: Point;
  P2: Point;
  Committed: boolean;
  // true = P1, false = P2
  LockedDirection: boolean;
}

interface Point {
  X: number;
  Y: number;
}

interface Connection {
  VertexId: number;
  BranchId: number;
}

export interface Vertex {
  Id: number;
  Children: Vertex[];
  Parents: Vertex[];
  BranchId: number;
  X: number;
  Connections: Connection[];
  Committed: boolean;
  Stash: boolean;
}

export interface Limb {
  Id: number;
  Color: number;
  Lines: Line[];
  UncommitedPoints: number;
  Merge: boolean;
}

export interface Graph {
  Vertices: Vertex[];
  Limbs: Limb[];
  Width: number;
  Height: number;
  Continues: boolean;
}

interface CommitData {
  commits: Commit[];
  HEAD: Ref;
  Graph: Graph;
  page: number;
}

function createCommitData() {
  const { subscribe, set } = writable({} as CommitData);
  
  return {
    subscribe,
    refresh: async (page: number = 0) => {
      // todo: page instead of count
      GetCommitList(COMMIT_LIST_PAGING * (page + 1), 0).then(result => {
        parseResponse(result, () => {
          if (result.Response === 'no-commits') {
            set({
              commits: [] as Commit[],
              HEAD: {} as Ref,
              Graph: {} as Graph,
              page: 0,
            });
          }
          else {
            set({
              commits: result.Data.Commits,
              HEAD: result.Data.HEAD,
              Graph: result.Data.Graph,
              page: page,
            });
          }
        }, () => {}, true);
      });
      commitSignData.refresh();
    },
    clear: async () => set({} as CommitData),
  };
}
export const commitData = createCommitData();
export const commits = derived(commitData, $commitData => $commitData?.commits);
export const commitsMap = derived(commitData, $commitData => {
  let map = new Map<string, number>();
  if ($commitData?.commits && $commitData.commits.length) {
    $commitData.commits.forEach((c, i) => map.set(c.Hash, i));
  }
  return map;
});
export const HEAD = derived(commitData, $commitData => $commitData?.HEAD);
export const graph = derived(commitData, $commitData => $commitData?.Graph);
export const commitsPage = derived(commitData, $commitData => $commitData?.page);
export const tree = derived(commitData, $commitData => ({
  svg: $commitData?.Graph ? drawGraph($commitData.Graph) : null,
  width: $commitData?.Graph ? getSVGWidth($commitData.Graph) : '0',
  height: $commitData?.Graph?.Height,
  continues: $commitData?.Graph?.Continues,
}));

export type CommitsSigned = {[hash: string]: string};

function createCommitSignData() {
  const { subscribe, set, update } = writable({
    commits: {} as CommitsSigned,
    page: 0,
  });
  
  return {
    subscribe,
    refresh: async (page: number = 0) => {
      if (!get(settings).DisplayCommitSignatureInList) {
        return;
      }
      set({
        commits: {},
        page: 0,
      });
      let miniPaging = 10;
      // todo: page instead of count
      for (let i = 0; i < COMMIT_LIST_PAGING; i += miniPaging) {
        await GetCommitsSignStatus(miniPaging * (page + 1), i * (page + 1)).then(result => {
          parseResponse(result, () => {
            update(cs => {
              return {
                commits: {...cs.commits, ...result.Data},
                page: page,
              }
            });
          });
        });
      }
    },
  };
}
export const commitSignData = createCommitSignData();
