package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/cheezecakee/ace/internal/engine"
	"github.com/cheezecakee/ace/internal/pack"
	"github.com/cheezecakee/ace/internal/session"
	"github.com/cheezecakee/ace/internal/ui/context"
	"github.com/cheezecakee/ace/internal/ui/widgets"
)

type RoleScreen struct {
	roles  []pack.Role
	widget *widgets.Widget
	ctx    *context.Context
}

func NewRoleScreen(ctx *context.Context) Screen {
	packTypes := pack.FromEngineTypes(ctx.Format.Question.Types)

	// Only get roles that have questions available for current difficulty and types
	roles := ctx.LookupCache.GetAvailableRoles(
		ctx.Format.Progression.Difficulty,
		packTypes,
	)

	items := make([]widgets.Item, 0, len(roles))
	for _, role := range roles {
		items = append(items, widgets.NewTextItem(string(role)))
	}

	return &RoleScreen{
		roles:  roles,
		widget: widgets.NewList(items),
		ctx:    ctx,
	}
}

func (m *RoleScreen) Init() tea.Cmd {
	return nil
}

func (m *RoleScreen) Update(msg tea.Msg) (Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		if dir, ok := widgets.DirectionFromKey(msg, m.ctx.Keys.Up, m.ctx.Keys.Down, m.ctx.Keys.Left, m.ctx.Keys.Right); ok {
			m.widget.Move(dir)
			return m, nil
		}
		if key.Matches(msg, m.ctx.Keys.Submit) {
			row := int(m.widget.Cursor.Row)
			role := m.roles[row]

			packTypes := pack.FromEngineTypes(m.ctx.Format.Question.Types)

			questionIDs := m.ctx.LookupCache.GetQuestionIDs(
				m.ctx.Format.Progression.Difficulty,
				role,
				packTypes,
			)

			if len(questionIDs) == 0 {
				fmt.Println("error: no questions found")
				return m, nil
			}

			questions := m.ctx.QuestionCache.Fetch(questionIDs)
			grader := engine.GetGrader(m.ctx.Mode)

			sess := session.NewSession(
				m.ctx.Format,
				questions,
				grader,
			)

			if err := sess.Begin(); err != nil {
				return m, nil
			}

			m.ctx.Session = sess
			return NewGameScreen(m.ctx), nil
		}

		if key.Matches(msg, m.ctx.Keys.Back) {
			return NewDifficultyScreen(m.ctx), nil
		}
	}

	return m, nil
}

func (m *RoleScreen) View() string {
	var s strings.Builder
	s.WriteString("Select Role\n\n")
	s.WriteString(m.widget.Render())
	return s.String()
}
