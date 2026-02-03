package screens

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/cheezecakee/ace/internal/ui/context"
	"github.com/cheezecakee/ace/internal/ui/widgets"
)

type PacksScreen struct {
	widget  *widgets.DualColumnWidget
	packIDs []string // Maps row index to pack ID
	ctx     *context.Context
}

func NewPacksScreen(ctx *context.Context) Screen {
	leftItems := make([]widgets.Item, 0)
	rightItems := make([]widgets.Item, 0)
	packIDs := make([]string, 0)

	// Row 0: Import button in left column
	leftItems = append(leftItems, widgets.NewTextItem("Import"))
	rightItems = append(rightItems, widgets.NewTextItem(""))
	packIDs = append(packIDs, "") // Sentinel for import row

	// Add packs - inactive in left, active in right
	for _, p := range ctx.Metadata.Packs {
		packIDs = append(packIDs, p.ID)

		if ctx.Packs[p.ID] {
			// Active pack - right column
			leftItems = append(leftItems, widgets.NewTextItem(""))
			rightItems = append(rightItems, widgets.NewTextItem(p.Name))
		} else {
			// Inactive pack - left column
			leftItems = append(leftItems, widgets.NewTextItem(p.Name))
			rightItems = append(rightItems, widgets.NewTextItem(""))
		}
	}

	return &PacksScreen{
		widget:  widgets.NewDualColumn(leftItems, rightItems),
		packIDs: packIDs,
		ctx:     ctx,
	}
}

func (m *PacksScreen) Init() tea.Cmd {
	return nil
}

func (m *PacksScreen) Update(msg tea.Msg) (Screen, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle navigation
		if dir, ok := widgets.DirectionFromKey(msg, m.ctx.Keys.Up, m.ctx.Keys.Down, m.ctx.Keys.Left, m.ctx.Keys.Right); ok {
			m.widget.Move(dir)
			return m, nil
		}

		// Handle selection/toggle (Space or Enter)
		if key.Matches(msg, m.ctx.Keys.Select) || key.Matches(msg, m.ctx.Keys.Submit) {
			row := int(m.widget.Cursor.Row)

			// Special case: Import button
			if row == 0 && m.widget.Cursor.Col == 0 {
				return NewImportScreen(m.ctx), nil
			}

			// Toggle pack between columns
			packID := m.packIDs[row]
			if packID != "" {
				// Update context
				m.ctx.Packs[packID] = !m.ctx.Packs[packID]

				// Move item in widget
				m.widget.MoveItem()
			}

			return m, nil
		}

		// Back to menu
		if key.Matches(msg, m.ctx.Keys.Back) {
			// Save active packs
			m.ctx.User.Settings.ActivePacks = m.ctx.GetActivePacks()
			_ = m.ctx.User.Save()
			_ = m.ctx.RebuildCache()

			return NewMenu(m.ctx), nil
		}
	}

	return m, nil
}

func (m *PacksScreen) View() string {
	var s strings.Builder
	s.WriteString("Packs\n\n")
	s.WriteString("Available            Active\n")
	s.WriteString("─────────            ──────\n")
	s.WriteString(m.widget.Render())
	s.WriteString("\n\nSpace/Enter: Toggle | Esc: Back")
	return s.String()
}
