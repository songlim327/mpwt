package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// status represents the state of status component
type status struct {
	message     string
	windowWidth int
	style       lipgloss.Style
}

// statusMsg represents a message struct to trigger status component updates
type statusMsg struct {
	message string
}

// newStatus creates a new status
func newStatus(defaultMessage string) *status {
	return &status{
		message: fmt.Sprintf("🍊 %s", defaultMessage),
		style: lipgloss.NewStyle().
			Height(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(BorderForegroundColor)),
	}
}

// updateWindowSize updates the main windows width to be used for the status component
func (s *status) updateWindowWidth(width int) {
	s.windowWidth = width
}

// getHeight returns the height of the status component
func (s *status) getHeight() int { return s.style.GetHeight() }

// Init is the bubbletea package ELM architecture specific functions
func (s status) Init() tea.Cmd { return nil }

// Update is the bubbletea package ELM architecture specific functions
func (s status) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statusMsg:
		s.message = fmt.Sprintf("🍊 %s", msg.message)
	}
	return s, nil
}

// View is the bubbletea package ELM architecture specific functions
func (s status) View() string {
	return s.style.Width(s.windowWidth).
		Render(s.message)
}
