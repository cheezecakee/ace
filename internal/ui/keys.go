package ui

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	// Global
	Quit   key.Binding
	Back   key.Binding
	Submit key.Binding

	// Navigation
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding

	Select key.Binding

	// Question interaction
	ToggleFocus  key.Binding // For text entry
	NextQuestion key.Binding
	PrevQuestion key.Binding

	// Toggle help
	Help key.Binding
}

func DefaultKeyMap() KeyMap {
	// Define all the keys with help text
	return KeyMap{
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c", "ctrl+q"),
			key.WithHelp("ctrl+q", "quit"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Submit: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "submit"),
		),

		NextQuestion: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next question"),
		),
		PrevQuestion: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "prev question"),
		),
		Up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("↓/j", "move down"),
		),
		Left: key.NewBinding(
			key.WithKeys("h", "left"),
			key.WithHelp("←/h", "move left"),
		),
		Right: key.NewBinding(
			key.WithKeys("l", "right"),
			key.WithHelp("→/l", "move right"),
		),

		Select: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "select"),
		),
		ToggleFocus: key.NewBinding(
			key.WithKeys("shift+enter"),
			key.WithHelp("shift+enter", "toggle focus"),
		),

		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
	}
}
