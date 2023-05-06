// Match to git.commits.go.
export const UNCOMMITED_HASH = "#";

// Match to _tree.scss.
const NUM_COLORS = 10;

// Adjust start of graph from top left.
const OFFSET_X = 12;
const OFFSET_Y = 12;

// Scale from graph coordinates to pixels.
const SCALE_X = 12;
// SCALE_Y should match height of <tr>.
const SCALE_Y = 24;

// Adjust curve of lines.
// SCALE_Y * 0  = straight lines, hard corners
// SCALE_Y * 1  = curved, but hits horizontal
// SCALE_Y * 1+ = inverts curve
const CURVE = SCALE_Y * 0.5;

// Dot size.
const VERTEX_RADIUS = 3;
const VERTEX_RADIUS_U = 4;

const SVG_NAMESPACE = "http://www.w3.org/2000/svg";

export interface Commit {
  Hash: string;
  Parents: string[];
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

interface Vertex {
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

interface Branch {
  Id: number;
  Color: number;
  Lines: Line[];
  UncommitedPoints: number;
  Merge: boolean;
}

interface Graph {
  Vertices: Vertex[];
  Branches: Branch[];
  Width: number;
  Height: number;
}

export function getLabelDist(x: number): number {
  return scaleX(x);
}

export function drawGraph(g: Graph): SVGSVGElement {
  let svg = document.createElementNS(SVG_NAMESPACE, "svg");
  let grp = document.createElementNS(SVG_NAMESPACE, "g");

  g.Branches.forEach((b) => {
    drawBranch(grp, b);
  });

  g.Vertices.forEach((v) => {
    if (v.BranchId != -1) {
      drawVertex(grp, v, g.Branches[v.BranchId]);
    }
  });

  svg.appendChild(grp);

  svg.setAttribute('height', scaleY(g.Height).toFixed(0).toString())
  svg.setAttribute('width', scaleY(g.Width).toFixed(0).toString())

  return svg;
}

function drawBranch(g: SVGGElement, b: Branch) {
  if (!b.Lines || !b.Lines.length) return;

  let color = (b.Color % NUM_COLORS).toString();

  // Remove middle points on consecutive straight lines.
  for (let i = 0; i < b.Lines.length - 1 /* iterate below */; ) {
    if (
      b.Lines[i].P1.X === b.Lines[i + 1].P1.X &&
      b.Lines[i].P1.Y === b.Lines[i + 1].P1.Y &&
      b.Lines[i].P2.X === b.Lines[i + 1].P2.X &&
      b.Lines[i].P2.Y === b.Lines[i + 1].P2.Y &&
      b.Lines[i].Committed === b.Lines[i + 1].Committed
    ) {
      b.Lines[i].P2.Y = b.Lines[i + 1].P2.Y;
      b.Lines.splice(i + 1, 1);
    } else {
      i++;
    }
  }

  let path = "";
  for (let i = 0; i < b.Lines.length; i++) {
    // If there's a current path and the new point is a different type of path.
    if (path && i && b.Lines[i].Committed !== b.Lines[i - 1].Committed) {
      let c = b.Lines[i - 1].Committed ? color : 'u';
      drawBranchPath(g, path, c);
      path = "";
    }

    let x1 = scaleX(b.Lines[i].P1.X).toFixed(0);
    let x2 = scaleX(b.Lines[i].P2.X).toFixed(0);
    let y1 = scaleY(b.Lines[i].P1.Y).toFixed(1);
    let y2 = scaleY(b.Lines[i].P2.Y).toFixed(1);

    // If no path or on different path
    if (
      !path ||
      (i &&
        (b.Lines[i].P1.X !== b.Lines[i - 1].P1.X ||
          b.Lines[i].P2.Y !== b.Lines[i - 1].P2.Y))
    ) {
      path += `M${x1},${y1}`;
    }

    // Vertical path
    if (x1 === x2) {
      path += `L${x2},${y2}`;
    }
    // Curved path
    else {
      let y1d = (scaleY(b.Lines[i].P1.Y) + CURVE).toFixed(1);
      let y2d = (scaleY(b.Lines[i].P2.Y) - CURVE).toFixed(1);
      path += `C${x1},${y1d} ${x2},${y2d} ${x2},${y2}`;
    }
  }

  if (path) {
    let c = b.Lines[b.Lines.length - 1].Committed ? color : 'u';
    drawBranchPath(g, path, c);
  }
}

function drawBranchPath(g: SVGGElement, path: string, color: string) {
  let l = document.createElementNS(SVG_NAMESPACE, "path");
  l.setAttribute("d", path);
  l.setAttribute("class", `b b-${color}`);
  g.appendChild(l);
}

function drawVertex(g: SVGGElement, v: Vertex, b: Branch) {
  let color = v.Committed ? (b.Color % NUM_COLORS).toString() : "u";

  let cx = scaleX(v.X).toString();
  let cy = scaleY(v.Id).toString();

  let c = document.createElementNS(SVG_NAMESPACE, "circle");
  c.setAttribute("cx", cx);
  c.setAttribute("cy", cy);
  if (v.Stash || !v.Committed) {
    c.setAttribute("r", VERTEX_RADIUS_U.toString());
    c.setAttribute("class", `v2 v-${color}`);
  }
  else {
    c.setAttribute("r", VERTEX_RADIUS.toString());
    c.setAttribute("class", `v v-${color}`);
  }
  g.appendChild(c);
}

function scaleX(x: number): number {
  return x * SCALE_X + OFFSET_X;
}

function scaleY(y: number): number {
  return y * SCALE_Y + OFFSET_Y;
}
