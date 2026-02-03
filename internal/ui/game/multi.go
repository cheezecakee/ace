package game

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/cheezecakee/ace/internal/engine"
	ctx "github.com/cheezecakee/ace/internal/ui/context"
	"github.com/cheezecakee/ace/internal/ui/widgets"
)

type MultiChoiceQuestionUI struct {
	q      engine.MultipleChoiceQuestion
	ctx    *ctx.Context
	widget *widgets.Widget
}

func NewMultiChoiceQuestionUI(q engine.MultipleChoiceQuestion, c *ctx.Context) *MultiChoiceQuestionUI {
	items := make([]widgets.Item, len(q.Options))
	for i, opt := range q.Options {
		items[i] = widgets.NewCheckItem(opt)
	}

	return &MultiChoiceQuestionUI{
		q:      q,
		ctx:    c,
		widget: widgets.NewMultipleChoiceWidget(items),
	}
}

func (m *MultiChoiceQuestionUI) Init() tea.Cmd {
	return nil
}

func (m *MultiChoiceQuestionUI) Update(msg tea.Msg) tea.Cmd {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		// Handle movement
		if dir, ok := widgets.DirectionFromKey(keyMsg, m.ctx.Keys.Up, m.ctx.Keys.Down, m.ctx.Keys.Left, m.ctx.Keys.Right); ok {
			m.widget.Move(dir)
			return nil
		}

		// Handle toggle (space)
		if key.Matches(keyMsg, m.ctx.Keys.Select) {
			m.widget.Toggle()
		}
	}
	return nil
}

func (m *MultiChoiceQuestionUI) View() string {
	return m.widget.Render()
}

func (m *MultiChoiceQuestionUI) Submit() (engine.Answer, bool) {
	// Get all selected items using the widget helper
	selectedItems := m.widget.GetSelectedItems()

	if len(selectedItems) == 0 {
		// No selection - not ready to submit
		return nil, false
	}

	// Convert selected items to indices
	selected := make([]int, 0, len(selectedItems))

	// Get selection state from widget
	if multiSel, ok := m.widget.Selection.(*widgets.Multi); ok {
		for cursor := range multiSel.Selected {
			if multiSel.Selected[cursor] {
				selected = append(selected, int(cursor.Row))
			}
		}
	}

	return engine.MultipleChoiceAnswer{
		Selected: selected,
	}, true
}
