package widgets

type Selection interface {
	Toggle(cursor Cursor)
	Select(cursor Cursor)
	IsSelectable(cursor Cursor) bool
	IsSelected(cursor Cursor) bool
}

type Single struct{} // menus

func (s *Single) Toggle(cursor Cursor) {
	// no-go for menus
}

func (s *Single) Select(cursor Cursor) {
	// navigation layer will react to cursor externally
}

func (s *Single) IsSelectable(cursor Cursor) bool {
	return true
}

func (s *Single) IsSelected(cursor Cursor) bool {
	return false
}

type Multi struct {
	Selected map[Cursor]bool
} // checkbox list

func (s *Multi) Toggle(cursor Cursor) {
	if s.Selected == nil {
		s.Selected = make(map[Cursor]bool)
	}
	s.Selected[cursor] = !s.Selected[cursor]
}

func (s *Multi) Select(cursor Cursor) {
	s.Toggle(cursor)
}

func (s *Multi) IsSelectable(cursor Cursor) bool {
	return true
}

func (s *Multi) IsSelected(cursor Cursor) bool {
	return s.Selected[cursor]
}

type ColumnBased struct {
	Left      map[Row]bool
	Right     map[Row]bool
	ActiveCol Col
} // dual list

func (s *ColumnBased) Toggle(cursor Cursor) {}

func (s *ColumnBased) Select(cursor Cursor) {}

func (s *ColumnBased) IsSelectable(cursor Cursor) bool {
	return false
}

func (s *ColumnBased) IsSelected(cursor Cursor) bool {
	if cursor.Col == 0 {
		return s.Left[cursor.Row]
	}
	if cursor.Col == 1 {
		return s.Right[cursor.Row]
	}
	return false
}

type NoSelection struct{}                   // pure cursor navigation
func (s *NoSelection) Toggle(cursor Cursor) {}

func (s *NoSelection) Select(cursor Cursor) {}

func (s *NoSelection) IsSelectable(cursor Cursor) bool {
	return false
}

func (s *NoSelection) IsSelected(cursor Cursor) bool {
	return false
}
