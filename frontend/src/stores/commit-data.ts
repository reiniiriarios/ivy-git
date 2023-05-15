import { removeTableSizing, updateSavedTableSizing } from 'scripts/commit-table-resize';
import { drawGraph, getSVGWidth } from 'scripts/graph';
import { parseResponse } from 'scripts/parse-response';
import { derived, writable } from 'svelte/store';
import { GetCommitList } from 'wailsjs/go/main/App';

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
  Remotes: Ref[];
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
  ShortName: string;
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
  NextParent: number;
  BranchId: number;
  X: number;
  XNext: number;
  Connections: Connection[];
  Committed: boolean;
  Stash: boolean;
}

export interface Branch {
  Id: number;
  Color: number;
  Lines: Line[];
  UncommitedPoints: number;
  Merge: boolean;
}

export interface Graph {
  Vertices: Vertex[];
  Branches: Branch[];
  Width: number;
  Height: number;
}

function createCommitData() {
  const { subscribe, set } = writable({
    commits: [] as Commit[],
    HEAD: {} as Ref,
    Graph: {} as Graph,
  });
  
  return {
    subscribe,
    refresh: async () => {
      GetCommitList().then(result => {
        console.log('um');
        parseResponse(result, () => {
          removeTableSizing().then(() => {
            set({
              commits: result.Commits,
              HEAD: result.HEAD,
              Graph: result.Graph,
            });
            console.log('HEAD', result.HEAD);
            console.log('commits', result.commits);
            console.log('branches', result.Graph.Branches);
            console.log('vertices', result.Graph.Vertices);
          });
        });
      });
    },
  };
}
export const commitData = createCommitData();
export const commits = derived(commitData, $commitData => $commitData.commits);
export const HEAD = derived(commitData, $commitData => $commitData.HEAD);
export const graph = derived(commitData, $commitData => $commitData.Graph);
export const tree = derived(commitData, $commitData => ({
  svg: drawGraph($commitData.Graph),
  width: getSVGWidth($commitData.Graph),
}));
