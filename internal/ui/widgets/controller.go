package widgets

type Controller struct {
	Cursor    Cursor
	Selection Selection
	Graph     *Graph
}

func (c *Controller) Move(dir Direction) {
	if c.Graph == nil {
		return
	}

	if next, ok := c.Graph.Move(c.Cursor, dir); ok {
		c.Cursor = next
	}
}

func (c *Controller) Toggle() {
	if c.Selection.IsSelectable(c.Cursor) {
		c.Selection.Toggle(c.Cursor)
	}
}

func (c *Controller) Select() {
	if c.Selection.IsSelectable(c.Cursor) {
		c.Selection.Select(c.Cursor)
	}
}
