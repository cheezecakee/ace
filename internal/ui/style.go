package ui

import (
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	Header lipgloss.Style
	Body   lipgloss.Style
	Footer lipgloss.Style

	Popup        lipgloss.Style
	PopupTitle   lipgloss.Style
	PopupContent lipgloss.Style

	Accent lipgloss.Style
	Muted  lipgloss.Style
}

func NewStyles(width int) Styles {
	return Styles{
		Header: lipgloss.NewStyle().
			Width(width).
			Align(lipgloss.Center).
			Padding(0, 1),

		Body: lipgloss.NewStyle().
			Width(width).
			Padding(1, 1),

		Footer: lipgloss.NewStyle().
			Width(width).
			Padding(0, 1),

		Popup: lipgloss.NewStyle().
			Width(40).
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2),

		Accent: lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")),
	}
}
