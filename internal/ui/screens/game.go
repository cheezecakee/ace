package screens

import (
	tea "github.com/charmbracelet/bubbletea"

	ctx "github.com/cheezecakee/ace/internal/ui/context"
	"github.com/cheezecakee/ace/internal/ui/game"
)

type GameScreen struct {
	ctx  *ctx.Context
	game *game.Screen
}

func NewGameScreen(ctx *ctx.Context) Screen {
	return &GameScreen{
		ctx:  ctx,
		game: game.NewScreen(ctx),
	}
}

func (m *GameScreen) Init() tea.Cmd {
	return m.game.Init()
}

func (m *GameScreen) Update(msg tea.Msg) (Screen, tea.Cmd) {
	done, cmd := m.game.Update(msg)
	if done {
		return NewCompleteScreen(m.ctx), nil
	}
	return m, cmd
}

func (m *GameScreen) View() string {
	return m.game.View()
}
