package widgets

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// DirectionFromKey converts a key message to a direction using key bindings
func DirectionFromKey(msg tea.KeyMsg, up, down, left, right key.Binding) (Direction, bool) {
	switch {
	case key.Matches(msg, up):
		return Top, true
	case key.Matches(msg, down):
		return Down, true
	case key.Matches(msg, left):
		return Left, true
	case key.Matches(msg, right):
		return Right, true
	default:
		return Direction{}, false
	}
}
