package screens

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/cheezecakee/ace/internal/ui/context"
	"github.com/cheezecakee/ace/internal/ui/widgets"
)

type GameOverScreen struct {
	widget *widgets.Widget
	ctx    *context.Context
}

func NewGameOverScreen(ctx *context.Context) Screen {
	items := []widgets.Item{
		widgets.NewTextItem("Back to menu"),
	}

	return &GameOverScreen{
		widget: widgets.NewList(items),
		ctx:    ctx,
	}
}

func (m *GameOverScreen) Init() tea.Cmd {
	return nil
}

func (m *GameOverScreen) Update(msg tea.Msg) (Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		if dir, ok := widgets.DirectionFromKey(msg, m.ctx.Keys.Up, m.ctx.Keys.Down, m.ctx.Keys.Left, m.ctx.Keys.Right); ok {
			m.widget.Move(dir)
			return m, nil
		}

		if key.Matches(msg, m.ctx.Keys.Submit) {
			return NewMenu(m.ctx), nil
		}
	}

	return m, nil
}

func (m *GameOverScreen) View() string {
	var s strings.Builder
	s.WriteString("🎉 Game Over!\n\n")
	s.WriteString("Oh no! You failed!\n\n")
	s.WriteString(m.widget.Render())
	return s.String()
}
