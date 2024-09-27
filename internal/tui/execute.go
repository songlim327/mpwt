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
	Back   key.Binding
	Quit   key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view
// It is part of the key.Map interface
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Launch, k.Back, k.Quit}
}

// FullHelp returns keybindings to be shown in the full help view
// It is part of the key.Map interface
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Launch, k.Back, k.Quit},
	}
}

// keys implements key.Map interface and defines keyMap  of help menu
var keys = keyMap{
	Launch: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "launch"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back to main menu"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}

// execute represents the command execution ui component
type execute struct {
	width    int
	height   int
	textarea textarea.Model
	help     help.Model
	keys     keyMap
	tc       *core.TerminalConfig
}

// newExecute creates a new execute view
func newExecute(tc *core.TerminalConfig) *execute {
	ta := textarea.New()
	ta.Placeholder = "..."
	ta.Focus()

	return &execute{
		textarea: ta,
		help:     help.New(),
		keys:     keys,
		tc:       tc,
	}
}

// Init is the bubbletea package ELM architecture specific functions
func (e *execute) Init() tea.Cmd {
	return nil
}

// Update is the bubbletea package ELM architecture specific functions
func (e *execute) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, e.keys.Quit):
			return e, tea.Quit

		case key.Matches(msg, e.keys.Back):
			return e, tea.Batch(
				func() tea.Msg {
					return viewportMsg{viewport: Main}
				},
				func() tea.Msg {
					return statusMsg{message: ""}
				},
			)

		case key.Matches(msg, e.keys.Launch):
			e.tc.Commands = strings.Split(e.textarea.Value(), "\n")
			err := core.OpenWt(e.tc)
			if err != nil {
				return e, func() tea.Msg {
					return statusMsg{message: err.Error()}
				}
			} else {
				return e, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	e.textarea, cmd = e.textarea.Update(msg)
	return e, cmd
}

// View is the bubbletea package ELM architecture specific functions
func (e *execute) View() string {
	e.help.Width = e.width
	e.textarea.SetWidth(e.width)
	e.textarea.SetHeight(e.height - 1) // height of help model

	return lipgloss.JoinVertical(lipgloss.Left,
		e.textarea.View(),
		e.help.View(e.keys),
	)
}
