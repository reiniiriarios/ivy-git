package main

const NULL_VERTEX = -1

type Line struct {
	P1        Point
	P2        Point
	Committed bool
	// true = P1, false = P2
	LockedDirection bool
	Merge           bool
}

type Point struct {
	X uint16
	Y int64
}

type Connection struct {
	VertexId int64
	BranchId int64
}

// Get vertices to plot based on commits and HEAD.
func (a *App) getVertices(commits []Commit, HEAD Ref) ([]Vertex, map[string]int64) {
	vertices := []Vertex{}

	lookup := make(map[string]int64)
	for i, c := range commits {
		vertices = append(vertices, Vertex{
			Id:          int64(i),
			Committed:   c.Hash != UNCOMMITED_HASH,
			Connections: make(map[uint16]Connection),
			Stash:       c.Stash,
		})
		lookup[c.Hash] = int64(i)
	}

	// Assign each vertex its parents.
	for i, commit := range commits {
		for _, parent_hash := range commit.Parents {
			if parent_id, exists := lookup[parent_hash]; exists {
				vertices[i].Parents = append(vertices[i].Parents, parent_id)
				vertices[parent_id].Children = append(vertices[parent_id].Children, int64(i))
			} else {
				vertices[i].Parents = append(vertices[i].Parents, NULL_VERTEX)
			}
		}
	}

	return vertices, lookup
}

type Vertex struct {
	Id          int64
	Children    []int64
	Parents     []int64
	NextParent  int
	Branch      *GraphBranch
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

func (v *Vertex) addUnavailPoint(x uint16, v2 *Vertex, b *GraphBranch) {
	if x == v.XNext {
		v.XNext = x + 1
		var vId int64 = -1
		if v2 != nil {
			vId = v2.Id
		}
		v.Connections[x] = Connection{
			VertexId: vId,
			BranchId: b.Id,
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
		if l.P2.X == 0 && l.P2.X < b.UncommitedPoints {
			b.UncommitedPoints = uint16(l.P2.Y)
		}
	} else {
		b.UncommitedPoints++
	}
}

type Graph struct {
	Vertices []Vertex
	Branches []GraphBranch
	Width    uint16
	Height   uint16
}

// Return whether the vertex is a merge commit.
func (g *Graph) isMergeCommit(v int64) bool {
	if g.Vertices[v].Branch == nil {
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
	if g.Vertices[p].Branch == nil {
		return false
	}

	return true
}

// Build all paths for the graph.
func (g *Graph) BuildPaths() {
	var color uint16 = 0
	for i := 0; i < len(g.Vertices); {
		// If the vertex has no parents or isn't on a branch yet, draw it.

		if g.Vertices[i].hasNextParent() || g.Vertices[i].Branch == nil {
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

	// Current vertex and current parent.
	v1 := &g.Vertices[v.Id]
	var p *Vertex = nil
	if v.hasNextParent() && v.getNextParent() != NULL_VERTEX {
		p = &g.Vertices[v.getNextParent()]
	}

	// Previous point.
	var p1 Point
	if v.Branch == nil {
		p1 = v.getNextPoint()
	} else {
		p1 = v.getPoint()
	}
	v.addUnavailPoint(p1.X, v, &b)

	// If vertex doesn't have a branch yet, assign it and set x to the previous point.
	if v.Branch == nil {
		v.Branch = &b
		v.X = p1.X
	}

	// Loop vertices after the current.
	var i int
	for i = int(v.Id) + 1; i < len(g.Vertices); i++ {
		var p2 Point
		if p != nil && p.Id == int64(i) && p.Branch != nil {
			p2 = g.Vertices[i].getPoint()
		} else {
			p2 = g.Vertices[i].getNextPoint()
		}

		// Add the line, mark the point unavail, move to next point.
		b.addLine(Line{
			P1:              p1,
			P2:              p2,
			Committed:       v1.Committed,
			LockedDirection: p1.X < p2.X,
			Merge:           false,
		})
		g.Vertices[i].addUnavailPoint(p2.X, p, &b)

		p1 = p2

		// If the parent of v1 has been reached,
		if p != nil && p.Id == int64(i) {
			// Is the parent already on a branch.
			pb := p.Branch != nil

			// Assign branch to parent.
			g.Vertices[p.Id].Branch = &b
			g.Vertices[p.Id].X = p2.X

			// Update the current vertex's parent.
			v1.NextParent++
			// Move the current parent to the current vertex.
			v1 = &g.Vertices[p.Id]
			// Move the next parent to the current parent.
			if v1.hasNextParent() {
				p = &g.Vertices[v1.getNextParent()]
			} else {
				p = nil
			}

			// If there's no next parent or if the new current vertex already has a branch.
			if p == nil || pb {
				break
			}
		}
	}

	// If we looped through every vertex and the parent is a null vertex.
	if i == len(g.Vertices) && p != nil && p.Id == NULL_VERTEX {
		v.NextParent++
	}

	b.Id = int64(len(g.Branches))
	g.Branches = append(g.Branches, b)
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
	if v1.Branch == nil {
		p1 = v1.getNextPoint()
	} else {
		p1 = v1.getPoint()
	}

	for _, v2 := range g.Vertices {
		// Get point connecting to parent vertex.
		var p2 Point
		found := false
		for i, c := range v2.Connections {
			if c.VertexId == p.Id && c.BranchId == p.Branch.Id {
				p2 = Point{
					X: uint16(i),
					Y: v2.Id,
				}
				found = true
			}
		}

		// If point wasn't found connected to parent, move to next point.
		if !found {
			p2 = v2.getNextPoint()
		}

		// Which point on the line is locked.
		dir := true
		if !found && v2.Id != p.Id {
			dir = p1.X < p2.X
		}

		p.Branch.addLine(Line{
			P1:              p1,
			P2:              p2,
			Committed:       v1.Committed,
			LockedDirection: dir,
			Merge:           true,
		})

		g.Vertices[v2.Id].addUnavailPoint(p2.X, p, p.Branch)

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
