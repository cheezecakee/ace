// Package widgets
package widgets

// Widget is a self-contained navigable UI component
// Combines navigation (graph, cursor, selection) with rendering
type Widget struct {
	Controller          // Embedded: provides Move(), Select(), Toggle(), and Graph/Cursor/Selection
	renderer   Renderer // How to draw this widget
}

// Render produces the widget's visual content (plain text, no styling)
func (w *Widget) Render() string {
	return w.renderer.Draw(RenderState{
		Cursor:    w.Cursor,
		Selection: w.Selection,
		Graph:     w.Graph,
	})
}

// GetItem returns the item at the current cursor position
func (w *Widget) GetItem() (Item, bool) {
	node := w.Graph.nodes[w.Cursor]
	if node == nil {
		return Item{}, false
	}
	return node.Meta.Item, true
}

// GetSelectedItems returns all selected items (for multi-select widgets)
func (w *Widget) GetSelectedItems() []Item {
	var selected []Item

	// Only works with Multi selection type
	multiSel, ok := w.Selection.(*Multi)
	if !ok {
		return selected
	}

	for cursor, isSelected := range multiSel.Selected {
		if !isSelected {
			continue
		}

		node := w.Graph.nodes[cursor]
		if node != nil {
			selected = append(selected, node.Meta.Item)
		}
	}

	return selected
}

// =============================================================================
// WIDGET FACTORY FUNCTIONS - Pre-configured widget types
// =============================================================================

// NewGrid creates a 2D grid widget (like a menu)
// items is a 2D slice where items[row][col] is the item at that position
func NewGrid(items [][]Item, cols int) *Widget {
	rows := len(items)

	// Build graph using GraphBuilder
	builder := &GraphBuilder{graph: NewGraph()}

	for r := 0; r < rows; r++ {
		for c := 0; c < len(items[r]); c++ {
			cursor := Cursor{Row: Row(r), Col: Col(c)}

			// Add node with item in metadata
			builder.Node(cursor, NodeMeta{
				Item:    items[r][c],
				Enabled: true,
				Empty:   false,
			})
		}
	}

	// Connect nodes bi-directionally
	for r := 0; r < rows; r++ {
		for c := 0; c < len(items[r]); c++ {
			cursor := Cursor{Row: Row(r), Col: Col(c)}

			// Connect right
			if c < len(items[r])-1 {
				builder.BiEdge(cursor, Right, Cursor{Row: Row(r), Col: Col(c + 1)})
			}

			// Connect down
			if r < rows-1 && c < len(items[r+1]) {
				builder.BiEdge(cursor, Down, Cursor{Row: Row(r + 1), Col: Col(c)})
			}
		}
	}

	graph := builder.Build()

	return &Widget{
		Controller: Controller{
			Cursor:    Cursor{Row: 0, Col: 0},
			Selection: &Single{},
			Graph:     graph,
		},
		renderer: &GridRenderer{items: items, cols: cols},
	}
}

// NewList creates a vertical list widget
func NewList(items []Item) *Widget {
	builder := &GraphBuilder{graph: NewGraph()}

	// Add all items as nodes
	for i, item := range items {
		cursor := Cursor{Row: Row(i), Col: 0}
		builder.Node(cursor, NodeMeta{
			Item:    item,
			Enabled: true,
			Empty:   false,
		})
	}

	// Connect nodes vertically
	for i := 0; i < len(items)-1; i++ {
		from := Cursor{Row: Row(i), Col: 0}
		to := Cursor{Row: Row(i + 1), Col: 0}
		builder.BiEdge(from, Down, to)
	}

	graph := builder.Build()

	return &Widget{
		Controller: Controller{
			Cursor:    Cursor{Row: 0, Col: 0},
			Selection: &Single{},
			Graph:     graph,
		},
		renderer: &ListRenderer{items: items},
	}
}

// NewCheckboxList creates a multi-select list widget
func NewCheckboxList(items []Item) *Widget {
	builder := &GraphBuilder{graph: NewGraph()}

	for i, item := range items {
		cursor := Cursor{Row: Row(i), Col: 0}
		builder.Node(cursor, NodeMeta{
			Item:    item,
			Enabled: true,
			Empty:   false,
		})
	}

	for i := 0; i < len(items)-1; i++ {
		from := Cursor{Row: Row(i), Col: 0}
		to := Cursor{Row: Row(i + 1), Col: 0}
		builder.BiEdge(from, Down, to)
	}

	graph := builder.Build()

	return &Widget{
		Controller: Controller{
			Cursor:    Cursor{Row: 0, Col: 0},
			Selection: &Multi{Selected: make(map[Cursor]bool)},
			Graph:     graph,
		},
		renderer: &CheckboxRenderer{items: items},
	}
}

