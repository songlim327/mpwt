package tui

import (
	"mpwt/internal/core"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tui struct {
	width    int
	height   int
	viewport string
	status   *status
	footer   *footer
	option   *option
	execute  *execute
}

// viewportMsg represents a message struct to trigger main window view changes
type viewportMsg struct {
	viewport string
}

func newTui(tc *core.TerminalConfig) *tui {
	return &tui{
		viewport: Main,
		status:   newStatus(""),
		footer:   newFooter(),
		option:   newOption(),
		execute:  newExecute(tc),
	}
}

func (t *tui) Init() tea.Cmd {
	return tea.SetWindowTitle("🍊 MPWT")
}

func (t *tui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		t.width = msg.Width
		t.height = msg.Height

	case viewportMsg:
		t.viewport = msg.viewport
		if t.viewport == Execute {
			s, cmd := t.status.Update(statusMsg{message: "Each line of command will spawn a new pane in terminal"})
			t.status = s.(*status)
			return t, cmd
		}

	case statusMsg:
		s, cmd := t.status.Update(msg)
		t.status = s.(*status)
		return t, cmd

	case tea.KeyMsg:
		switch t.viewport {
		case Execute:
			e, cmd := t.execute.Update(msg)
			t.execute = e.(*execute)
			return t, cmd
		default:
			o, cmd := t.option.Update(msg)
			t.option = o.(*option)
			return t, cmd
		}
	}
	return t, nil
}

func (t *tui) View() string {
	margin := 2
	padding := 1
	gap := 1

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(BorderForegroundColor))

	// Window and children size calculations
	boxWidth := t.width - margin*2
	boxHeight := t.height - margin*2 - boxStyle.GetBorderTopSize()*2 - t.status.style.GetHeight() - t.status.style.GetBorderTopSize()*2

	t.status.width = boxWidth
	t.footer.width = boxWidth - padding*2

	var view string
	switch t.viewport {
	case Execute:
		t.execute.width = boxWidth - padding*2
		t.execute.height = boxHeight - padding - t.footer.style.GetHeight()
		view = t.execute.View()
	default:
		t.option.width = boxWidth - padding*2
		t.option.height = boxHeight - padding - t.footer.style.GetHeight()
		view = t.option.View()
	}

	// Content box
	return lipgloss.NewStyle().
		Margin(margin).
		Render(
			t.status.View(),
			boxStyle.
				Width(boxWidth).
				Height(boxHeight).
				Padding(padding, padding, 0, padding).
				MarginTop(gap).Render(
				lipgloss.JoinVertical(lipgloss.Left,
					view,
					t.footer.View(),
				),
			),
		)

}

// InitTea intialize a new tea program with user interactions
func InitTea(tc *core.TerminalConfig) error {
	p := tea.NewProgram(
		newTui(tc),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
