package screens

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/cheezecakee/ace/internal/engine"
	"github.com/cheezecakee/ace/internal/ui/context"
	"github.com/cheezecakee/ace/internal/ui/widgets"
)

type DifficultyScreen struct {
	difficulties []engine.Difficulty
	widget       *widgets.Widget
	ctx          *context.Context
}

func NewDifficultyScreen(ctx *context.Context) Screen {
	difficulties := []engine.Difficulty{
		engine.Entry,
		engine.Junior,
		engine.Mid,
		engine.Senior,
	}

	items := make([]widgets.Item, 0, len(difficulties))
	for _, d := range difficulties {
		items = append(items, widgets.NewTextItem(d.String()))
	}

	return &DifficultyScreen{
		difficulties: difficulties,
		widget:       widgets.NewList(items),
		ctx:          ctx,
	}
}

func (m *DifficultyScreen) Init() tea.Cmd {
	return nil
}

func (m *DifficultyScreen) Update(msg tea.Msg) (Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if dir, ok := widgets.DirectionFromKey(msg, m.ctx.Keys.Up, m.ctx.Keys.Down, m.ctx.Keys.Left, m.ctx.Keys.Right); ok {
			m.widget.Move(dir)
			return m, nil
		}

		if key.Matches(msg, m.ctx.Keys.Submit) {
			row := int(m.widget.Cursor.Row)
			selected := m.difficulties[row]

			mode := engine.GetGameMode(m.ctx.Mode)
			m.ctx.Format = mode.Format(selected)

			return NewRoleScreen(m.ctx), nil
		}

		if key.Matches(msg, m.ctx.Keys.Back) {
			return NewMenu(m.ctx), nil
		}
	}

	return m, nil
}

func (m *DifficultyScreen) View() string {
	var s strings.Builder
	s.WriteString("Select Difficulty\n\n")
	s.WriteString(m.widget.Render())
	return s.String()
}

func ModeDifficultyScreen(mode engine.ModeID) func(*context.Context) Screen {
	return func(ctx *context.Context) Screen {
		ctx.Mode = mode
		return NewDifficultyScreen(ctx)
	}
}