// NewBar creates a horizontal bar widget
func NewBar(items []Item) *Widget {
	builder := &GraphBuilder{graph: NewGraph()}

	// Add all items as nodes in a horizontal row
	for i, item := range items {
		cursor := Cursor{Row: 0, Col: Col(i)}
		builder.Node(cursor, NodeMeta{
			Item:    item,
			Enabled: true,
			Empty:   false,
		})
	}

	// Connect nodes horizontally
	for i := 0; i < len(items)-1; i++ {
		from := Cursor{Row: 0, Col: Col(i)}
		to := Cursor{Row: 0, Col: Col(i + 1)}
		builder.BiEdge(from, Right, to)
	}

	graph := builder.Build()

	return &Widget{
		Controller: Controller{
			Cursor:    Cursor{Row: 0, Col: 0},
			Selection: &Single{},
			Graph:     graph,
		},
		renderer: &BarRenderer{items: items},
	}
}

// =============================================================================
// QUESTION WIDGETS - For quiz/game question rendering
// =============================================================================

// NewChoiceWidget creates a single-choice question widget with letter labels [A], [B], [C]
func NewChoiceWidget(items []Item) *Widget {
	builder := &GraphBuilder{graph: NewGraph()}

	// Add all items as nodes
	for i, item := range items {
		cursor := Cursor{Row: Row(i), Col: 0}
		builder.Node(cursor, NodeMeta{
			Item:    item,
			Enabled: true,
			Empty:   false,
		})
	}

	// Connect nodes vertically
	for i := 0; i < len(items)-1; i++ {
		from := Cursor{Row: Row(i), Col: 0}
		to := Cursor{Row: Row(i + 1), Col: 0}
		builder.BiEdge(from, Down, to)
	}

	graph := builder.Build()

	return &Widget{
		Controller: Controller{
			Cursor:    Cursor{Row: 0, Col: 0},
			Selection: &Single{},
			Graph:     graph,
		},
		renderer: &ChoiceRenderer{},
	}
}

// NewMultipleChoiceWidget creates a multi-select question widget with checkboxes and letter labels
func NewMultipleChoiceWidget(items []Item) *Widget {
	builder := &GraphBuilder{graph: NewGraph()}

	for i, item := range items {
		cursor := Cursor{Row: Row(i), Col: 0}
		builder.Node(cursor, NodeMeta{
			Item:    item,
			Enabled: true,
			Empty:   false,
		})
	}

	for i := 0; i < len(items)-1; i++ {
		from := Cursor{Row: Row(i), Col: 0}
		to := Cursor{Row: Row(i + 1), Col: 0}
		builder.BiEdge(from, Down, to)
	}

	graph := builder.Build()

	return &Widget{
		Controller: Controller{
			Cursor:    Cursor{Row: 0, Col: 0},
			Selection: &Multi{Selected: make(map[Cursor]bool)},
			Graph:     graph,
		},
		renderer: &MultipleChoiceRenderer{},
	}
}

// NewBoolWidget creates a True/False question widget (no items needed, always True/False)
func NewBoolWidget() *Widget {
	builder := &GraphBuilder{graph: NewGraph()}

	// Add True and False as nodes
	builder.Node(Cursor{Row: 0, Col: 0}, NodeMeta{
		Item:    Item{Label: "True"},
		Enabled: true,
		Empty:   false,
	})
	builder.Node(Cursor{Row: 1, Col: 0}, NodeMeta{
		Item:    Item{Label: "False"},
		Enabled: true,
		Empty:   false,
	})

	// Connect True -> False
	builder.BiEdge(Cursor{Row: 0, Col: 0}, Down, Cursor{Row: 1, Col: 0})

	graph := builder.Build()

	return &Widget{
		Controller: Controller{
			Cursor:    Cursor{Row: 0, Col: 0},
			Selection: &Single{},
			Graph:     graph,
		},
		renderer: &BoolRenderer{},
	}
}
