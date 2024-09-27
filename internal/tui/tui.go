package tui

import (
	"mpwt/internal/core"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type mainWindow struct {
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

func initialModel(tc *core.TerminalConfig) mainWindow {
	return mainWindow{
		viewport: Main,
		status:   newStatus(""),
		footer:   newFooter(),
		option:   newOption(),
		execute:  newExecute(tc),
	}
}

func (m mainWindow) Init() tea.Cmd {
	return textarea.Blink
}

func (m mainWindow) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case viewportMsg:
		m.viewport = msg.viewport
		if m.viewport == Execute {
			s, cmd := m.status.Update(statusMsg{message: "Each line of command will spawn a new pane in terminal"})
			m.status = s.(*status)
			return m, cmd
		}

	case statusMsg:
		s, cmd := m.status.Update(msg)
		m.status = s.(*status)
		return m, cmd

	case tea.KeyMsg:
		switch m.viewport {
		case Execute:
			e, cmd := m.execute.Update(msg)
			m.execute = e.(*execute)
			return m, cmd
		default:
			o, cmd := m.option.Update(msg)
			m.option = o.(*option)
			return m, cmd
		}
	}
	return m, nil
}

func (m mainWindow) View() string {
	margin := 2
	padding := 1
	gap := 1

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(BorderForegroundColor))

	// Window and children size calculations
	boxWidth := m.width - margin*2
	boxHeight := m.height - margin*2 - boxStyle.GetBorderTopSize()*2 - m.status.style.GetHeight() - m.status.style.GetBorderTopSize()*2

	m.status.width = boxWidth
	m.footer.width = boxWidth - padding*2

	var view string
	switch m.viewport {
	case Execute:
		m.execute.width = boxWidth - padding*2
		m.execute.height = boxHeight - padding - m.footer.style.GetHeight()
		view = m.execute.View()
	default:
		m.option.width = boxWidth - padding*2
		m.option.height = boxHeight - padding - m.footer.style.GetHeight()
		view = m.option.View()
	}

	// Content box
	return lipgloss.NewStyle().
		Margin(margin).
		Render(
			m.status.View(),
			boxStyle.
				Width(boxWidth).
				Height(boxHeight).
				Padding(padding, padding, 0, padding).
				MarginTop(gap).Render(
				lipgloss.JoinVertical(lipgloss.Left,
					view,
					m.footer.View(),
				),
			),
		)

}

// InitTea intialize a new tea program with user interactions
func InitTea(tc *core.TerminalConfig) error {
	p := tea.NewProgram(
		initialModel(tc),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		return err
	}

	return nil
}
