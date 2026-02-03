package widgets

//================= GRAPH BUILDER ===================//

type GraphBuilder struct {
	graph *Graph
}

func (b *GraphBuilder) Node(c Cursor, meta NodeMeta) *GraphBuilder {
	b.graph.AddNode(c, meta)
	return b
}

func (b *GraphBuilder) Edge(from Cursor, dir Direction, to Cursor) *GraphBuilder {
	b.graph.Connect(from, dir, to)
	return b
}

func (b *GraphBuilder) BiEdge(a Cursor, dir Direction, bCursor Cursor) *GraphBuilder {
	b.graph.Connect(a, dir, bCursor)
	b.graph.Connect(bCursor, invert(dir), a)
	return b
}

func (b *GraphBuilder) Build() *Graph {
	return b.graph
}

func invert(d Direction) Direction {
	switch d {
	case Left:
		return Right
	case Right:
		return Left
	case Top:
		return Down
	case Down:
		return Top
	}

	return d
}

//==================== GRAPH =======================//

type Graph struct {
	nodes map[Cursor]*Node
}

type Edges map[Direction]Cursor

type Node struct {
	Cursor Cursor
	Edges  Edges
	Meta   NodeMeta
}

type NodeMeta struct {
	Item    Item
	Enabled bool
	Empty   bool
	Tags    []string
}

func NewGraph() *Graph {
	return &Graph{
		nodes: make(map[Cursor]*Node),
	}
}

func (g *Graph) AddNode(c Cursor, meta NodeMeta) {
	g.nodes[c] = &Node{
		Cursor: c,
		Edges:  make(map[Direction]Cursor),
		Meta:   meta,
	}
}

func (g *Graph) Connect(from Cursor, dir Direction, to Cursor) {
	if from == to {
		return
	}

	node := g.nodes[from]
	if node == nil {
		return
	}

	node.Edges[dir] = to
}

func (g *Graph) Move(from Cursor, dir Direction) (Cursor, bool) {
	node := g.nodes[from]
	if node == nil || !node.Meta.Enabled {
		return from, false
	}

	next, ok := node.Edges[dir]
	if !ok {
		return from, false
	}

	target := g.nodes[next]
	if target == nil || !target.Meta.Enabled || target.Meta.Empty {
		return from, false
	}

	return next, true
}
