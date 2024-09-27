package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	textStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(TextColor))
)

// footer represents the state of footer component
type footer struct {
	width int
}

// newFooter creates a new footer instance
func newFooter() *footer {
	return &footer{}
}

// Init is the bubbletea package ELM architecture specific functions
func (f *footer) Init() tea.Cmd { return nil }

// Update is the bubbletea package ELM architecture specific functions
func (f *footer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return f, nil
}

// View is the bubbletea package ELM architecture specific functions
func (f *footer) View() string {
	miscBox := textStyle.
		Width(f.width / 3).
		Align(lipgloss.Left).
		Render("songlim327")
	titleBox := textStyle.
		Width(f.width/3 - 1).
		Align(lipgloss.Center).
		Render("MPWT")
	versionBox := textStyle.
		Width(f.width / 3).
		Align(lipgloss.Right).
		Bold(true).
		Render("0.1.1")

	return lipgloss.JoinHorizontal(lipgloss.Right, miscBox, titleBox, versionBox)
}
