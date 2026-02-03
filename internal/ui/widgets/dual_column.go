package widgets

// DualColumnWidget is a specialized widget for moving items between two columns
// Commonly used for "available vs selected" patterns
type DualColumnWidget struct {
	*Widget    // Embed base widget
	leftItems  []Item
	rightItems []Item
}

// NewDualColumn creates a two-column widget for moving items between lists
func NewDualColumn(leftItems, rightItems []Item) *DualColumnWidget {
	widget := newDualColumnBase(leftItems, rightItems)

	return &DualColumnWidget{
		Widget:     widget,
		leftItems:  leftItems,
		rightItems: rightItems,
	}
}

// Internal helper to build the base widget
func newDualColumnBase(leftItems, rightItems []Item) *Widget {
	builder := &GraphBuilder{graph: NewGraph()}

	maxRows := len(leftItems)
	if len(rightItems) > maxRows {
		maxRows = len(rightItems)
	}

	// Add left column nodes
	for i := 0; i < len(leftItems); i++ {
		cursor := Cursor{Row: Row(i), Col: 0}
		builder.Node(cursor, NodeMeta{
			Item:    leftItems[i],
			Enabled: true,
			Empty:   false,
		})
	}

	// Add right column nodes
	for i := 0; i < len(rightItems); i++ {
		cursor := Cursor{Row: Row(i), Col: 1}
		builder.Node(cursor, NodeMeta{
			Item:    rightItems[i],
			Enabled: true,
			Empty:   false,
		})
	}

	// Connect vertically within columns
	for i := 0; i < len(leftItems)-1; i++ {
		from := Cursor{Row: Row(i), Col: 0}
		to := Cursor{Row: Row(i + 1), Col: 0}
		builder.BiEdge(from, Down, to)
	}

	for i := 0; i < len(rightItems)-1; i++ {
		from := Cursor{Row: Row(i), Col: 1}
		to := Cursor{Row: Row(i + 1), Col: 1}
		builder.BiEdge(from, Down, to)
	}

	// Connect horizontally between columns
	for i := 0; i < maxRows; i++ {
		if i < len(leftItems) && i < len(rightItems) {
			left := Cursor{Row: Row(i), Col: 0}
			right := Cursor{Row: Row(i), Col: 1}
			builder.BiEdge(left, Right, right)
		}
	}

	graph := builder.Build()

	return &Widget{
		Controller: Controller{
			Cursor: Cursor{Row: 0, Col: 0},
			Selection: &ColumnBased{
				Left:      make(map[Row]bool),
				Right:     make(map[Row]bool),
				ActiveCol: 0,
			},
			Graph: graph,
		},
		renderer: &DualColumnRenderer{leftItems: leftItems, rightItems: rightItems},
	}
}

// MoveItem moves the item at the current cursor position to the other column
// Returns true if an item was moved
func (w *DualColumnWidget) MoveItem() bool {
	row := int(w.Cursor.Row)
	col := int(w.Cursor.Col)

	// Moving from left to right
	if col == 0 && row < len(w.leftItems) {
		item := w.leftItems[row]

		// Skip empty items
		if item.Label == "" {
			return false
		}

		// Move item
		w.rightItems = append(w.rightItems, item)
		w.leftItems[row] = Item{Label: ""} // Mark as empty

		// Move cursor to right column, same row
		w.Cursor = Cursor{Row: Row(len(w.rightItems) - 1), Col: 1}

		// Rebuild widget with new items
		w.rebuild()
		return true
	}

	// Moving from right to left
	if col == 1 && row < len(w.rightItems) {
		item := w.rightItems[row]

		// Skip empty items
		if item.Label == "" {
			return false
		}

		// Move item
		w.leftItems = append(w.leftItems, item)
		w.rightItems[row] = Item{Label: ""} // Mark as empty

		// Move cursor to left column, same row
		w.Cursor = Cursor{Row: Row(len(w.leftItems) - 1), Col: 0}

		// Rebuild widget with new items
		w.rebuild()
		return true
	}

	return false
}

// rebuild reconstructs the widget's graph and renderer with current items
func (w *DualColumnWidget) rebuild() {
	newWidget := newDualColumnBase(w.leftItems, w.rightItems)

	// Preserve cursor position (already updated by MoveItem)
	newWidget.Cursor = w.Cursor

	// Update the embedded widget
	w.Graph = newWidget.Graph
	w.renderer = newWidget.renderer
}

// GetLeftItems returns all non-empty items in the left column
func (w *DualColumnWidget) GetLeftItems() []Item {
	var items []Item
	for _, item := range w.leftItems {
		if item.Label != "" {
			items = append(items, item)
		}
	}
	return items
}

// GetRightItems returns all non-empty items in the right column
func (w *DualColumnWidget) GetRightItems() []Item {
	var items []Item
	for _, item := range w.rightItems {
		if item.Label != "" {
			items = append(items, item)
		}
	}
	return items
}

// GetRightItemLabels returns just the labels of items in the right column
// Useful for getting IDs or names of selected items
func (w *DualColumnWidget) GetRightItemLabels() []string {
	var labels []string
	for _, item := range w.rightItems {
		if item.Label != "" {
			labels = append(labels, item.Label)
		}
	}
	return labels
}

// GetLeftItemLabels returns just the labels of items in the left column
func (w *DualColumnWidget) GetLeftItemLabels() []string {
	var labels []string
	for _, item := range w.leftItems {
		if item.Label != "" {
			labels = append(labels, item.Label)
		}
	}
	return labels
}
