package git

const NULL_VERTEX = -1

type Line struct {
	P1        Point
	P2        Point
	Committed bool
	// true = P1, false = P2
	LockedDirection bool
	Merge           bool
	NullParent      bool
}

type Point struct {
	X uint16
	Y int64
}

type Connection struct {
	VertexId int64
	BranchId int64
}

type Graph struct {
	Vertices     []Vertex
	CommitLookup map[string]int64
	Branches     []GraphBranch
	Width        uint16
	Height       uint16
	Continues    bool
}

func (g *Git) BuildGraph(HEAD Ref, commits []Commit) Graph {
	// Build all graph data.
	graph := Graph{}
	graph.addVertices(commits, HEAD)
	graph.buildPaths()

	return graph
}

// Get vertices to plot based on commits and HEAD.
func (g *Graph) addVertices(commits []Commit, HEAD Ref) {
	g.CommitLookup = make(map[string]int64)
	for i, c := range commits {
		g.Vertices = append(g.Vertices, Vertex{
			Id:          int64(i),
			BranchId:    -1,
			Committed:   c.Hash != UNCOMMITED_HASH,
			Connections: make(map[uint16]Connection),
			Stash:       c.Stash,
		})
		g.CommitLookup[c.Hash] = int64(i)
	}

	// Assign each vertex its parents.
	for i, commit := range commits {
		for n, parent_hash := range commit.Parents {
			// Only use the first stash parent to build the graph. The second
			// is the stash's commit and should not be drawn.
			if commit.Stash && n > 0 {
				break
			}
			if parent_id, exists := g.CommitLookup[parent_hash]; exists {
				g.Vertices[i].Parents = append(g.Vertices[i].Parents, parent_id)
				g.Vertices[parent_id].Children = append(g.Vertices[parent_id].Children, int64(i))
			} else {
				g.Vertices[i].Parents = append(g.Vertices[i].Parents, NULL_VERTEX)
				g.Continues = true
			}
		}
	}
}

// Return whether the vertex is a merge commit.
func (g *Graph) isMergeCommit(v int64) bool {
	if g.Vertices[v].BranchId == -1 {
		return false
	}
	if len(g.Vertices[v].Parents) <= 1 {
		return false
	}
	if !g.Vertices[v].hasNextParent() {
		return false
	}
	p := g.Vertices[v].getNextParent()
	if p == NULL_VERTEX {
		return false
	}
	if g.Vertices[p].BranchId == -1 {
		return false
	}

	return true
}

// Build all paths for the graph.
func (g *Graph) buildPaths() {
	var color uint16 = 0
	for i := 0; i < len(g.Vertices); {
		// If the vertex has no parents or isn't on a branch yet, draw it.

		if g.Vertices[i].hasNextParent() || g.Vertices[i].BranchId == -1 {
			if g.isMergeCommit(int64(i)) {
				g.buildMergePath(&g.Vertices[i])
			} else {
				g.buildNormalPath(&g.Vertices[i], color)
				color++
			}
		} else {
			i++
		}
	}

	g.Width = g.getWidth()
	g.Height = g.getHeight()
}

// Build a path for the graph (that isn't a merge commit).
func (g *Graph) buildNormalPath(v *Vertex, color uint16) {
	// Create new branch.
	b := GraphBranch{
		Color: color,
	}
	b.Id = int64(len(g.Branches))
	g.Branches = append(g.Branches, b)

	// Current vertex and current parent.
	v1 := &g.Vertices[v.Id]
	var p *Vertex = nil
	null_parent := false
	if v.hasNextParent() && v.getNextParent() != NULL_VERTEX {
		p = &g.Vertices[v.getNextParent()]
	} else if v.getNextParent() == NULL_VERTEX {
		null_parent = true
	}

	// Previous point.
	var p1 Point
	if v.BranchId == -1 {
		p1 = v.getNextPoint()
	} else {
		p1 = v.getPoint()
	}
	v.addUnavailPoint(p1.X, v, b.Id)

	// If vertex doesn't have a branch yet, assign it and set x to the previous point.
	if v.BranchId == -1 {
		v.BranchId = b.Id
		v.X = p1.X
	}

	// Loop vertices after the current.
	var i int
	for i = int(v.Id) + 1; i < len(g.Vertices); i++ {
		var p2 Point
		if p != nil && p.Id == int64(i) && p.BranchId != -1 {
			p2 = g.Vertices[i].getPoint()
		} else {
			p2 = g.Vertices[i].getNextPoint()
		}

		// Add the line, mark the point unavail, move to next point.
		g.Branches[b.Id].addLine(Line{
			P1:              p1,
			P2:              p2,
			Committed:       v1.Committed,
			LockedDirection: p1.X < p2.X,
			Merge:           false,
		})
		g.Vertices[i].addUnavailPoint(p2.X, p, b.Id)

		p1 = p2

		// If the parent of v1 has been reached,
		if p != nil && p.Id == int64(i) {
			// Is the parent already on a branch.
			pb := p.BranchId != -1

			// Assign branch to parent.
			if !pb {
				g.Vertices[p.Id].BranchId = b.Id
				g.Vertices[p.Id].X = p2.X
			}

			// Update the current vertex's parent.
			v1.NextParent++
			// Move the current parent to the current vertex.
			v1 = &g.Vertices[p.Id]
			// Move the next parent to the current parent.
			if v1.hasNextParent() && v1.getNextParent() != NULL_VERTEX && int(v1.getNextParent()) < len(g.Vertices) {
				p = &g.Vertices[v1.getNextParent()]
				null_parent = false
			} else {
				p = nil
				if v1.hasNextParent() && v1.getNextParent() == NULL_VERTEX {
					null_parent = true
				}
				// If p2 is the final commit in the displayed commits and it has a parent.
				if int(p2.Y) == len(g.Vertices)-1 && len(g.Vertices[len(g.Vertices)-1].Parents) > 0 {
					null_point := Point{
						X: p2.X,
						Y: -1,
					}

					g.Branches[b.Id].addLine(Line{
						P1:              p2,
						P2:              null_point,
						Committed:       true,
						LockedDirection: false,
						Merge:           false,
					})
				}
			}

			// If there's no next parent or if the new current vertex already has a branch.
			if p == nil || pb {
				break
			}
		}
	}

	if i == len(g.Vertices) && len(g.Branches[b.Id].Lines) > 0 {
		p1 := g.Branches[b.Id].Lines[len(g.Branches[b.Id].Lines)-1].P2
		null_point := Point{
			X: p1.X,
			Y: -1,
		}

		g.Branches[b.Id].addLine(Line{
			P1:              p1,
			P2:              null_point,
			Committed:       true,
			LockedDirection: false,
			Merge:           false,
		})
	}

	// If we looped through every vertex and the parent is a null vertex.
	if i == len(g.Vertices) && null_parent {
		v.NextParent++
	}
}

