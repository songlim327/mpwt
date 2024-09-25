package tui

import (
	"mpwt/internal/core"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// keyMap defines a set of keybindings.
type keyMap struct {
	Launch key.Binding
	Quit   key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view
// It is part of the key.Map interface
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Launch, k.Quit}
}

// FullHelp returns keybindings to be shown in the full help view
// It is part of the key.Map interface
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Launch, k.Quit},
	}
}

// keys implements key.Map interface and defines keyMap  of help menu
var keys = keyMap{
	Launch: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "launch"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc"),
		key.WithHelp("esc/q", "quit"),
	),
}

type mainWindow struct {
	width    int
	height   int
	tc       *core.TerminalConfig
	textarea textarea.Model
	help     help.Model
	keys     keyMap
	status   *status
}

func initialModel(tc *core.TerminalConfig) mainWindow {
	ta := textarea.New()
	ta.Placeholder = "..."
	ta.Focus()

	return mainWindow{
		textarea: ta,
		help:     help.New(),
		keys:     keys,
		tc:       tc,
		status:   newStatus(""),
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

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Launch):
			m.tc.Commands = strings.Split(m.textarea.Value(), "\n")
			err := core.OpenWt(m.tc)
			if err != nil {
				m.status.Update(statusMsg{message: err.Error()})
				return m, nil
			} else {
				return m, tea.Quit
			}
		default:
			// Send all other keypresses to the textarea.
			var cmd tea.Cmd
			m.textarea, cmd = m.textarea.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

func (m mainWindow) View() string {
	margin := 2
	padding := 1

	boxWidth := m.width - margin*2
	boxHeight := m.height - margin*2 - m.status.GetHeight()

	m.help.Width = boxWidth - padding*2
	m.textarea.SetWidth(boxWidth - padding*2)
	// Calculate the height of other components and minus off
	m.textarea.SetHeight(boxHeight - 4)
	m.status.UpdateWindowWidth(boxWidth)

	// Footer
	miscBox := lipgloss.NewStyle().
		Width(boxWidth / 3).
		Align(lipgloss.Left).
		Render("Author: songlim327")
	titleBox := lipgloss.NewStyle().
		Width(boxWidth/3 - 1).
		Align(lipgloss.Center).
		Render("🍊 MPWT 🍊")
	versionBox := lipgloss.NewStyle().
		Width(boxWidth / 3).
		Align(lipgloss.Right).
		Bold(true).
		Render("0.1.1")
	footer := lipgloss.JoinHorizontal(lipgloss.Left, miscBox, titleBox, versionBox)

	// Content box
	return lipgloss.NewStyle().MarginLeft(margin).MarginTop(margin).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			m.status.View(),
			lipgloss.NewStyle().
				Width(boxWidth).
				Height(boxHeight).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color(BorderForegroundColor)).
				PaddingLeft(padding).
				PaddingRight(padding).
				PaddingTop(padding).
				Render(
					lipgloss.JoinVertical(lipgloss.Center,
						"Each line command will spawn a new terminal",
						m.textarea.View(),
						m.help.View(m.keys),
					),
					footer,
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
