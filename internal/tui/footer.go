package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// footer represents the state of footer component
type footer struct {
	width int
	style lipgloss.Style
}

// newFooter creates a new footer instance
func newFooter() *footer {
	return &footer{
		style: lipgloss.NewStyle().Height(1).Foreground(lipgloss.Color(TextColor)),
	}
}

// Init is the bubbletea package ELM architecture specific functions
func (f *footer) Init() tea.Cmd { return nil }

// Update is the bubbletea package ELM architecture specific functions
func (f *footer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return f, nil
}

// View is the bubbletea package ELM architecture specific functions
func (f *footer) View() string {
	return f.style.Width(f.width).AlignHorizontal(lipgloss.Right).Render(
		// lipgloss.NewStyle().Foreground(lipgloss.Color(Yellow)).Underline(true).Render("Github"),
		lipgloss.NewStyle().Foreground(lipgloss.Color(GreenColor)).Render("0.1.1"),
	)
}
