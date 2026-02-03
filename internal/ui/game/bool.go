package game

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/cheezecakee/ace/internal/engine"
	ctx "github.com/cheezecakee/ace/internal/ui/context"
	"github.com/cheezecakee/ace/internal/ui/widgets"
)

type BoolQuestionUI struct {
	q      engine.BoolQuestion
	ctx    *ctx.Context
	widget *widgets.Widget
}

func NewBoolQuestionUI(q engine.BoolQuestion, c *ctx.Context) *BoolQuestionUI {
	return &BoolQuestionUI{
		q:      q,
		ctx:    c,
		widget: widgets.NewBoolWidget(),
	}
}

func (b *BoolQuestionUI) Init() tea.Cmd {
	return nil
}

func (b *BoolQuestionUI) Update(msg tea.Msg) tea.Cmd {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if dir, ok := widgets.DirectionFromKey(keyMsg, b.ctx.Keys.Up, b.ctx.Keys.Down, b.ctx.Keys.Left, b.ctx.Keys.Right); ok {
			b.widget.Move(dir)
		}
	}
	return nil
}

func (b *BoolQuestionUI) View() string {
	return b.widget.Render()
}

func (b *BoolQuestionUI) Submit() (engine.Answer, bool) {
	return engine.BoolAnswer{
		Answer: b.widget.Cursor.Row == 0, // Row 0 = True, Row 1 = False
	}, true
}
