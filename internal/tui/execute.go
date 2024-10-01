package tui

import (
	"mpwt/internal/core"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// keyMap defines a set of keybindings.
type keyMap struct {
	launch key.Binding
	back   key.Binding
	quit   key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view
// It is part of the key.Map interface
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.launch, k.back, k.quit}
}

// FullHelp returns keybindings to be shown in the full help view
// It is part of the key.Map interface
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.launch, k.back, k.quit},
	}
}

// keys implements key.Map interface and defines keyMap  of help menu
var keys = keyMap{
	launch: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "launch"),
	),
	back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back to main menu"),
	),
	quit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
}

// execute represents the command execution ui component
type execute struct {
	width     int
	height    int
	textarea  textarea.Model
	help      help.Model
	keys      keyMap
	tuiConfig *TuiConfig
}

// newExecute creates a new execute view
func newExecute(tuiConf *TuiConfig) *execute {
	ta := textarea.New()
	ta.Placeholder = "..."
	ta.Focus()

	return &execute{
		textarea:  ta,
		help:      help.New(),
		keys:      keys,
		tuiConfig: tuiConf,
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
		case key.Matches(msg, e.keys.quit):
			return e, tea.Quit

		case key.Matches(msg, e.keys.back):
			return e, tea.Batch(
				sendViewportUpdate(MainView),
				sendStatusUpdate(""),
			)

		case key.Matches(msg, e.keys.launch):
			// Split user input and compute the command
			cmds := strings.Split(e.textarea.Value(), "\n")
			e.tuiConfig.TerminalConfig.Commands = cmds
			cmdStr, err := core.OpenWt(e.tuiConfig.TerminalConfig)
			if err != nil {
				return e, sendStatusUpdate(err.Error())
			}

			// Execute the command
			cmd := exec.Command("cmd", "/C", cmdStr)
			if err := cmd.Run(); err != nil {
				return e, sendStatusUpdate(err.Error())
			}

			// Add command history to the database
			err = e.tuiConfig.Repository.InsertHistory(cmds, cmdStr)
			if err != nil {
				return e, sendStatusUpdate(err.Error())
			}

			return e, tea.Quit
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
