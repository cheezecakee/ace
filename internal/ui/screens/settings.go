package screens

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/cheezecakee/ace/internal/ui/context"
	"github.com/cheezecakee/ace/internal/ui/widgets"
)

type SettingsScreen struct {
	widget *widgets.Widget
	ctx    *context.Context
}

func NewSettingsScreen(ctx *context.Context) Screen {
	items := []widgets.Item{
		widgets.NewTextItem("Language"),
		widgets.NewTextItem("Verify/Repair"),
		widgets.NewTextItem("Reset"),
	}

	return &SettingsScreen{
		widget: widgets.NewList(items),
		ctx:    ctx,
	}
}

func (m *SettingsScreen) Init() tea.Cmd {
	return nil
}

func (m *SettingsScreen) Update(msg tea.Msg) (Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		if dir, ok := widgets.DirectionFromKey(msg, m.ctx.Keys.Up, m.ctx.Keys.Down, m.ctx.Keys.Left, m.ctx.Keys.Right); ok {
			m.widget.Move(dir)
			return m, nil
		}

		if key.Matches(msg, m.ctx.Keys.Submit) {
			// For now, just go back to menu as a test
			return NewMenu(m.ctx), nil
		}

		if key.Matches(msg, m.ctx.Keys.Back) {
			return NewMenu(m.ctx), nil
		}
	}

	return m, nil
}

func (m *SettingsScreen) View() string {
	var s strings.Builder
	s.WriteString("Select Settings\n\n")
	s.WriteString(m.widget.Render())
	return s.String()
}
