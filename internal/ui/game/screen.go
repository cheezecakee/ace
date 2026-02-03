package game

import (
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/cheezecakee/ace/internal/engine"
	"github.com/cheezecakee/ace/internal/ui"
	"github.com/cheezecakee/ace/internal/ui/components"
	ctx "github.com/cheezecakee/ace/internal/ui/context"
)

type Screen struct {
	ctx *ctx.Context

	questionUI QuestionUI
	help       help.Model

	lastTick    time.Time
	gameStarted bool
}

type TickMsg time.Time

func NewScreen(c *ctx.Context) *Screen {
	q := c.Session.GetCurrentQuestion()

	return &Screen{
		ctx:         c,
		questionUI:  NewQuestionUI(q, c),
		help:        help.New(),
		lastTick:    time.Now(),
		gameStarted: false,
	}
}

func (s *Screen) Init() tea.Cmd {
	return tea.Batch(
		s.questionUI.Init(),
		s.tickCmd(),
	)
}

func (s *Screen) tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (s *Screen) Update(msg tea.Msg) (bool, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case TickMsg:
		if s.gameStarted {
			now := time.Time(msg)
			elapsed := now.Sub(s.lastTick)
			s.lastTick = now

			format := s.ctx.Session.GetFormat()
			if format.Time.Control != engine.Unlimited {
				expired, state := s.ctx.Session.Tick(elapsed)
				if expired || state.IsTerminal() {
					return true, nil
				}
			}
		}
		cmds = append(cmds, s.tickCmd())

	case tea.KeyMsg:
		if !s.gameStarted {
			s.gameStarted = true
			s.lastTick = time.Now()
			return false, nil
		}

		switch {
		case key.Matches(msg, s.ctx.Keys.NextQuestion):
			if !s.ctx.Session.CanNavigateBack() {
				return s.submit()
			}
			if err := s.ctx.Session.NextQuestion(); err == nil {
				s.loadCurrentQuestion()
			}

		case key.Matches(msg, s.ctx.Keys.PrevQuestion):
			if s.ctx.Session.CanNavigateBack() {
				if err := s.ctx.Session.PrevQuestion(); err == nil {
					s.loadCurrentQuestion()
				}
			}

		case key.Matches(msg, s.ctx.Keys.Submit):
			return s.submit()

		case key.Matches(msg, s.ctx.Keys.Help):
			s.help.ShowAll = !s.help.ShowAll
		}

		cmds = append(cmds, s.questionUI.Update(msg))
	}

	return false, tea.Batch(cmds...)
}

func (s *Screen) submit() (bool, tea.Cmd) {
	ans, ok := s.questionUI.Submit()
	if !ok {
		return false, nil
	}

	_ = s.ctx.Session.SubmitAnswer(ans)

	if s.ctx.Session.IsCompleted() {
		return true, nil
	}

	if !s.ctx.Session.CanNavigateBack() {
		if err := s.ctx.Session.NextQuestion(); err == nil {
			s.loadCurrentQuestion()
		}
	}

	return false, nil
}

func (s *Screen) loadCurrentQuestion() {
	q := s.ctx.Session.GetCurrentQuestion()
	s.questionUI = NewQuestionUI(q, s.ctx)
}

func (s *Screen) View() string {
	r := ui.NewRender(s.ctx.Styles)

	q := s.ctx.Session.GetCurrentQuestion()
	index := s.ctx.Session.GetCurrentIndex()
	format := s.ctx.Session.GetFormat()

	// Header
	var lives, timer string
	if format.Lives.Enabled {
		lives = components.LivesView(s.ctx.Session.GetLivesRemaining())
	}
	if format.Time.Control != engine.Unlimited {
		timer = components.TimerView(s.ctx.Session.GetTimeRemaining())
	}

	header := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Width(s.ctx.Width/2).Align(lipgloss.Left).Render(lives),
		lipgloss.NewStyle().Width(s.ctx.Width/2).Align(lipgloss.Right).Render(timer),
	)
	r.Head = r.Styles.Header.Render(header)

	// Body
	r.Body = components.QuestionView(q.GetPrompt())
	r.Body += s.questionUI.View() + "\n"

	// Footer
	total := len(s.ctx.Session.GetResults().Answers)
	answered := make([]bool, total)
	for i := range answered {
		answered[i] = s.ctx.Session.IsAnswered(i)
	}

	// fmt.Printf("answered: %v\n", answered)

	r.Footer = components.PaginationView(index, total, answered)
	r.Footer += "\n\n" + s.help.View(s.ctx.Keys)

	return r.Build()
}
