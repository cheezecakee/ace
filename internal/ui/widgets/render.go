package widgets

import "strings"

// Renderer draws a widget's content (no styling)
type Renderer interface {
	Draw(state RenderState) string
}

// RenderState contains everything a renderer needs
type RenderState struct {
	Cursor    Cursor
	Selection Selection
	Graph     *Graph // Contains items in NodeMeta
}

// =============================================================================
// GRID RENDERER - For 2D grid layouts (like menus)
// =============================================================================

type GridRenderer struct {
	items [][]Item // Original 2D structure for easy rendering
	cols  int
}

func (r *GridRenderer) Draw(state RenderState) string {
	var b strings.Builder

	for rowIdx, row := range r.items {
		for colIdx, item := range row {
			cursor := Cursor{Row: Row(rowIdx), Col: Col(colIdx)}

			// Highlight current cursor position
			if state.Cursor == cursor {
				b.WriteString("> ")
			} else {
				b.WriteString("  ")
			}

			b.WriteString(item.Label)
			b.WriteString("  ")
		}
		b.WriteString("\n")
	}

	return b.String()
}

// =============================================================================
// LIST RENDERER - For vertical lists
// =============================================================================

type ListRenderer struct {
	items []Item
}

func (r *ListRenderer) Draw(state RenderState) string {
	var b strings.Builder

	for i, item := range r.items {
		cursor := Cursor{Row: Row(i), Col: 0}

		if state.Cursor == cursor {
			b.WriteString("> ")
		} else {
			b.WriteString("  ")
		}

		b.WriteString(item.Label)
		b.WriteString("\n")
	}

	return b.String()
}

// =============================================================================
// CHECKBOX RENDERER - For multi-select lists
// =============================================================================

type CheckboxRenderer struct {
	items []Item
}

func (r *CheckboxRenderer) Draw(state RenderState) string {
	var b strings.Builder

	for i, item := range r.items {
		cursor := Cursor{Row: Row(i), Col: 0}

		// Show checkbox state
		if state.Selection.IsSelected(cursor) {
			b.WriteString("[x] ")
		} else {
			b.WriteString("[ ] ")
		}

		// Show cursor
		if state.Cursor == cursor {
			b.WriteString("> ")
		} else {
			b.WriteString("  ")
		}

		b.WriteString(item.Label)
		b.WriteString("\n")
	}

	return b.String()
}

// =============================================================================
// DUAL COLUMN RENDERER - For side-by-side lists
// =============================================================================

type DualColumnRenderer struct {
	leftItems  []Item
	rightItems []Item
}

func (r *DualColumnRenderer) Draw(state RenderState) string {
	var b strings.Builder

	maxRows := len(r.leftItems)
	if len(r.rightItems) > maxRows {
		maxRows = len(r.rightItems)
	}

	for row := 0; row < maxRows; row++ {
		// Left column
		if row < len(r.leftItems) {
			cursor := Cursor{Row: Row(row), Col: 0}

			if state.Cursor == cursor {
				b.WriteString("> ")
			} else {
				b.WriteString("  ")
			}

			b.WriteString(r.leftItems[row].Label)
		} else {
			b.WriteString("  ")
		}

		b.WriteString("    ") // Spacing between columns

		// Right column
		if row < len(r.rightItems) {
			cursor := Cursor{Row: Row(row), Col: 1}

			if state.Cursor == cursor {
				b.WriteString("> ")
			} else {
				b.WriteString("  ")
			}

			b.WriteString(r.rightItems[row].Label)
		}

		b.WriteString("\n")
	}

	return b.String()
}

// =============================================================================
// BAR RENDERER - For horizontal bars
// =============================================================================

type BarRenderer struct {
	items []Item
}

func (r *BarRenderer) Draw(state RenderState) string {
	var b strings.Builder

	for i, item := range r.items {
		cursor := Cursor{Row: 0, Col: Col(i)}

		if state.Cursor == cursor {
			b.WriteString("[")
			b.WriteString(item.Label)
			b.WriteString("]")
		} else {
			b.WriteString(" ")
			b.WriteString(item.Label)
			b.WriteString(" ")
		}

		b.WriteString("  ")
	}

	return b.String()
}

// =============================================================================
// CHOICE RENDERER - For single-choice questions with letter labels
// =============================================================================

type ChoiceRenderer struct{}

func (r *ChoiceRenderer) Draw(state RenderState) string {
	var b strings.Builder

	// Get items from the graph
	var items []Item
	i := 0
	for {
		cursor := Cursor{Row: Row(i), Col: 0}
		node := state.Graph.nodes[cursor]
		if node == nil {
			break
		}
		items = append(items, node.Meta.Item)
		i++
	}

	for i, item := range items {
		cursor := Cursor{Row: Row(i), Col: 0}

		// Show cursor
		if state.Cursor == cursor {
			b.WriteString("> ")
		} else {
			b.WriteString("  ")
		}

		// Show letter label [A], [B], [C], etc.
		label := string(rune('A' + i))
		b.WriteString("[" + label + "] ")

		b.WriteString(item.Label)
		b.WriteString("\n")
	}

	return b.String()
}

// =============================================================================
// MULTIPLE CHOICE RENDERER - For multi-select questions with letter labels
// =============================================================================

type MultipleChoiceRenderer struct{}

func (r *MultipleChoiceRenderer) Draw(state RenderState) string {
	var b strings.Builder

	// Get items from the graph
	var items []Item
	i := 0
	for {
		cursor := Cursor{Row: Row(i), Col: 0}
		node := state.Graph.nodes[cursor]
		if node == nil {
			break
		}
		items = append(items, node.Meta.Item)
		i++
	}

	for i, item := range items {
		cursor := Cursor{Row: Row(i), Col: 0}

		// Show checkbox state
		if state.Selection.IsSelected(cursor) {
			b.WriteString("[x] ")
		} else {
			b.WriteString("[ ] ")
		}

		// Show letter label [A], [B], [C], etc.
		label := string(rune('A' + i))
		b.WriteString("[" + label + "] ")

		b.WriteString(item.Label)
		b.WriteString("\n")
	}

	return b.String()
}

// =============================================================================
// BOOL RENDERER - For True/False questions
// =============================================================================

type BoolRenderer struct{}

func (r *BoolRenderer) Draw(state RenderState) string {
	var b strings.Builder

	// True option (Row 0)
	trueCursor := Cursor{Row: 0, Col: 0}
	if state.Cursor == trueCursor {
		b.WriteString("> ")
	} else {
		b.WriteString("  ")
	}
	b.WriteString("True\n")

	// False option (Row 1)
	falseCursor := Cursor{Row: 1, Col: 0}
	if state.Cursor == falseCursor {
		b.WriteString("> ")
	} else {
		b.WriteString("  ")
	}
	b.WriteString("False\n")

	return b.String()
}
