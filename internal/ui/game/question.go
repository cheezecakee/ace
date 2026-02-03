package game

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/cheezecakee/ace/internal/engine"
)

type QuestionUI interface {
	Init() tea.Cmd
	Update(msg tea.Msg) tea.Cmd
	View() string

	// Returns (answer, ok)
	// ok == false, means "not ready to submit"
	Submit() (engine.Answer, bool)
}
