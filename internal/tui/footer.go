package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// footer represents the state of footer component
type footer struct {
	windowWidth int
}

// newFooter creates a new footer instance
func newFooter() *footer {
	return &footer{}
}

// updateWindowWidth updates the main window width tobe used for the footer component
func (f *footer) updateWindowWidth(width int) {
	f.windowWidth = width
}

// Init is the bubbletea package ELM architecture specific functions
func (f *footer) Init() tea.Cmd { return nil }

// Update is the bubbletea package ELM architecture specific functions
func (f *footer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return f, nil
}

// View is the bubbletea package ELM architecture specific functions
func (f *footer) View() string {
	miscBox := lipgloss.NewStyle().
		Width(f.windowWidth / 3).
		Align(lipgloss.Left).
		Render("Author: songlim327")
	titleBox := lipgloss.NewStyle().
		Width(f.windowWidth/3 - 1).
		Align(lipgloss.Center).
		Render("🍊 MPWT 🍊")
	versionBox := lipgloss.NewStyle().
		Width(f.windowWidth / 3).
		Align(lipgloss.Right).
		Bold(true).
		Render("0.1.1")
	return lipgloss.JoinHorizontal(lipgloss.Left, miscBox, titleBox, versionBox)
}
