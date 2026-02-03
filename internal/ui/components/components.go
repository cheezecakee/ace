package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

func LivesView(lives int) string {
	if lives <= 0 {
		return ""
	}
	return strings.Repeat("❤️ ", lives)
}

func TimerView(d time.Duration) string {
	if d <= 0 {
		return "⏱ 0:00"
	}

	min := int(d.Minutes())
	sec := int(d.Seconds()) % 60
	return fmt.Sprintf("⏱ %d:%02d", min, sec)
}

func PaginationView(current, total int, answered []bool) string {
	var b strings.Builder

	// Debug: show the raw data
	// fmt.Printf("DEBUG Pagination - current: %d, total: %d, answered: %v\n", current, total, answered)

	for i := range total {
		switch {
		case i == current:
			b.WriteString("[●]") // Current question - filled with brackets
		case answered != nil && i < len(answered) && answered[i]:
			b.WriteString(" ● ") // Answered question - filled
		default:
			b.WriteString(" ○ ") // Unanswered question - empty
		}
		b.WriteString(" ") // Add space between dots
	}

	return b.String() + "\n"
}

func OptionView(options []string, selected []int, multiple bool) string {
	if len(options) == 2 {
		return boolView(selected)
	}

	if multiple {
		return multipleView(options, selected)
	}

	return choiceView(options, selected)
}

func choiceView(options []string, selected []int) string {
	var b strings.Builder

	for i, opt := range options {
		cursor := " "
		if i == selected[0] {
			cursor = ">"
		}

		label := string(rune('A' + i))
		b.WriteString(cursor + " [" + label + "] " + opt + "\n")
	}

	return b.String()
}

func multipleView(options []string, selected []int) string {
	var b strings.Builder

	selectedMap := make(map[int]bool, len(selected))
	for _, i := range selected {
		selectedMap[i] = true
	}

	for i, opt := range options {
		marker := " "
		if selectedMap[i] {
			marker = "✓"
		}

		label := string(rune('A' + i))
		b.WriteString("[" + marker + "] [" + label + "] " + opt + "\n")
	}
	return b.String()
}

func boolView(selected []int) string {
	var b strings.Builder

	trueMark := " "
	falseMark := " "

	switch selected[0] {
	case 0:
		trueMark = "✓"

	case 1:
		falseMark = "✓"
	}

	b.WriteString("(" + trueMark + ") True\n")
	b.WriteString("(" + falseMark + ") False\n")

	return b.String()
}

func TextEntryAnswerView(input string) string {
	return "Answer:\n" + input + "\n"
}

func QuestionView(q string) string {
	return q + "\n\n"
}

func NewTextArea(placeholder string, focused bool) textarea.Model {
	ta := textarea.New()
	ta.Placeholder = placeholder
	ta.Focus()
	ta.CharLimit = 500
	ta.SetWidth(50)
	ta.SetHeight(4)

	if focused {
		ta.Focus()
		ta.FocusedStyle.Base.BorderStyle(lipgloss.NormalBorder())
	} else {
		ta.Blur()
		ta.FocusedStyle.Base.BorderStyle(lipgloss.HiddenBorder())
	}

	return ta
}
