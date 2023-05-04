package main

type Line struct {
	P1        Point
	P2        Point
	Committed bool
	// true = P1, false = P2
	LockedDirection bool
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
func (a *App) getVertices(commits []Commit, HEAD Ref) []Vertex {
	var vertices []Vertex
	nullVertex := Vertex{
		Id: -1,
	}

	lookup := make(map[string]int)
	for i, c := range commits {
		vertices[i] = Vertex{
			Id:        int64(i),
			Committed: c.Hash != UNCOMMITED_HASH,
		}
		lookup[c.Hash] = i
	}

	// Assign each vertex its parents.
	for i, commit := range commits {
		for _, parent_hash := range commit.Parents {
			if parent_id, exists := lookup[parent_hash]; exists {
				vertices[i].Parents = append(vertices[i].Parents, vertices[parent_id])
			} else {
				vertices[i].Parents = append(vertices[i].Parents, nullVertex)
			}
		}
	}

	return vertices
}

type Vertex struct {
	Id          int64
	Children    []Vertex
	Parents     []Vertex
	NextParent  int
	Branch      *GraphBranch
	X           uint16
	XNext       uint16
	Connections []Connection
	Committed   bool
}

func (v *Vertex) getNextParent() *Vertex {
	if v.NextParent < len(v.Parents) {
		return &v.Parents[v.NextParent]
	}
	return nil
}

func (v *Vertex) getPoint() *Point {
	return &Point{
		X: v.X,
		Y: v.Id,
	}
}

func (v *Vertex) getNextPoint() *Point {
	return &Point{
		X: v.XNext,
		Y: v.Id,
	}
}

func (v *Vertex) addUnavailPoint(x uint16, v2 *Vertex, b *GraphBranch) {
	if x == v.XNext {
		v.XNext = x + 1
		v.Connections = append(v.Connections, Connection{
			VertexId: v2.Id,
			BranchId: b.Id,
		})
	}
}

// Return whether the vertex is a merge commit.
func (v *Vertex) isMergeCommit() bool {
	// If there's a next parent vertex that isn't a null vertex,
	// the vertex has more than one parent,
	// and the vertex and next parent both have a branch.
	p := v.getNextParent()
	return p != nil && p.Id != -1 && len(v.Parents) > 1 && v.Branch != nil && p.Branch != nil
}

type GraphBranch struct {
	Id               int64
	Color            uint16
	Lines            []Line
	FinalVertexId    int64
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

// Build all paths for the graph.
func (g *Graph) BuildPaths() {
	var color uint16 = 0
	for _, v := range g.Vertices {
		// If there are no parents and there's already a branch, skip.
		if len(v.Parents) == 0 && v.Branch != nil {
			return
		}

		if v.isMergeCommit() {
			g.buildMergePath(&v)
		} else {
			g.buildNormalPath(&v, color)
			color++
		}
	}

	g.Width = g.getWidth()
	g.Height = g.getHeight()
}

// Build a path for the graph (that isn't a merge commit).
func (g *Graph) buildNormalPath(v *Vertex, color uint16) {
	// Create new branch, assign it to the vertex.
	b := GraphBranch{
		Color: color,
	}
	if v.Branch != nil {
		v.Branch = &b
	}

	// Current vertex and current parent.
	v1 := v
	p := v.getNextParent()

	// Previous point.
	var p1 *Point
	if v.Branch == nil {
		p1 = v.getNextPoint()
	} else {
		p1 = v.getPoint()
	}
	v.addUnavailPoint(p1.X, v, &b)

	// Loop vertices after the current.
	var i int64 = v.Id + 1
	for _, v2 := range g.Vertices[i:] {
		var p2 *Point
		if p.Id == v2.Id && p.Branch != nil {
			p2 = v2.getPoint()
		} else {
			p2 = v2.getNextPoint()
		}

		// Add the line, mark the point unavail, move to next point.
		b.addLine(Line{
			P1:              *p1,
			P2:              *p2,
			Committed:       v.Committed,
			LockedDirection: p1.X < p2.X,
		})
		v2.addUnavailPoint(p2.X, p, &b)
		p1 = p2

		// If the parent of v1 has been reached,
		if p.Id == v2.Id {
			// Assign branch to parent.
			p.Branch = &b
			p.X = p2.X

			// Update the current vertex's parent.
			v1.NextParent++
			// Move the current parent to the current vertex.
			v1 = p
			// Move the next parent to the current parent.
			p = v1.getNextParent()

			// If there's no next parent or if the new current vertex already has a branch.
			if p == nil || v1.Branch != nil {
				break
			}
		}

		i++
	}
	// If we looped through every vertex and
	// either there's no parent or the parent is a null vertex,
	if int(i) == len(g.Vertices) && (p != nil || p.Id == -1) {
		v.NextParent++
	}

	b.FinalVertexId = i // VERIFY!!!!!
	b.Id = int64(len(g.Branches))
	g.Branches = append(g.Branches, b)
}

// Build a merge path for the graph.
func (g *Graph) buildMergePath(v1 *Vertex) {
	// Parent
	p := v1.getNextParent()

	// Previous point.
	var p1 *Point
	if v1.Branch == nil {
		p1 = v1.getNextPoint()
	} else {
		p1 = v1.getPoint()
	}

	for _, v2 := range g.Vertices {
		// Get point connecting to parent vertex.
		var p2 *Point
		for i, c := range v2.Connections {
			if c.VertexId == p.Id && c.BranchId == p.Branch.Id {
				p2 = &Point{
					X: uint16(i),
					Y: v2.Id,
				}
			}
		}

		// If point was found connected to parent, move to next point.
		found := p2 != nil
		if found {
			p2 = v2.getNextPoint()
		}

		// Which point on the line is locked.
		dir := true
		if found || v2.Id == p.Id {
			dir = p1.X < p2.X
		}

		p.Branch.addLine(Line{
			P1:              *p1,
			P2:              *p2,
			Committed:       v1.Committed,
			LockedDirection: dir,
		})

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
