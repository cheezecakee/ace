// Package screens
package screens

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/cheezecakee/ace/internal/engine"
	"github.com/cheezecakee/ace/internal/ui"
	ctx "github.com/cheezecakee/ace/internal/ui/context"
)

type Screen interface {
	Update(msg tea.Msg) (Screen, tea.Cmd)
	View() string
	Init() tea.Cmd // Optional initilization
}

type Model struct {
	currentScreen Screen

	// Shared state that all screens might need
	ctx *ctx.Context
}

func NewModel(ctx *ctx.Context) *Model {
	return &Model{
		currentScreen: NewMenu(ctx), // Start with menu
		ctx:           ctx,
	}
}

func (m *Model) Init() tea.Cmd {
	return m.currentScreen.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle global keys (like quit)
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.ctx.Width = msg.Width
		m.ctx.Height = msg.Height
		m.ctx.Styles = ui.NewStyles(msg.Width)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.ctx.Keys.Quit):
			return m, tea.Quit
		}

	case ctx.SetFormatMsg:
		m.ctx.Format = engine.Format(msg)
	}

	newScreen, cmd := m.currentScreen.Update(msg)
	m.currentScreen = newScreen
	return m, cmd
}

func (m *Model) View() string {
	return m.currentScreen.View()
}
