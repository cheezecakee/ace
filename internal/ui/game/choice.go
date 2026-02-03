package game

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/cheezecakee/ace/internal/engine"
	ctx "github.com/cheezecakee/ace/internal/ui/context"
	"github.com/cheezecakee/ace/internal/ui/widgets"
)

type ChoiceQuestionUI struct {
	q      engine.ChoiceQuestion
	ctx    *ctx.Context
	widget *widgets.Widget
}

func NewChoiceQuestionUI(q engine.ChoiceQuestion, c *ctx.Context) *ChoiceQuestionUI {
	items := make([]widgets.Item, len(q.Options))
	for i, opt := range q.Options {
		items[i] = widgets.NewTextItem(opt)
	}

	return &ChoiceQuestionUI{
		q:      q,
		ctx:    c,
		widget: widgets.NewChoiceWidget(items),
	}
}

func (c *ChoiceQuestionUI) Init() tea.Cmd {
	return nil
}

func (c *ChoiceQuestionUI) Update(msg tea.Msg) tea.Cmd {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if dir, ok := widgets.DirectionFromKey(keyMsg, c.ctx.Keys.Up, c.ctx.Keys.Down, c.ctx.Keys.Left, c.ctx.Keys.Right); ok {
			c.widget.Move(dir)
		}
	}
	return nil
}

func (c *ChoiceQuestionUI) View() string {
	return c.widget.Render()
}

func (c *ChoiceQuestionUI) Submit() (engine.Answer, bool) {
	return engine.ChoiceAnswer{
		Selected: int(c.widget.Cursor.Row),
	}, true
}