// Build a merge path for the graph.
func (g *Graph) buildMergePath(v1 *Vertex) {
	// Parent
	var p *Vertex = nil
	if v1.hasNextParent() {
		p = &g.Vertices[v1.getNextParent()]
	}

	// Previous point.
	var p1 Point
	if v1.BranchId == -1 {
		p1 = v1.getNextPoint()
	} else {
		p1 = v1.getPoint()
	}

	var i int
	for i = int(v1.Id) + 1; i < len(g.Vertices); i++ {
		// Get point connecting to parent vertex.
		var p2 Point
		found := false
		for n, c := range g.Vertices[i].Connections {
			if c.VertexId == p.Id && c.BranchId == p.BranchId {
				p2 = Point{
					X: uint16(n),
					Y: int64(i),
				}
				found = true
			}
		}

		// If point wasn't found connected to parent, move to next point.
		if !found {
			p2 = g.Vertices[i].getNextPoint()
		}

		// Which point on the line is locked.
		dir := true
		if !found && int64(i) != p.Id {
			dir = p1.X < p2.X
		}

		g.Branches[p.BranchId].addLine(Line{
			P1:              p1,
			P2:              p2,
			Committed:       v1.Committed,
			LockedDirection: dir,
			Merge:           true,
		})

		g.Vertices[i].addUnavailPoint(p2.X, p, p.BranchId)

		p1 = p2

		// If point was found connected to parent, move vertex to next parent and be done.
		if found {
			v1.NextParent++
			break
		}
	}
}

// Get the width of the graph in vertices.
func (g *Graph) getWidth() uint16 {
	var x uint16 = 0
	for _, v := range g.Vertices {
		p := v.getNextPoint()
		if p.X > x {
			x = p.X
		}
	}
	return x
}

// Get the height of the graph in commits.
func (g *Graph) getHeight() uint16 {
	return uint16(len(g.Vertices))
}

type Vertex struct {
	Id          int64
	Children    []int64
	Parents     []int64
	NextParent  int
	BranchId    int64
	X           uint16
	XNext       uint16
	Connections map[uint16]Connection
	Committed   bool
	Stash       bool
}

func (v *Vertex) hasNextParent() bool {
	return v.NextParent < len(v.Parents)
}

func (v *Vertex) getNextParent() int64 {
	return v.Parents[v.NextParent]
}

func (v *Vertex) getPoint() Point {
	return Point{
		X: v.X,
		Y: v.Id,
	}
}

func (v *Vertex) getNextPoint() Point {
	return Point{
		X: v.XNext,
		Y: v.Id,
	}
}

func (v *Vertex) addUnavailPoint(x uint16, v2 *Vertex, b int64) {
	if x == v.XNext {
		v.XNext = x + 1
		var vId int64 = NULL_VERTEX
		if v2 != nil {
			vId = v2.Id
		}
		v.Connections[x] = Connection{
			VertexId: vId,
			BranchId: b,
		}
	}
}

type GraphBranch struct {
	Id               int64
	Color            uint16
	Lines            []Line
	UncommitedPoints uint16
}

// Add line to branch.
func (b *GraphBranch) addLine(l Line) {
	b.Lines = append(b.Lines, l)
	if l.Committed {
		if l.P2.X == 0 && l.P2.X < b.UncommitedPoints && l.P2.Y != NULL_VERTEX {
			b.UncommitedPoints = uint16(l.P2.Y)
		}
	} else {
		b.UncommitedPoints++
	}
}
