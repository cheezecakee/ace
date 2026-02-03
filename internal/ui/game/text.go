package game

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/cheezecakee/ace/internal/engine"
	ctx "github.com/cheezecakee/ace/internal/ui/context"
)

type TextQuestionUI struct {
	q        engine.TextEntryQuestion
	ctx      *ctx.Context
	textArea textarea.Model
}

func NewTextQuestionUI(q engine.TextEntryQuestion, c *ctx.Context) *TextQuestionUI {
	ta := textarea.New()
	ta.Placeholder = "Type your answer..."
	ta.CharLimit = 500
	ta.SetWidth(50)
	ta.SetHeight(4)
	ta.Focus() // Start focused

	return &TextQuestionUI{
		q:        q,
		ctx:      c,
		textArea: ta,
	}
}

func (t *TextQuestionUI) Init() tea.Cmd {
	return t.textArea.Cursor.BlinkCmd()
}

func (t *TextQuestionUI) Update(msg tea.Msg) tea.Cmd {
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		// Handle focus toggle
		if key.Matches(keyMsg, t.ctx.Keys.ToggleFocus) {
			if t.textArea.Focused() {
				t.textArea.Blur()
			} else {
				t.textArea.Focus()
			}
			return nil
		}
	}

	// Update textarea
	var cmd tea.Cmd
	t.textArea, cmd = t.textArea.Update(msg)
	return cmd
}

func (t *TextQuestionUI) View() string {
	return t.textArea.View()
}

func (t *TextQuestionUI) Submit() (engine.Answer, bool) {
	text := t.textArea.Value()

	if text == "" {
		// Not ready to submit
		return nil, false
	}

	return engine.TextEntryAnswer{
		Text: text,
	}, true
}
