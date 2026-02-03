package screens

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/cheezecakee/ace/internal/ui/context"
	"github.com/cheezecakee/ace/internal/ui/widgets"
)

type ImportScreen struct {
	widget *widgets.Widget
	ctx    *context.Context
}

func NewImportScreen(ctx *context.Context) Screen {
	items := []widgets.Item{
		widgets.NewTextItem("Import pack from file"),
	}

	return &ImportScreen{
		widget: widgets.NewList(items),
		ctx:    ctx,
	}
}

func (m *ImportScreen) Init() tea.Cmd {
	return nil
}

func (m *ImportScreen) Update(msg tea.Msg) (Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		if dir, ok := widgets.DirectionFromKey(msg, m.ctx.Keys.Up, m.ctx.Keys.Down, m.ctx.Keys.Left, m.ctx.Keys.Right); ok {
			m.widget.Move(dir)
			return m, nil
		}

		if key.Matches(msg, m.ctx.Keys.Submit) {
			// TODO: bubbles file picker goes here
			return m, nil
		}

		if key.Matches(msg, m.ctx.Keys.Back) {
			return NewPacksScreen(m.ctx), nil
		}
	}

	return m, nil
}

func (m *ImportScreen) View() string {
	var s strings.Builder
	s.WriteString("Import Pack\n\n")
	s.WriteString(m.widget.Render())
	s.WriteString("\n\n(file picker coming soon)")
	return s.String()
}
