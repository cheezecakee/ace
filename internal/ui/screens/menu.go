package screens

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/cheezecakee/ace/internal/engine"
	"github.com/cheezecakee/ace/internal/ui/context"
	"github.com/cheezecakee/ace/internal/ui/widgets"
)

type Menu struct {
	widget *widgets.Widget
	ctx    *context.Context
}

func NewMenu(ctx *context.Context) Screen {
	// Items have their actions built-in!
	items := [][]widgets.Item{
		{
			widgets.NewButtonItem("Standard", func() any {
				return ModeDifficultyScreen(engine.StandardMode)(ctx)
			}),
			widgets.NewButtonItem("Rapid", func() any {
				return ModeDifficultyScreen(engine.RapidMode)(ctx)
			}),
			widgets.NewButtonItem("Hardcore", func() any {
				return ModeDifficultyScreen(engine.HardcoreMode)(ctx)
			}),
			widgets.NewButtonItem("Custom", func() any {
				return ModeDifficultyScreen(engine.CustomMode)(ctx)
			}),
		},
		{widgets.NewButtonItem("Quick Start", func() any {
			return NewDifficultyScreen(ctx)
		})},
		{widgets.NewButtonItem("Packs", func() any {
			return NewPacksScreen(ctx)
		})},
		{widgets.NewButtonItem("Stats", func() any {
			// TODO: implement stats screen
			return NewMenu(ctx)
		})},
		{widgets.NewButtonItem("Settings", func() any {
			return NewSettingsScreen(ctx)
		})},
	}

	return &Menu{
		widget: widgets.NewGrid(items, 4),
		ctx:    ctx,
	}
}

func (m *Menu) Update(msg tea.Msg) (Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle navigation
		if dir, ok := widgets.DirectionFromKey(msg, m.ctx.Keys.Up, m.ctx.Keys.Down, m.ctx.Keys.Left, m.ctx.Keys.Right); ok {
			m.widget.Move(dir)
			return m, nil
		}

		// Handle selection - execute the item's action!
		if key.Matches(msg, m.ctx.Keys.Submit) {
			item, ok := m.widget.GetItem()
			if ok && item.Action != nil {
				result := item.Action.Exec()
				if screen, ok := result.(Screen); ok {
					return screen, nil
				}
			}
		}
	}

	return m, nil
}

func (m *Menu) View() string {
	var s strings.Builder
	s.WriteString("Menu\n\n")
	s.WriteString(m.widget.Render())
	s.WriteString("\n\n")
	return s.String()
}

func (m *Menu) Init() tea.Cmd {
	return nil
}
